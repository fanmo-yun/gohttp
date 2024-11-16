package main

import (
	"gohttpd/banner"
	server "gohttpd/internal"
)

func main() {
	banner.ShowBanner()
	server.ServerRun()
}
