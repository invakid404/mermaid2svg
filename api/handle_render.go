package api

import (
	"errors"
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
	Content string `json:"content"`
}

func (body *RenderRequest) Bind(*http.Request) error {
	if body.Content == "" {
		return errors.New("missing required `content` field")
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
	svg, err := driver.Render(body.Content)
	if err != nil {
		_ = render.Render(res, req, ErrInternalServerError(err))
		return
	}

	res.Header().Set("Content-Type", "image/svg+xml")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(svg))
}
