package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func zapLogger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		log := log.With(
			"remote", c.ClientIP(),
			"host", c.Request.Host,
			"method", c.Request.Method,
			"uri", c.Request.URL.String(),
			"duration", duration,
			"status", status,
		)
		if err := c.Errors.String(); err != "" {
			log = log.With("err", err)
		}
		switch {
		case status >= 500 || status == 0:
			log.Error()
		case status >= 400:
			log.Warn()
		default:
			log.Info()
		}

	}
}
