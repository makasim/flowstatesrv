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

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/memdriver"
	"github.com/makasim/flowstatesrv/internal/api/corsmiddleware"
	"github.com/makasim/flowstatesrv/internal/api/serverservicev1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
	"github.com/makasim/flowstatesrv/ui"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
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
	d := memdriver.New()
	e, err := flowstate.NewEngine(d, a.l)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	addr := `127.0.0.1:8080`
	if os.Getenv(`FLOWSTATESRV_ADDR`) != `` {
		addr = os.Getenv(`FLOWSTATESRV_ADDR`)
	}

	corsMW := corsmiddleware.New(os.Getenv(`CORS_ENABLED`) == `true`)

	mux := http.NewServeMux()

	mux.Handle(corsMW.WrapPath(flowstatev1connect.NewServerServiceHandler(serverservicev1.New(e, d))))

	mux.Handle("/", corsMW.Wrap(http.FileServerFS(ui.PublicFS())))

	a.l.Info("http server starting", "addr", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
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

	if err := e.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("engine: shutdown: %w", err))
	}

	return shutdownRes
}
