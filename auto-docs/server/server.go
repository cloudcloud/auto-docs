// Package server provides capabilities to manage and process
// incoming HTTP requests. As this is the main driver for doc
// surfacing and interaction, special attention is paid to
// openly handling generated paths.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	autodocs "github.com/cloudcloud/auto-docs"
	"github.com/cloudcloud/auto-docs/auto-docs/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server holds details around the executing HTTP server based
// on the provided configuration.
type Server struct {
	// Config contains the configuration data that helped start
	// this instance of the Server.
	Config *autodocs.Config

	// Engine holds the gin instance for handling HTTP.
	Engine *gin.Engine

	// Listen is the directive upon which the server will accept
	// connections and process requests.
	Listen string

	// Origins contains all CORS allowed origins.
	Origins []string
}

// New is a short-hand to give a functional method for
// generating a new server.
func New(c *autodocs.Config) *Server {
	gin.SetMode(gin.ReleaseMode)

	return &Server{
		Config:  c,
		Engine:  gin.Default(),
		Listen:  c.Listen,
		Origins: []string{"*"},
	}
}

// Start will finish preparing and begin serving HTTP.
func (s *Server) Start() {
	d := data.Prep(s.Config.Git)
	t := time.NewTicker(
		time.Duration(s.Config.Git.Period) * time.Second,
	)
	defer t.Stop()

	go func(tick *time.Ticker) {
		for {
			select {
			case t := <-tick.C:
				d.Fetch(t)
			}
		}
	}(t)

	// TODO: Allow for graceful server shutdown.

	s.addMiddleware().addAPI().addHelpers().Serve()
}

// addAPI will add the route handling for API methods.
func (s *Server) addAPI() *Server {
	api := s.Engine.Group("/_api")
	api.GET("pages", pages)
	api.GET("page/*path", page)

	return s
}

// addHelpers will add additional routes for internal working.
func (s *Server) addHelpers() *Server {
	s.Engine.GET("/_health", health)
	s.Engine.NoRoute(root)

	return s
}

// addMiddleware will setup our required middleware methods on
// the internal engine.
func (s *Server) addMiddleware() *Server {
	handleFiles(s.Engine)

	s.Engine.Use(
		cors.New(cors.Config{
			AllowOrigins: s.Origins,
			AllowMethods: []string{"GET", "POST", "PUT", "OPTIONS", "HEAD", "DELETE"},
			AllowHeaders: []string{"Origin", "X-Client", "Content-Type", "Connection"},
		}),
	)

	return s
}

// Serve will begin to HTTP server execution.
func (s *Server) Serve() {
	srv := &http.Server{
		Addr:    s.Listen,
		Handler: s.Engine,
	}

	// this will be s.Listen
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting.")
}
