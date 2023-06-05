package webdriver

import (
	_ "embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

const (
	geckodriverPort = 9559
)

var (
	//go:embed index.html
	mermaidHTML []byte

	//go:embed geckodriver.sh
	geckodriverEntrypoint []byte
)

type Driver struct {
	service           *selenium.Service
	serviceEntrypoint string
	webDriver         selenium.WebDriver
	server            *http.Server
	serverPort        int
	renderTasks       chan renderTask
	init              atomic.Bool
}

func New() (*Driver, error) {
	router := chi.NewRouter()
	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		_, err := res.Write(mermaidHTML)
		if err != nil {
			log.Println("failed to serve mermaid html:", err)
		}
	})

	server := &http.Server{
		Addr:    ":0",
		Handler: router,
	}

	return &Driver{
		service:     nil,
		server:      server,
		renderTasks: make(chan renderTask, 16),
	}, nil
}

func (driver *Driver) Start() error {
	if !driver.init.CompareAndSwap(false, true) {
		return nil
	}

	listener, err := net.Listen("tcp", driver.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", driver.server.Addr, err)
	}

	driver.serverPort = listener.Addr().(*net.TCPAddr).Port

	go func() {
		_ = driver.server.Serve(listener)
	}()

	entrypoint, err := os.CreateTemp("", "geckodriver")
	if err != nil {
		return fmt.Errorf("failed to create geckodriver entrypoint: %w", err)
	}
	driver.serviceEntrypoint = entrypoint.Name()

	if err = entrypoint.Chmod(0700); err != nil {
		return fmt.Errorf("failed to make geckodriver entrypoint executable: %w", err)
	}

	if _, err = entrypoint.Write(geckodriverEntrypoint); err != nil {
		return fmt.Errorf("failed to write geckodriver entrypoint: %w", err)
	}

	if err = entrypoint.Close(); err != nil {
		return fmt.Errorf("failed to close geckodriver entrypoint: %w", err)
	}

	service, err := selenium.NewGeckoDriverService(driver.serviceEntrypoint, geckodriverPort)
	if err != nil {
		return fmt.Errorf("failed to start geckodriver: %w", err)
	}
	driver.service = service

	capabilities := selenium.Capabilities{"browserName": "firefox"}
	capabilities.AddFirefox(firefox.Capabilities{
		Args: []string{"--headless"},
	})

	webDriver, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d", geckodriverPort))
	if err != nil {
		log.Fatalln("failed to create web driver:", err)
	}
	driver.webDriver = webDriver

	go driver.renderThread()

	return nil
}

func (driver *Driver) renderThread() {
	for task := range driver.renderTasks {
		driver.renderOne(task)
	}
}

func (driver *Driver) renderOne(task renderTask) {
	defer func() {
		close(task.output)
		close(task.err)
	}()

	if err := driver.webDriver.Get(fmt.Sprintf("http://localhost:%d/", driver.serverPort)); err != nil {
		task.err <- fmt.Errorf("failed to get internal page: %w", err)
		return
	}

	result, err := driver.webDriver.ExecuteScriptAsync(
		// language=javascript
		`
			const [source, options, callback] = arguments;
			mermaid.initialize(options);

			(async () => {
			  const { svg } = await mermaid.render("container", source);
			  document.documentElement.innerHTML = svg;

			  const container = document.querySelector("#container");
			  const bbox = container.getBBox();
			  container.setAttribute(
				"viewBox",
				[bbox.x, bbox.y, bbox.width, bbox.height].join(" ")
			  );

			  callback(container.outerHTML);
			})();
		`,
		[]any{task.input, task.options},
	)

	if err != nil {
		task.err <- fmt.Errorf("failed to render diagram: %w", err)
		return
	}

	task.output <- result.(string)
}

type renderTask struct {
	input   string
	options map[string]any
	output  chan<- string
	err     chan<- error
}

func (driver *Driver) enqueueRender(input string, options map[string]any) (<-chan string, <-chan error) {
	outputChan := make(chan string, 1)
	errChan := make(chan error, 1)

	driver.renderTasks <- renderTask{
		input:   input,
		options: options,
		output:  outputChan,
		err:     errChan,
	}

	return outputChan, errChan
}

func (driver *Driver) Render(input string, options map[string]any) (string, error) {
	outputChan, errChan := driver.enqueueRender(input, options)

	select {
	case result := <-outputChan:
		return result, nil
	case err := <-errChan:
		return "", err
	}
}

func (driver *Driver) Stop() error {
	close(driver.renderTasks)

	time.Sleep(time.Second * 2)

	if driver.webDriver != nil {
		if err := driver.webDriver.Quit(); err != nil {
			return fmt.Errorf("failed to quit webdriver web driver: %w", err)
		}
	}

	if driver.service != nil {
		if err := driver.service.Stop(); err != nil {
			return fmt.Errorf("failed to stop webdriver service: %w", err)
		}
	}

	if driver.server != nil {
		if err := driver.server.Close(); err != nil {
			return fmt.Errorf("failed to close http server: %w", err)
		}
	}

	if driver.serviceEntrypoint != "" {
		if err := os.Remove(driver.serviceEntrypoint); err != nil {
			return fmt.Errorf("failed to remove service entrypoint: %w", err)
		}
	}

	return nil
}
