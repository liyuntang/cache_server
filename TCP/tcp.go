package TCP

import (
	"cache_server/cache"
	"net"
)

type Server struct {
	cache.Cache
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
		//defer func() {
		//	fmt.Println("关闭客户端连接")
		//	c.Close()
		//	fmt.Println("关闭客户端连接成功")
		//}()
		go s.process(c)
	}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}





















