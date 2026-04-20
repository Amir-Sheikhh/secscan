package scanner

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/amir-sheikh/secscan/backend/internal/model"
)

type PortsModule struct{}

func NewPortsModule() Module { return PortsModule{} }

func (PortsModule) Name() string { return "ports" }

func (PortsModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	commonPorts := []int{21, 22, 25, 53, 80, 110, 143, 443, 445, 587, 993, 995, 3000, 3306, 5432, 6379, 8080, 8443}
	services := map[int]string{
		21: "ftp", 22: "ssh", 25: "smtp", 53: "dns", 80: "http", 110: "pop3", 143: "imap",
		443: "https", 445: "smb", 587: "smtp-submission", 993: "imaps", 995: "pop3s",
		3000: "dev-http", 3306: "mysql", 5432: "postgres", 6379: "redis", 8080: "http-alt", 8443: "https-alt",
	}
	sensitive := map[int]struct{}{21: {}, 22: {}, 445: {}, 3306: {}, 5432: {}, 6379: {}}

	var (
		mu       sync.Mutex
		openList []map[string]any
		findings []model.Finding
	)
	semaphore := make(chan struct{}, 6)
	var wg sync.WaitGroup

	for _, port := range commonPorts {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
			}
			defer func() { <-semaphore }()

			dialer := &net.Dialer{Timeout: 400 * time.Millisecond}
			conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(target.ResolvedIP.String(), fmt.Sprint(port)))
			if err != nil {
				return
			}
			_ = conn.Close()

			mu.Lock()
			defer mu.Unlock()
			openList = append(openList, map[string]any{
				"port":    port,
				"service": services[port],
			})
			if _, isSensitive := sensitive[port]; isSensitive {
				findings = append(findings, model.Finding{
					Title:          fmt.Sprintf("Sensitive service exposed on %d", port),
					Severity:       "medium",
					Category:       "attack-surface",
					Description:    "A commonly abused administrative or database port is reachable from the target.",
					Recommendation: "Restrict exposure with a firewall, reverse proxy, or network ACL.",
					Evidence:       fmt.Sprintf("open port %d (%s)", port, services[port]),
				})
			}
		}(port)
	}
	wg.Wait()

	slices.SortFunc(openList, func(a, b map[string]any) int {
		return a["port"].(int) - b["port"].(int)
	})

	score := 100 - len(findings)*12
	if score < 40 {
		score = 40
	}
	if len(openList) == 0 {
		score = 95
	}

	return model.ModuleResult{
		Name:     "ports",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  fmt.Sprintf("%d common ports open", len(openList)),
		Findings: findings,
		Details: map[string]any{
			"openPorts": openList,
			"scanned":   len(commonPorts),
		},
	}
}

type HeadersModule struct{}

func NewHeadersModule() Module { return HeadersModule{} }

func (HeadersModule) Name() string { return "headers" }

func (HeadersModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target.RawURL, nil)
	req.Header.Set("User-Agent", "SecScan/1.0")
	resp, err := target.Client.Do(req)
	if err != nil {
		return failedModule("headers", "header probe failed", err)
	}
	defer resp.Body.Close()

	required := []string{
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"Referrer-Policy",
		"Permissions-Policy",
		"Cross-Origin-Opener-Policy",
		"Cross-Origin-Resource-Policy",
	}

	headers := map[string]any{}
	findings := []model.Finding{}
	present := 0

	for _, header := range required {
		value := resp.Header.Get(header)
		headers[header] = value
		if value != "" {
			present++
			continue
		}
		if header == "Strict-Transport-Security" && target.URL.Scheme != "https" {
			continue
		}
		findings = append(findings, model.Finding{
			Title:          header + " missing",
			Severity:       "medium",
			Category:       "security-headers",
			Description:    "The response is missing a recommended defensive HTTP header.",
			Recommendation: "Set the missing header on the reverse proxy or application layer.",
			Evidence:       header,
		})
	}

	score := int(float64(present) / float64(len(required)) * 100)
	if target.URL.Scheme != "https" && score > 10 {
		score -= 10
	}

	return model.ModuleResult{
		Name:     "headers",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  fmt.Sprintf("%d/%d key security headers present", present, len(required)),
		Findings: findings,
		Details: map[string]any{
			"statusCode": resp.StatusCode,
			"headers":    headers,
		},
	}
}

