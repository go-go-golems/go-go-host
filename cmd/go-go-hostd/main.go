package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/httpapi"
	"github.com/go-go-golems/go-go-host/internal/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "go-go-hostd",
		Short: "Run the go-go-host daemon",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			return run(cmd.Context(), cfg)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "Path to YAML or JSON daemon config")
	return cmd
}

func run(ctx context.Context, cfg config.Config) error {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parse log level %q: %w", cfg.LogLevel, err)
	}
	zerolog.SetGlobalLevel(level)

	st, err := store.Open(ctx, cfg.ControlDBDSN)
	if err != nil {
		return err
	}
	defer st.Close()
	if err := st.ApplyMigrations(ctx); err != nil {
		return err
	}
	if err := st.ReconcileStaleRuntimeStatuses(ctx); err != nil {
		return err
	}

	core := control.NewCoreWithStore(cfg, st)
	if err := core.Deployments.RestoreActiveRuntimes(ctx); err != nil {
		return err
	}
	server := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      httpapi.NewHandler(core),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", cfg.ListenAddr).Msg("starting go-go-host daemon")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-runCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}
