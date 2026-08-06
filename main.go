package main

import "learn-gin/internal/app"

func main() {
	application := app.New()
	application.Run()
}
