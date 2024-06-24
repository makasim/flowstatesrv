package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/memdriver"
	"github.com/makasim/flowstatesrv/internal/api/enginehandlerv1alpha1"
	"github.com/makasim/flowstatesrv/internal/protogen/flowstate/v1alpha1/flowstatev1alpha1connect"
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
	d := memdriver.New()
	d.SetFlow("test", flowstate.FlowFunc(func(stateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {
		log.Println("test flow executed!")
		return flowstate.End(stateCtx), nil
	}))

	e, err := flowstate.NewEngine(d)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	mux := http.NewServeMux()

	mux.Handle(flowstatev1alpha1connect.NewEngineServiceHandler(enginehandlerv1alpha1.New(e)))

	srv := &http.Server{
		Addr:    `127.0.0.1:8080`,
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
