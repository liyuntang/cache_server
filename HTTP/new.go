package HTTP

import (
	"cache_server/cache"
	"cache_server/cluster"
	"net/http"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server)Listen()  {
	http.Handle("/cache/",s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	http.ListenAndServe(":12345", nil)
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}