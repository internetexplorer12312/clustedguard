package health

import "ClusterGuard/internal/domain"

// Registry maps check types to HealthChecker implementations.
type Registry struct {
	checkers map[string]domain.HealthChecker
}

func NewRegistry(checkers ...domain.HealthChecker) *Registry {
	m := make(map[string]domain.HealthChecker, len(checkers))
	for _, c := range checkers {
		m[c.Type()] = c
	}
	return &Registry{checkers: m}
}

func (r *Registry) Get(checkType string) (domain.HealthChecker, bool) {
	if checkType == "" {
		checkType = "tcp"
	}
	c, ok := r.checkers[checkType]
	return c, ok
}

var _ domain.HealthCheckerRegistry = (*Registry)(nil)
