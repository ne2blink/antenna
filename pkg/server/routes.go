package server

func (s *Server) registerRoutes() {
	s.engine.POST("/antenna/:id", s.auth, s.broadcast)
}
