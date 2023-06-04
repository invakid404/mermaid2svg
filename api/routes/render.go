package routes

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/invakid404/mermaid2svg/api/utils"
	"github.com/invakid404/mermaid2svg/webdriver"
	"net/http"
)

type renderImpl struct{}

func registerRender(router chi.Router) {
	impl := &renderImpl{}

	router.Post("/v1/render", impl.render)
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

func (impl *renderImpl) render(res http.ResponseWriter, req *http.Request) {
	body := &RenderRequest{}
	if err := render.Bind(req, body); err != nil {
		_ = render.Render(res, req, utils.ErrInvalidRequest(err))
		return
	}

	app := req.Context().Value("driver").(*webdriver.Driver)
	svg, err := app.Render(body.Content)
	if err != nil {
		_ = render.Render(res, req, utils.ErrInternalServerError(err))
		return
	}

	res.Header().Set("Content-Type", "image/svg+xml")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(svg))
}
