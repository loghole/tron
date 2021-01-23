package tron

import (
	"errors"
	"sync/atomic"

	"github.com/loghole/tron/healthcheck"
)

var errAppIsNotRunning = errors.New("app is not run yet")

type health struct {
	healthcheck.Checker

	state int32
}

func (h *health) init() {
	h.Checker = healthcheck.NewChecker()
	h.Checker.AddReadiness("tron", h.readiness)
}

func (h *health) readiness() error {
	if atomic.LoadInt32(&h.state) != 1 {
		return errAppIsNotRunning
	}

	return nil
}

func (h *health) setReady() {
	atomic.StoreInt32(&h.state, 1)
}

func (h *health) setUnready() {
	atomic.StoreInt32(&h.state, 0)
}
