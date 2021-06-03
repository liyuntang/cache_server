package TCP

import (
	"cache_server/cache"
	"cache_server/cluster"
	"net"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server)Listen()  {
	l, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}

	for {
		//fmt.Println("accept....................")
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go s.process(c)
	}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}





















