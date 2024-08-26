package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/internal/api/serverservicev1"
	"github.com/makasim/flowstatesrv/internal/driver"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
}

type App struct {
	cfg Config
}

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	d := driver.New()
	e, err := flowstate.NewEngine(d)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	addr := `127.0.0.1:8080`
	if os.Getenv(`FLOWSTATESRV_ADDR`) != `` {
		addr = os.Getenv(`FLOWSTATESRV_ADDR`)
	}

	mux := http.NewServeMux()

	mux.Handle(flowstatev1connect.NewServerServiceHandler(serverservicev1.New(e, d)))

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
