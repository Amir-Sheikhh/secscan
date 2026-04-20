package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/amir-sheikh/secscan/backend/internal/model"
	"github.com/amir-sheikh/secscan/backend/internal/scanner"
	"github.com/amir-sheikh/secscan/backend/internal/storage"
)

type Config struct {
	AllowedOrigins      []string
	ScanTimeout         time.Duration
	EnableActiveProbes  bool
	EnableExternalIntel bool
}

type Server struct {
	service        *scanner.Service
	allowedOrigins []string
}

func NewServer(cfg Config) *Server {
	store := storage.NewMemoryStore()
	service := scanner.NewService(store, scanner.Config{
		ScanTimeout:         cfg.ScanTimeout,
		EnableActiveProbes:  cfg.EnableActiveProbes,
		EnableExternalIntel: cfg.EnableExternalIntel,
	})

	return &Server{
		service:        service,
		allowedOrigins: cfg.AllowedOrigins,
	}
}

func (s *Server) Router() http.Handler {
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
		}
	}()

	s.applyCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/health":
		s.health(w)
	case r.Method == http.MethodPost && r.URL.Path == "/api/scan":
		s.createScan(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/scan/"):
		s.dispatchScanRoutes(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) dispatchScanRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/scan/")
	switch {
	case r.Method == http.MethodGet && strings.HasSuffix(path, "/stream"):
		id := strings.TrimSuffix(strings.TrimSuffix(path, "/stream"), "/")
		s.streamScan(w, r, id)
	case r.Method == http.MethodGet && strings.HasSuffix(path, "/report.pdf"):
		id := strings.TrimSuffix(strings.TrimSuffix(path, "/report.pdf"), "/")
		s.reportPDF(w, id)
	case r.Method == http.MethodGet && path != "" && !strings.Contains(path, "/"):
		s.getScan(w, path)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) applyCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	allowed := false
	for _, candidate := range s.allowedOrigins {
		if candidate == "*" || strings.EqualFold(candidate, origin) {
			allowed = true
			break
		}
	}
	if !allowed {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
}

func (s *Server) health(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"time":    time.Now().UTC().Format(time.RFC3339),
		"modules": model.DefaultModules,
	})
}

func (s *Server) createScan(w http.ResponseWriter, r *http.Request) {
	var req model.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	scan, err := s.service.Start(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusAccepted, scan)
}

func (s *Server) getScan(w http.ResponseWriter, id string) {
	scan, err := s.service.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, scan)
}

func (s *Server) streamScan(w http.ResponseWriter, r *http.Request, id string) {
	ch, unsubscribe, err := s.service.Subscribe(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	defer unsubscribe()

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, errors.New("streaming unsupported"))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case scan := <-ch:
			payload, _ := json.Marshal(scan)
			_, _ = w.Write([]byte("data: " + string(payload) + "\n\n"))
			flusher.Flush()
			if scan.Status == model.ScanCompleted || scan.Status == model.ScanFailed {
				return
			}
		case <-ticker.C:
			_, _ = w.Write([]byte(": keepalive\n\n"))
			flusher.Flush()
		}
	}
}

func (s *Server) reportPDF(w http.ResponseWriter, id string) {
	scan, err := s.service.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	data := renderPDF(scan)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=secscan-"+scan.ID+".pdf")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"error":"json encode failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}
