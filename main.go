package main

import (
	"cache_server/HTTP"
	"cache_server/TCP"
	"cache_server/cache"
)

func main()  {
	ca := cache.New("inmemory")
	go TCP.New(ca).Listen()
	HTTP.New(ca).Listen()
}
