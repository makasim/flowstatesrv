package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/badgerdriver"
	"github.com/makasim/flowstate/memdriver"
	"github.com/makasim/flowstate/netdriver"
	"github.com/makasim/flowstate/netflow"
	"github.com/makasim/flowstate/pgdriver"
	"github.com/makasim/flowstatesrv/ui"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type BadgerDriverConfig struct {
	InMemory bool
	Path     string
}

type PostgresDriverConfig struct {
	ConnString string
}

type Config struct {
	Driver         string
	BadgerDriver   BadgerDriverConfig
	PostgresDriver PostgresDriverConfig
}

type App struct {
	cfg Config
	l   *slog.Logger
}

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
		l:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (a *App) Run(ctx context.Context) error {
	var d flowstate.Driver
	switch a.cfg.Driver {
	case "memdriver":
		a.l.Info("init memdriver")
		d = memdriver.New(a.l)
	case "badgerdriver":
		a.l.Info("init badgerdriver")

		badgerCfg := badger.DefaultOptions(a.cfg.BadgerDriver.Path).
			WithInMemory(a.cfg.BadgerDriver.InMemory).
			WithLoggingLevel(2)
		db, err := badger.Open(badgerCfg)
		if err != nil {
			return fmt.Errorf("badger: open: %w", err)
		}
		defer db.Close()

		d0, err := badgerdriver.New(db)
		if err != nil {
			return fmt.Errorf("badgerdriver: new: %w", err)
		}
		defer d0.Shutdown(context.Background())

		d = d0
	case "pgdriver":
		a.l.Info("init pgdriver")
		conn, err := pgxpool.New(context.Background(), a.cfg.PostgresDriver.ConnString)
		if err != nil {
			return fmt.Errorf("pgxpool: new: %w", err)
		}
		defer conn.Close()

		d = pgdriver.New(conn, a.l)
	default:
		return fmt.Errorf("unknown driver: %s; support: memdriver, badgerdriver", a.cfg.Driver)
	}

	httpHost := `http://localhost:8080`
	if os.Getenv(`FLOWSTATESRV_HTTP_HOST`) != `` {
		httpHost = os.Getenv(`FLOWSTATESRV_HTTP_HOST`)
	}
	fr := netflow.NewRegistry(httpHost, d, a.l)
	defer fr.Close()

	e, err := flowstate.NewEngine(d, fr, a.l)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	r, err := flowstate.NewRecoverer(e, a.l)
	if err != nil {
		return fmt.Errorf("recoverer: new: %w", err)
	}

	dlr, err := flowstate.NewDelayer(e, a.l)
	if err != nil {
		return fmt.Errorf("delayer: new: %w", err)
	}

	addr := `0:8080`
	if os.Getenv(`FLOWSTATESRV_ADDR`) != `` {
		addr = os.Getenv(`FLOWSTATESRV_ADDR`)
	}

	uiH := http.FileServerFS(ui.PublicFS())

	a.l.Info("http server starting", "addr", addr)
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if handleCORS(rw, r) {
				return
			}
			if netdriver.HandleAll(rw, r, d, a.l) {
				return
			}
			if netflow.HandleExecute(rw, r, e) {
				return
			}

			uiH.ServeHTTP(rw, r)
		}), &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("WARN: http server: listen and serve: %s", err)
		}
	}()

	<-ctx.Done()

	var shutdownRes error
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer shutdownCtxCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("http server: shutdown: %w", err))
	}

	if err := r.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("recovery: shutdown: %w", err))
	}

	if err := dlr.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("delayer: shutdown: %w", err))
	}

	if err := e.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("engine: shutdown: %w", err))
	}

	return shutdownRes
}
