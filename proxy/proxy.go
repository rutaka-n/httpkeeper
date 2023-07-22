package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"httpkeeper/config"
	"httpkeeper/token"

	gs "github.com/rutaka-n/genericset"
)

// Proxy - contains settings for proxy server
type Proxy struct {
	Services          map[string]*Service
	Secret            []byte
	ServiceName       string
	server            http.Server
	InvalidatedTokens gs.Set[string]
}

// New returns Proxy initialized with addr
func New(addr string) *Proxy {
	return &Proxy{
		server: http.Server{
			Addr: addr,
		},
	}
}

// SetSecret sets the secret value, that used to validate JWT
func (p *Proxy) SetSecret(secret []byte) {
	p.Secret = secret
}

// SetServiceName sets name of issuer, that used to validate JWT
func (p *Proxy) SetServiceName(name string) {
	p.ServiceName = name
}

// SetServices sets list of config.Service
func (p *Proxy) SetServices(services []config.Service) {
	svcMap := make(map[string]*Service, len(services))
	for i := range services {
		svcMap[services[i].VirtualHost] = NewService(&services[i])
	}
	p.Services = svcMap
}

// SetInvalidatedTokens sets the list of invalidated tokens
func (p *Proxy) SetInvalidatedTokens(tokens []string) {
	set := gs.New[string]()
	set.Add(tokens...)
	p.InvalidatedTokens = set
}

// ListenAndServe binds port and starts listening
func (p *Proxy) ListenAndServe() error {
	p.server.Handler = http.HandlerFunc(p.handler)
	return p.server.ListenAndServe()
}

// Shutdown gracefully stops the http-server
func (p *Proxy) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *Proxy) getHost(req *http.Request) string {
	splited := strings.Split(req.Host, ":")
	if len(splited) > 0 {
		return splited[0]
	}
	return ""
}

func (p *Proxy) authWithToken(req *http.Request) error {
	value, ok := req.Header["Authorization"]
	if !ok {
		return errors.New("Authorization header not found")
	}
	bearer := []string{}
	if len(value) > 0 {
		bearer = strings.Split(value[0], " ")
	}
	if len(bearer) == 2 && bearer[0] == "Bearer" {
		if p.InvalidatedTokens.IsElement(bearer[1]) {
			return errors.New("token is invalidated")
		}
		payload, err := token.Validate(p.Secret, bearer[1])
		if err != nil {
			return err
		}
		if payload.Issuer != p.ServiceName {
			return errors.New("token issuer mismatch")
		}
	} else {
		return errors.New("Authorization header is invalid")
	}
	return nil
}

func (p *Proxy) handler(rw http.ResponseWriter, req *http.Request) {
	requestLogMsg := fmt.Sprintf("received request %s %s %v from %s", req.Method, req.Host, req.RequestURI, req.RemoteAddr)
	host := p.getHost(req)
	service, ok := p.Services[host]
	if !ok {
		log.Printf(requestLogMsg+" - %d, service not configured", http.StatusForbidden)
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	if service.IsJWTEnabled() {
		err := p.authWithToken(req)
		if err != nil {
			log.Printf(requestLogMsg+" - %d, %v", http.StatusProxyAuthRequired, err)
			rw.WriteHeader(http.StatusProxyAuthRequired)
			return
		}
	}

	if !service.Allow() {
		log.Printf(requestLogMsg+" - %d, rate limit exceeded", http.StatusTooManyRequests)
		rw.WriteHeader(http.StatusTooManyRequests)
		return
	}

	// set req Host, URL and Request URI to forward a request to the origin server
	req.Host = service.URL.Host
	req.URL.Host = service.URL.Host
	req.URL.Scheme = service.URL.Scheme
	req.RequestURI = "" // Set to empty line since RequestURI shouldn't set for clients request
	if service.IsBasicAuth() {
		req.SetBasicAuth(service.BasicAuth.User, service.BasicAuth.Password)
	}

	// save the response from the origin server
	originServerResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(rw, err)
		log.Printf(requestLogMsg+" - %d", http.StatusInternalServerError)
		return
	}

	// return response to the client
	io.Copy(rw, originServerResponse.Body)
	log.Printf(requestLogMsg+" - %d", http.StatusOK)
}
