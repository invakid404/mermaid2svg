package main

import (
	_ "embed"
	"fmt"
	"github.com/invakid404/mermaid2svg/webdriver"
	"log"
	"sync"
)

const (
	example = `pie title NETFLIX
         "Time spent looking for movie" : 90
         "Time spent watching it" : 10`
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

	wg := sync.WaitGroup{}
	wg.Add(3)

	for i := 0; i < 3; i++ {
		i := i
		go func() {
			defer wg.Done()

			svg, err := driver.Render(example)
			if err != nil {
				log.Println("failed to render:", err)
			}

			fmt.Println(i, svg)
		}()
	}

	wg.Wait()
}
