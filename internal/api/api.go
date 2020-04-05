package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.unanet.io/devops/eve-bot/internal/api/middleware"
	"gitlab.unanet.io/devops/eve-bot/internal/evelogger"
	"gitlab.unanet.io/devops/eve-bot/internal/servicefactory"
	"go.uber.org/zap"
)

// App is the interface for the Application
type App interface {
	Serve() error
}

// app is the main server that serves up the API
type app struct {
	done       chan bool
	sigChannel chan os.Signal
	healthy    int32
	router     *mux.Router
	server     *http.Server
	logger     evelogger.Container
	svcFactory *servicefactory.Container
}

// New creates a new Application server (our Rest API)
// this is where the mux router is created and decorated with middleware
func New(svcFactory *servicefactory.Container) App {
	router := mux.NewRouter().StrictSlash(true)

	a := &app{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", svcFactory.Config.API.Port),
			Handler: handlers.CORS(
				handlers.AllowedOrigins(svcFactory.Config.API.AllowedOrigins),
				handlers.AllowedHeaders(svcFactory.Config.API.AllowedHeaders),
				handlers.AllowedMethods(svcFactory.Config.API.AllowedMethods),
			)(middleware.Adapt(
				router,
				middleware.LogMetricsHandler(svcFactory.Logger, svcFactory.Metrics),
				middleware.TimeoutHandler(svcFactory.Config.API.TimeoutSecs))),
			ReadTimeout:  time.Duration(svcFactory.Config.API.ReadTimeOutSecs) * time.Second,
			WriteTimeout: time.Duration(svcFactory.Config.API.WriteTimeOutSecs) * time.Second,
			IdleTimeout:  time.Duration(svcFactory.Config.API.IdleTimeOutSecs) * time.Second,
		},
		done:       make(chan bool),
		logger:     svcFactory.Logger.With(zap.String("package", "api")),
		sigChannel: make(chan os.Signal, 1024),
		healthy:    0,
		router:     router,
		svcFactory: svcFactory,
	}

	a.registerRoutes(svcFactory)
	return a
}

// Handle SIGNALS
func (a *app) sigHandler() {
	for {
		sig := <-a.sigChannel
		switch sig {
		case syscall.SIGHUP:
			a.logger.Bg().Info("reload config not setup up")
		case os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
			a.logger.Bg().Info("caught shutdown signal", zap.String("signal", sig.String()))
			a.gracefulShutdown()
		}
	}
}

func (a *app) gracefulShutdown() {
	// Store atomic "health" value
	// Prevents new requests from coming in while draining
	atomic.StoreInt32(&a.healthy, 0)

	// Pause the Context for `ShutdownTimeoutSecs` config value
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.svcFactory.Config.API.ShutdownTimeoutSecs)*time.Second)
	defer cancel()

	// Turn off keepalive
	a.server.SetKeepAlivesEnabled(false)

	// Attempt to shutdown cleanly
	if err := a.server.Shutdown(ctx); err != nil {
		// YIKES! Shutdown failed, time to panic
		panic("http server failed graceful shutdown")
	}
	close(a.done)
}

// Serve is the entrypoint into the api
func (a *app) Serve() error {
	// Serve Up Metrics on a separte port
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%v", a.svcFactory.Config.API.MetricsPort), nil)

	// signal the token-svc channel whenever an OS.Interrupt or SIGHUP occur
	// (both currently terminate. would like to use the SIGHUP for config reload)
	signal.Notify(a.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// goroutine to handle signals
	go a.sigHandler()

	// atomically store the health as "healthy=1"
	atomic.StoreInt32(&a.healthy, 1)

	// log server start details
	a.logger.Bg().Debug("server up ===========>",
		zap.String("port", a.svcFactory.Config.API.Port),
		zap.String("metrics_port", a.svcFactory.Config.API.MetricsPort),
		zap.String("full_version", a.svcFactory.VersionInfo.FullVersionNumber(true)),
		zap.String("build_date", a.svcFactory.VersionInfo.BuildDate),
		zap.String("metadata", a.svcFactory.VersionInfo.VersionMetadata),
		zap.String("prerelease", a.svcFactory.VersionInfo.VersionPrerelease),
		zap.String("version", a.svcFactory.VersionInfo.Version),
		zap.String("revision", a.svcFactory.VersionInfo.Revision),
		zap.String("author", a.svcFactory.VersionInfo.Author),
		zap.String("branch", a.svcFactory.VersionInfo.Branch),
		zap.String("builder", a.svcFactory.VersionInfo.BuildUser),
		zap.String("host", a.svcFactory.VersionInfo.BuildHost))

	// Tally Build Metrics
	a.svcFactory.Metrics.StatBuildInfo.WithLabelValues(
		a.svcFactory.Config.API.ServiceName,
		a.svcFactory.VersionInfo.Revision,
		a.svcFactory.VersionInfo.Branch,
		a.svcFactory.VersionInfo.Version,
		a.svcFactory.VersionInfo.Author,
		a.svcFactory.VersionInfo.BuildDate,
		a.svcFactory.VersionInfo.BuildUser,
		a.svcFactory.VersionInfo.BuildHost).Set(1)

	// serve up the api by listening on the configured port
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Bg().Error("failed ListenAndServe", zap.String("port", a.svcFactory.Config.API.Port), zap.Error(err))
		return err
	}

	// signal caught in a.gracefulShutdown() and close(a.done) called
	<-a.done

	// log server shutdown details
	a.logger.Bg().Info("graceful server shutdown  ===========>", zap.String("lifetime", a.svcFactory.VersionInfo.UpTime().String()))
	return nil
}
