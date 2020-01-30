package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ne2blink/antenna/pkg/antenna"
	"github.com/ne2blink/antenna/pkg/storage"
	"go.uber.org/zap"
)

// Server listens for the content.
type Server struct {
	engine  *gin.Engine
	store   storage.AppStore
	antenna *antenna.Antenna
	log     *zap.SugaredLogger
}

// New creates a new server instance.
func New(store storage.AppStore, antenna *antenna.Antenna, log *zap.SugaredLogger) *Server {
	engine := gin.New()
	engine.Use(zapLogger(log), gin.Recovery())
	s := &Server{
		engine:  engine,
		store:   store,
		antenna: antenna,
		log:     log,
	}
	s.registerRoutes()
	return s
}

// Listen starts the server.
func (s *Server) Listen(addr string) error {
	defer s.log.Sync()
	s.log.Infof("starting at %s", addr)

	web := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Errorf("listen: %s", err.Error())
		return err
	}

	return nil
}
