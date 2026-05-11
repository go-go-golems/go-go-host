package control

import (
	"github.com/go-go-golems/go-go-host/internal/config"
	hostruntime "github.com/go-go-golems/go-go-host/internal/runtime"
	"github.com/go-go-golems/go-go-host/internal/store"
)

// Core is the transport-agnostic application entrypoint. It owns product
// services and keeps HTTP/CLI code from reaching directly into persistence.
type Core struct {
	Config     config.Config
	Store      *store.Store
	Supervisor *hostruntime.Supervisor

	Orgs  *OrgService
	Sites *SiteService
}

func NewCore(cfg config.Config) *Core {
	return NewCoreWithStore(cfg, nil)
}

func NewCoreWithStore(cfg config.Config, st *store.Store) *Core {
	c := &Core{Config: cfg, Store: st, Supervisor: hostruntime.NewSupervisor(hostruntime.WithStatusRecorder(runtimeStatusRecorder{store: st}))}
	c.Orgs = &OrgService{store: st}
	c.Sites = &SiteService{store: st, baseDomain: cfg.BaseDomain}
	return c
}
