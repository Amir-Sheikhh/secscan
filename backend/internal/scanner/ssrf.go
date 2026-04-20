package scanner

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidScheme = errors.New("only http and https targets are allowed")
	ErrBlockedTarget = errors.New("target resolves to a blocked IP range")
)

var blockedCIDRs = mustCIDRs([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"169.254.0.0/16",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
})

type Target struct {
	RawURL     string
	URL        *url.URL
	Hostname   string
	Port       string
	ResolvedIP net.IP
	Client     *http.Client
}

func PrepareTarget(raw string) (*Target, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, fmt.Errorf("parse target: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, ErrInvalidScheme
	}
	if parsed.Hostname() == "" {
		return nil, errors.New("target hostname is required")
	}

	hostname := parsed.Hostname()
	ip, err := resolveAndValidate(hostname)
	if err != nil {
		return nil, err
	}

	port := parsed.Port()
	if port == "" {
		if parsed.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return &Target{
		RawURL:     parsed.String(),
		URL:        parsed,
		Hostname:   hostname,
		Port:       port,
		ResolvedIP: ip,
		Client:     safeHTTPClient(hostname, port, ip),
	}, nil
}

func resolveAndValidate(hostname string) (net.IP, error) {
	if ip := net.ParseIP(hostname); ip != nil {
		if isBlockedIP(ip) {
			return nil, ErrBlockedTarget
		}
		return ip, nil
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("dns lookup failed: %w", err)
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return nil, ErrBlockedTarget
		}
	}
	for _, ip := range ips {
		if ip.To4() != nil || ip.To16() != nil {
			return ip, nil
		}
	}
	return nil, errors.New("no routable IP found for target")
}

func safeHTTPClient(hostname, port string, ip net.IP) *http.Client {
	dialer := &net.Dialer{Timeout: 4 * time.Second}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: hostname,
			MinVersion: tls.VersionTLS12,
		},
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		},
		Proxy:               nil,
		ForceAttemptHTTP2:   true,
		ResponseHeaderTimeout: 6 * time.Second,
		DisableKeepAlives:   false,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   8 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			if req.URL.Hostname() != hostname {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

func isBlockedIP(ip net.IP) bool {
	if ip == nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	for _, cidr := range blockedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func mustCIDRs(raw []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(raw))
	for _, item := range raw {
		_, block, err := net.ParseCIDR(item)
		if err != nil {
			panic(err)
		}
		nets = append(nets, block)
	}
	return nets
}
