package main

import (
	"beyimtech-test/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	app, err := app.InitApp(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Start(ctx); err != nil {
		log.Fatalln(err)
	}
}
