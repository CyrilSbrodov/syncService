package main

import "github.com/CyrilSbrodov/syncService/internal/app"

func main() {
	srv := app.NewServerApp()
	srv.Run()
}
