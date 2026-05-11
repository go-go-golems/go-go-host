package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	StatusStarting Status = "starting"
	StatusReady    Status = "ready"
	StatusFailed   Status = "failed"
	StatusStopped  Status = "stopped"
	StatusDraining Status = "draining"
)

type RuntimeStatus struct {
	SiteID        string    `json:"siteId"`
	OrgID         string    `json:"orgId"`
	DeploymentID  string    `json:"deploymentId"`
	Hosts         []string  `json:"hosts"`
	Status        Status    `json:"status"`
	StartedAt     time.Time `json:"startedAt"`
	LastError     string    `json:"lastError,omitempty"`
	RequestsTotal uint64    `json:"requestsTotal"`
	ErrorsTotal   uint64    `json:"errorsTotal"`
}

type Summary struct {
	ActiveSites int             `json:"activeSites"`
	Hosts       []string        `json:"hosts"`
	Runtimes    []RuntimeStatus `json:"runtimes"`
}

type Supervisor struct {
	mu     sync.RWMutex
	bySite map[string]*SiteRuntime
	byHost map[string]*SiteRuntime
	status map[string]RuntimeStatus
	specs  map[string]Spec
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		bySite: map[string]*SiteRuntime{},
		byHost: map[string]*SiteRuntime{},
		status: map[string]RuntimeStatus{},
		specs:  map[string]Spec{},
	}
}

func (s *Supervisor) Activate(ctx context.Context, spec Spec) error {
	if spec.SiteID == "" {
		return fmt.Errorf("site id is required")
	}
	if len(spec.Hosts) == 0 {
		return fmt.Errorf("at least one host is required")
	}
	s.setStatus(spec.SiteID, RuntimeStatus{SiteID: spec.SiteID, OrgID: spec.OrgID, DeploymentID: spec.DeploymentID, Hosts: normalizeHosts(spec.Hosts), Status: StatusStarting})
	next, err := NewSiteRuntime(ctx, spec)
	if err != nil {
		s.setStatus(spec.SiteID, RuntimeStatus{SiteID: spec.SiteID, OrgID: spec.OrgID, DeploymentID: spec.DeploymentID, Hosts: normalizeHosts(spec.Hosts), Status: StatusFailed, LastError: err.Error()})
		return err
	}
	if err := next.HealthCheck(ctx); err != nil {
		_ = next.Close(ctx)
		s.setStatus(spec.SiteID, RuntimeStatus{SiteID: spec.SiteID, OrgID: spec.OrgID, DeploymentID: spec.DeploymentID, Hosts: normalizeHosts(spec.Hosts), Status: StatusFailed, LastError: err.Error()})
		return err
	}

	var old *SiteRuntime
	s.mu.Lock()
	old = s.bySite[spec.SiteID]
	if old != nil {
		for _, host := range old.spec.Hosts {
			delete(s.byHost, normalizeHost(host))
		}
	}
	s.bySite[spec.SiteID] = next
	for _, host := range spec.Hosts {
		s.byHost[normalizeHost(host)] = next
	}
	previous := s.status[spec.SiteID]
	s.status[spec.SiteID] = RuntimeStatus{SiteID: spec.SiteID, OrgID: spec.OrgID, DeploymentID: spec.DeploymentID, Hosts: normalizeHosts(spec.Hosts), Status: StatusReady, StartedAt: next.StartedAt(), RequestsTotal: previous.RequestsTotal, ErrorsTotal: previous.ErrorsTotal}
	s.specs[spec.SiteID] = spec
	s.mu.Unlock()

	if old != nil {
		go func() { _ = old.Close(context.Background()) }()
	}
	return nil
}

func (s *Supervisor) Restart(ctx context.Context, siteID string) error {
	s.mu.RLock()
	spec, ok := s.specs[siteID]
	s.mu.RUnlock()
	if !ok {
		return ErrRuntimeNotFound
	}
	return s.Activate(ctx, spec)
}

func (s *Supervisor) Stop(ctx context.Context, siteID string) error {
	s.mu.Lock()
	rt := s.bySite[siteID]
	if rt == nil {
		s.mu.Unlock()
		return ErrRuntimeNotFound
	}
	delete(s.bySite, siteID)
	for _, host := range rt.spec.Hosts {
		delete(s.byHost, normalizeHost(host))
	}
	delete(s.specs, siteID)
	st := s.status[siteID]
	st.Status = StatusStopped
	s.status[siteID] = st
	s.mu.Unlock()
	return rt.Close(ctx)
}

func (s *Supervisor) GetByHost(host string) (*SiteRuntime, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rt, ok := s.byHost[normalizeHost(host)]
	return rt, ok
}

func (s *Supervisor) Status(siteID string) (RuntimeStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.status[siteID]
	return st, ok
}

func (s *Supervisor) Summary() Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hosts := make([]string, 0, len(s.byHost))
	for host := range s.byHost {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	runtimes := make([]RuntimeStatus, 0, len(s.status))
	for _, st := range s.status {
		runtimes = append(runtimes, st)
	}
	sort.Slice(runtimes, func(i, j int) bool { return runtimes[i].SiteID < runtimes[j].SiteID })
	return Summary{ActiveSites: len(s.bySite), Hosts: hosts, Runtimes: runtimes}
}

func (s *Supervisor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt, ok := s.GetByHost(r.Host)
	if !ok {
		http.Error(w, "unknown go-go-host site", http.StatusNotFound)
		return
	}
	rec := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	rt.ServeHTTP(rec, r)
	s.recordRequest(rt.spec.SiteID, rec.statusCode)
}

var ErrRuntimeNotFound = errors.New("runtime not found")

func normalizeHosts(hosts []string) []string {
	out := make([]string, 0, len(hosts))
	seen := map[string]struct{}{}
	for _, host := range hosts {
		n := normalizeHost(host)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	host = strings.TrimSuffix(host, ".")
	if i := strings.LastIndex(host, ":"); i >= 0 && strings.Count(host, ":") == 1 {
		host = host[:i]
	}
	return host
}

func (s *Supervisor) setStatus(siteID string, status RuntimeStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	previous := s.status[siteID]
	status.RequestsTotal = previous.RequestsTotal
	status.ErrorsTotal = previous.ErrorsTotal
	s.status[siteID] = status
}

func (s *Supervisor) recordRequest(siteID string, statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.status[siteID]
	st.RequestsTotal++
	if statusCode >= 500 {
		st.ErrorsTotal++
	}
	s.status[siteID] = st
}
