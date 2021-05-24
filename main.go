package main

import (
	"cache_server/HTTP"
	"cache_server/cache"
)

func main()  {
	c := cache.New("inmemory")
	HTTP.New(c).Listen()
}
