package api

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type renderAPI struct {
	api *API
}

func (r *renderAPI) Register(router chi.Router) {
	router.Post("/v1/render", r.render)
}

type RenderRequest struct {
	Content string         `json:"content"`
	Options map[string]any `json:"options"`
}

var (
	forbiddenOptions = []string{
		"maxTextSize",
		"securityLevel",
		"secure",
		"startOnLoad",
	}
)

func (body *RenderRequest) Bind(*http.Request) error {
	if body.Content == "" {
		return errors.New("missing required `content` field")
	}

	if body.Options == nil {
		body.Options = make(map[string]any)
	}

	for _, option := range forbiddenOptions {
		if _, ok := body.Options[option]; ok {
			return fmt.Errorf("usage of option %s is prohibited", option)
		}
	}

	return nil
}

func (r *renderAPI) render(res http.ResponseWriter, req *http.Request) {
	body := &RenderRequest{}
	if err := render.Bind(req, body); err != nil {
		_ = render.Render(res, req, ErrInvalidRequest(err))
		return
	}

	driver := r.api.driver
	svg, err := driver.Render(body.Content, body.Options)
	if err != nil {
		_ = render.Render(res, req, ErrInternalServerError(err))
		return
	}

	res.Header().Set("Content-Type", "image/svg+xml")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(svg))
}
