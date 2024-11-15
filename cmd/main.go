package main

import (
	"gohttp/banner"
	server "gohttp/internal"
)

func main() {
	banner.ShowBanner()
	server.ServerRun()
}
