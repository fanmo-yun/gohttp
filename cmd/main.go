package main

import (
	"gohttp/banner"
	"gohttp/server"
)

func main() {
	banner.ShowBanner()
	server.ServerRun()
}
