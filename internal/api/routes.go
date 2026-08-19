package api

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/routes", s.handleRoutes)
	s.mux.HandleFunc("/api/routes/", s.handleRouteID)
	s.mux.HandleFunc("/api/try", s.handleTry)
	s.mux.HandleFunc("/api/proxy/", s.handleProxy)
	if s.opts.WebDir != "" {
		s.mux.Handle("/", staticHandler(s.opts.WebDir))
	}
}
