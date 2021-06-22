package main

import (
	"cache_server/HTTP"
	"cache_server/TCP"
	"cache_server/cache"
	"cache_server/cluster"
	"flag"
	_ "net/http/pprof"
)

func aaa()  {
	typ := flag.String("type", "inmemory", "cache type")
	ttl := flag.Int("ttl", 30, "cache time to live")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	c := cache.New(*typ, *ttl)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go TCP.New(c, n).Listen()
	HTTP.New(c, n).Listen()
}
