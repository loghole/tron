package healthcheck

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// Check is a basic health check func.
type Check func() error

// Checker is a basic health checker.
type Checker interface {
	// AddLiveness adds a check that indicates that this instance of the
	// application should be destroyed or restarted. Every liveness check
	// is also included as a readiness check.
	AddLiveness(name string, check Check)

	// AddReadiness adds a check that indicates that this instance of the
	// application is currently unable to serve requests because of an upstream
	// or some transient failure.
	AddReadiness(name string, check Check)

	// LivenessHandler is the HTTP handler for just the /live endpoint.
	LivenessHandler(w http.ResponseWriter, r *http.Request)

	// ReadinessHandler is the HTTP handler for just the /ready endpoint.
	ReadinessHandler(w http.ResponseWriter, r *http.Request)
}

type checker struct {
	liveness  map[string]Check
	readiness map[string]Check
	mu        sync.RWMutex
}

// NewChecker creates a new basic health checker.
func NewChecker() Checker {
	handler := &checker{
		liveness:  make(map[string]Check),
		readiness: make(map[string]Check),
	}

	return handler
}

func (c *checker) AddLiveness(name string, check Check) {
	c.mu.Lock()
	c.liveness[name] = check
	c.mu.Unlock()
}

func (c *checker) AddReadiness(name string, check Check) {
	c.mu.Lock()
	c.readiness[name] = check
	c.mu.Unlock()
}

func (c *checker) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.handle(w, r, c.liveness)
}

func (c *checker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.handle(w, r, c.liveness, c.readiness)
}

func (c *checker) handle(w http.ResponseWriter, r *http.Request, list ...map[string]Check) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	var (
		result = make(map[string]string)
		status = http.StatusOK
	)

	for _, checks := range list {
		for name, check := range checks {
			if err := check(); err != nil {
				result[name] = err.Error()
			}
		}
	}

	if len(result) > 0 {
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(status)

	if r.URL.Query().Get("full") != "1" {
		_, _ = io.WriteString(w, "{}\n")

		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(result) // nolint:errchkjson // not need check error.
}
