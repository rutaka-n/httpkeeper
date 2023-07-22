package proxy

import (
	"net/url"
	"time"

	"httpkeeper/config"

	"golang.org/x/time/rate"
)

// BasicAuth contains credantials for basic auth
type BasicAuth struct {
	User     string
	Password string
}

// Service represent service for the virtual host
type Service struct {
	VirtualHost string
	URL         url.URL
	BasicAuth   *BasicAuth
	RateLimit   *rate.Limiter
	JWT         bool
}

// NewService returns Service for with configured parameters
func NewService(cfg *config.Service) *Service {
	s := Service{
		VirtualHost: cfg.VirtualHost,
		URL:         url.URL(cfg.URL),
	}

	s.JWT = cfg.JWT == "enabled" || cfg.JWT == "true"

	if cfg.User != "" && cfg.Password != "" {
		s.BasicAuth = &BasicAuth{User: cfg.User, Password: cfg.Password}
	}
	if cfg.RateLimit > 0 {
		s.RateLimit = rate.NewLimiter(rate.Every(time.Second/time.Duration(cfg.RateLimit)), cfg.RateLimit)
	}

	return &s
}

// IsBasicAuth return true if basic auth credantials are set
func (s *Service) IsBasicAuth() bool {
	return s.BasicAuth != nil
}

// IsJWTEnabled returns true if JWT auth is enabled
func (s *Service) IsJWTEnabled() bool {
	return s.JWT
}

// Allow returns true if rate limiter value is not exceeded
func (s *Service) Allow() bool {
	return s.RateLimit.Allow()
}
