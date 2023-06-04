package main

import (
	_ "embed"
	"github.com/invakid404/mermaid2svg/api"
	"github.com/invakid404/mermaid2svg/webdriver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	driver, err := webdriver.New()
	if err != nil {
		log.Fatalln("failed to init webdriver: %w", err)
	}

	defer func(driver *webdriver.Driver) {
		if err := driver.Stop(); err != nil {
			log.Println("failed to close web driver:", err)
		}
	}(driver)

	if err = driver.Start(); err != nil {
		log.Fatalln("failed to start driver:", err)
	}

	app := api.New(api.Options{
		Driver: driver,
	})

	defer func(app *api.API) {
		err := app.Stop()
		if err != nil {
			log.Println("failed to close api:", err)
		}
	}(app)

	if err = app.Start(); err != nil {
		log.Fatalln("failed to start app:", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