type TLSModule struct{}

func NewTLSModule() Module { return TLSModule{} }

func (TLSModule) Name() string { return "tls" }

func (TLSModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	if target.URL.Scheme != "https" {
		return model.ModuleResult{
			Name:     "tls",
			Status:   model.ModuleCompleted,
			Score:    45,
			Severity: "high",
			Summary:  "target does not use HTTPS",
			Findings: []model.Finding{{
				Title:          "HTTPS disabled",
				Severity:       "high",
				Category:       "transport-security",
				Description:    "The application is served over plain HTTP.",
				Recommendation: "Serve the application via HTTPS and redirect HTTP traffic.",
				Evidence:       target.RawURL,
			}},
		}
	}

	type versionProbe struct {
		Label   string
		Version uint16
	}

	probes := []versionProbe{
		{"TLS 1.0", tls.VersionTLS10},
		{"TLS 1.1", tls.VersionTLS11},
		{"TLS 1.2", tls.VersionTLS12},
		{"TLS 1.3", tls.VersionTLS13},
	}

	supported := []string{}
	findings := []model.Finding{}
	for _, probe := range probes {
		ok, cert, cipher, err := probeTLSVersion(ctx, target, probe.Version)
		if err != nil || !ok {
			continue
		}
		supported = append(supported, probe.Label)
		if probe.Version < tls.VersionTLS12 {
			findings = append(findings, model.Finding{
				Title:          probe.Label + " supported",
				Severity:       "high",
				Category:       "transport-security",
				Description:    "Legacy TLS versions are still accepted by the server.",
				Recommendation: "Disable TLS 1.0/1.1 and keep only TLS 1.2+ enabled.",
				Evidence:       cipher,
			})
		}
		if cert != nil {
			daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
			if daysLeft < 30 {
				findings = append(findings, model.Finding{
					Title:          "Certificate expires soon",
					Severity:       "medium",
					Category:       "transport-security",
					Description:    "The leaf certificate will expire in less than 30 days.",
					Recommendation: "Rotate or renew the certificate before expiration.",
					Evidence:       cert.NotAfter.Format(time.RFC3339),
				})
			}
		}
	}

	score := 100
	for _, finding := range findings {
		switch finding.Severity {
		case "high":
			score -= 25
		case "medium":
			score -= 10
		}
	}
	if score < 35 {
		score = 35
	}

	return model.ModuleResult{
		Name:     "tls",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  fmt.Sprintf("%d TLS versions accepted", len(supported)),
		Findings: findings,
		Details: map[string]any{
			"supportedVersions": supported,
		},
	}
}

type FuzzModule struct{}

func NewFuzzModule() Module { return FuzzModule{} }

func (FuzzModule) Name() string { return "fuzz" }

func (FuzzModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	paths := []string{"admin", "login", "dashboard", "api", "robots.txt", ".env", "backup", "uploads"}
	found := []map[string]any{}
	findings := []model.Finding{}

	for _, path := range paths {
		joined := strings.TrimRight(target.RawURL, "/") + "/" + path
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, joined, nil)
		req.Header.Set("User-Agent", "SecScan/1.0")
		resp, err := target.Client.Do(req)
		if err != nil {
			continue
		}
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			continue
		}
		found = append(found, map[string]any{
			"path":   path,
			"status": resp.StatusCode,
		})
		if path == ".env" || path == "backup" {
			findings = append(findings, model.Finding{
				Title:          "Potential sensitive path exposed",
				Severity:       "high",
				Category:       "content-discovery",
				Description:    "A path commonly associated with secrets or backups responded to a probe.",
				Recommendation: "Remove sensitive files from the web root and block direct access.",
				Evidence:       fmt.Sprintf("/%s -> %d", path, resp.StatusCode),
			})
		}
	}

	score := 100 - len(findings)*20
	if score < 45 {
		score = 45
	}
	if len(found) == 0 {
		score = 92
	}

	return model.ModuleResult{
		Name:     "fuzz",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  fmt.Sprintf("%d interesting paths found", len(found)),
		Findings: findings,
		Details: map[string]any{
			"paths": found,
		},
	}
}

