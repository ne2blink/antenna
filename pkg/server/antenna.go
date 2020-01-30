package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) broadcast(c *gin.Context) {
	id := c.Param("id")

	var parseMode string
	switch c.ContentType() {
	case "text/markdown":
		parseMode = "Markdown"
	case "text/html":
		parseMode = "HTML"
	}

	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.antenna.Broadcast(id, string(data), parseMode); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
