package server

func (s *Server) EnableCompression(status bool) {
	s.enableCompression.Set(status)
	// s.createWSUpgrader()
}