type XSSModule struct{}

func NewXSSModule() Module { return XSSModule{} }

func (XSSModule) Name() string { return "xss" }

func (XSSModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	marker := "secscan-xss-marker"
	probeURL := withQuery(target.URL, "secscan_reflection", marker)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, probeURL, nil)
	req.Header.Set("User-Agent", "SecScan/1.0")
	resp, err := target.Client.Do(req)
	if err != nil {
		return failedModule("xss", "reflection probe failed", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	text := string(body)
	reflected := strings.Contains(text, marker)
	escaped := strings.Contains(text, url.QueryEscape(marker)) || strings.Contains(text, "secscan_reflection")

	findings := []model.Finding{}
	score := 95
	summary := "no direct reflection observed"
	if reflected && !escaped {
		score = 68
		summary = "reflection detected in response body"
		findings = append(findings, model.Finding{
			Title:          "Reflected marker found in HTML response",
			Severity:       "medium",
			Category:       "xss",
			Description:    "User-controlled input was reflected without clear output encoding.",
			Recommendation: "HTML-encode reflected content and enforce a strict Content-Security-Policy.",
			Evidence:       marker,
		})
	}

	return model.ModuleResult{
		Name:     "xss",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  summary,
		Findings: findings,
		Details: map[string]any{
			"probeUrl":   probeURL,
			"reflected":  reflected,
			"bodyLength": len(body),
		},
	}
}

type SQLIModule struct {
	enableActive bool
}

func NewSQLIModule(enableActive bool) Module {
	return SQLIModule{enableActive: enableActive}
}

func (SQLIModule) Name() string { return "sqli" }

func (m SQLIModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	probes := []struct {
		name  string
		value string
	}{
		{name: "baseline", value: "1"},
		{name: "quote", value: "'"},
	}
	if m.enableActive {
		probes = append(probes,
			struct {
				name  string
				value string
			}{name: "booleanTrue", value: "1 OR 1=1"},
			struct {
				name  string
				value string
			}{name: "booleanFalse", value: "1 AND 1=2"},
		)
	}

	results := map[string]any{}
	findings := []model.Finding{}
	baselineSize := 0
	booleanDiff := false

	for _, probe := range probes {
		probeURL := withQuery(target.URL, "secscan_id", probe.value)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, probeURL, nil)
		req.Header.Set("User-Agent", "SecScan/1.0")
		start := time.Now()
		resp, err := target.Client.Do(req)
		if err != nil {
			results[probe.name] = map[string]any{"error": err.Error()}
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
		_ = resp.Body.Close()
		content := string(body)
		results[probe.name] = map[string]any{
			"status":     resp.StatusCode,
			"durationMs": time.Since(start).Milliseconds(),
			"size":       len(body),
		}

		if probe.name == "baseline" {
			baselineSize = len(body)
		}
		if probe.name == "quote" && containsSQLError(content) {
			findings = append(findings, model.Finding{
				Title:          "SQL error signature detected",
				Severity:       "high",
				Category:       "sqli",
				Description:    "The response included database error text after a malformed input probe.",
				Recommendation: "Use parameterized queries and avoid echoing backend errors to the client.",
				Evidence:       "database error text in response body",
			})
		}
		if probe.name == "booleanFalse" && baselineSize > 0 && absInt(len(body)-baselineSize) > 120 {
			booleanDiff = true
		}
	}

	if booleanDiff {
		findings = append(findings, model.Finding{
			Title:          "Probe responses changed materially",
			Severity:       "medium",
			Category:       "sqli",
			Description:    "Boolean-style test inputs produced a notably different response shape.",
			Recommendation: "Review server-side input handling and verify all queries are parameterized.",
			Evidence:       "response length delta exceeded heuristic threshold",
		})
	}

	score := 96
	if len(findings) > 0 {
		score = 60
	}

	return model.ModuleResult{
		Name:     "sqli",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  "heuristic SQL injection probes completed",
		Findings: findings,
		Details: map[string]any{
			"activeProbesEnabled": m.enableActive,
			"results":             results,
		},
	}
}

type CVEModule struct {
	enableIntel bool
}

func NewCVEModule(enableIntel bool) Module {
	return CVEModule{enableIntel: enableIntel}
}

func (CVEModule) Name() string { return "cve" }

func (m CVEModule) Run(ctx context.Context, target *Target) model.ModuleResult {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target.RawURL, nil)
	req.Header.Set("User-Agent", "SecScan/1.0")
	resp, err := target.Client.Do(req)
	if err != nil {
		return failedModule("cve", "technology fingerprint failed", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	detections := detectTechnologies(resp, string(body))

	findings := []model.Finding{}
	advisories := []map[string]any{}

	for _, tech := range detections {
		if tech.Version == "" {
			continue
		}
		if !m.enableIntel {
			continue
		}
		refs, err := queryOSV(ctx, tech.Name, tech.Version)
		if err != nil {
			continue
		}
		for _, ref := range refs {
			advisories = append(advisories, ref)
			findings = append(findings, model.Finding{
				Title:          fmt.Sprintf("%s advisory candidate", tech.Name),
				Severity:       severityFromString(fmt.Sprint(ref["severity"])),
				Category:       "cve",
				Description:    "An external advisory lookup returned a potential match for the detected version.",
				Recommendation: "Verify the advisory and upgrade to a fixed version when applicable.",
				Evidence:       fmt.Sprintf("%v", ref["id"]),
			})
		}
	}

	score := 90
	if len(findings) > 0 {
		score = 65
	}
	if len(detections) == 0 {
		score = 80
	}

	details := map[string]any{
		"technologies": detections,
		"advisories":   advisories,
	}

	return model.ModuleResult{
		Name:     "cve",
		Status:   model.ModuleCompleted,
		Score:    score,
		Severity: maxSeverity(findings),
		Summary:  fmt.Sprintf("%d technologies detected", len(detections)),
		Findings: findings,
		Details:  details,
	}
}

type detectedTech struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Source  string `json:"source"`
}

func detectTechnologies(resp *http.Response, body string) []detectedTech {
	seen := map[string]struct{}{}
	output := []detectedTech{}

	push := func(name, version, source string) {
		key := name + "@" + version + "#" + source
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		output = append(output, detectedTech{Name: name, Version: version, Source: source})
	}

	if server := resp.Header.Get("Server"); server != "" {
		name, version := splitTech(server)
		push(name, version, "Server")
	}
	if powered := resp.Header.Get("X-Powered-By"); powered != "" {
		name, version := splitTech(powered)
		push(name, version, "X-Powered-By")
	}

	metaGenerator := regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']([^"']+)["']`)
	if match := metaGenerator.FindStringSubmatch(body); len(match) == 2 {
		name, version := splitTech(match[1])
		push(name, version, "meta-generator")
	}

	scriptVersion := regexp.MustCompile(`(?i)(jquery|bootstrap|wordpress)[^0-9]{0,4}([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	for _, match := range scriptVersion.FindAllStringSubmatch(body, -1) {
		if len(match) == 3 {
			push(strings.Title(strings.ToLower(match[1])), match[2], "html")
		}
	}

	slices.SortFunc(output, func(a, b detectedTech) int {
		return strings.Compare(a.Name, b.Name)
	})
	return output
}

func splitTech(raw string) (string, string) {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, "(", " ")
	raw = strings.ReplaceAll(raw, ")", " ")
	versionRegex := regexp.MustCompile(`([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	version := ""
	if match := versionRegex.FindStringSubmatch(raw); len(match) == 2 {
		version = match[1]
		raw = strings.Replace(raw, match[1], "", 1)
	}
	fields := strings.Fields(raw)
	name := "unknown"
	if len(fields) > 0 {
		name = strings.Trim(fields[0], "/")
	}
	if name == "" {
		name = "unknown"
	}
	return name, version
}

func queryOSV(ctx context.Context, name, version string) ([]map[string]any, error) {
	payload := map[string]any{
		"package": map[string]string{
			"name":      name,
			"ecosystem": "OSS-Fuzz",
		},
		"version": version,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.osv.dev/v1/query", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decoded struct {
		Vulns []struct {
			ID       string `json:"id"`
			Severity []struct {
				Type  string `json:"type"`
				Score string `json:"score"`
			} `json:"severity"`
		} `json:"vulns"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	out := make([]map[string]any, 0, len(decoded.Vulns))
	for _, vuln := range decoded.Vulns {
		level := "medium"
		if len(vuln.Severity) > 0 {
			level = vuln.Severity[0].Score
		}
		out = append(out, map[string]any{
			"id":       vuln.ID,
			"severity": level,
		})
	}
	return out, nil
}

func probeTLSVersion(ctx context.Context, target *Target, version uint16) (bool, *x509.Certificate, string, error) {
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	rawConn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(target.ResolvedIP.String(), target.Port))
	if err != nil {
		return false, nil, "", err
	}
	defer rawConn.Close()

	conn := tls.Client(rawConn, &tls.Config{
		ServerName:         target.Hostname,
		MinVersion:         version,
		MaxVersion:         version,
		InsecureSkipVerify: false,
	})
	if err := conn.HandshakeContext(ctx); err != nil {
		return false, nil, "", err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	var cert *x509.Certificate
	if len(state.PeerCertificates) > 0 {
		cert = state.PeerCertificates[0]
	}
	return true, cert, tls.CipherSuiteName(state.CipherSuite), nil
}

func failedModule(name, summary string, err error) model.ModuleResult {
	return model.ModuleResult{
		Name:     name,
		Status:   model.ModuleFailed,
		Score:    0,
		Severity: "high",
		Summary:  summary,
		Error:    err.Error(),
		Findings: []model.Finding{{
			Title:          name + " probe failed",
			Severity:       "high",
			Category:       "scanner",
			Description:    "The module could not complete its probe.",
			Recommendation: "Check target reachability and retry.",
			Evidence:       err.Error(),
		}},
	}
}

func withQuery(base *url.URL, key, value string) string {
	nextURL := *base
	query := nextURL.Query()
	query.Set(key, value)
	nextURL.RawQuery = query.Encode()
	return nextURL.String()
}

func containsSQLError(body string) bool {
	patterns := []string{
		"sql syntax",
		"mysql",
		"postgresql",
		"sqlite",
		"syntax error",
		"ora-",
		"unclosed quotation mark",
	}
	lower := strings.ToLower(body)
	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func maxSeverity(findings []model.Finding) string {
	severityOrder := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
	best := "info"
	rank := 0
	for _, finding := range findings {
		level := strings.ToLower(finding.Severity)
		if severityOrder[level] > rank {
			rank = severityOrder[level]
			best = level
		}
	}
	if best == "info" && len(findings) == 0 {
		return "info"
	}
	return best
}

func severityFromString(value string) string {
	lower := strings.ToLower(value)
	switch {
	case strings.Contains(lower, "critical"), strings.HasPrefix(lower, "9"), strings.HasPrefix(lower, "10"):
		return "critical"
	case strings.Contains(lower, "high"), strings.HasPrefix(lower, "8"), strings.HasPrefix(lower, "7"):
		return "high"
	case strings.Contains(lower, "medium"), strings.HasPrefix(lower, "6"), strings.HasPrefix(lower, "5"), strings.HasPrefix(lower, "4"):
		return "medium"
	default:
		return "low"
	}
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
