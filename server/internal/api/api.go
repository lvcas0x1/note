// Package api registers the read-only note API with Huma on Gin.
package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"notes/internal/content"
	"notes/internal/search"
	"notes/internal/store"
)

// Version is the API version reported in the OpenAPI document.
const Version = "1.0.0"

// MaxCalendarDays is the longest calendar range one request may ask for.
const MaxCalendarDays = 366

// Builds gives handlers the build that is currently served.
type Builds interface {
	// Reader returns the current build, or store.ErrNoBuild.
	Reader() (*store.Reader, error)
}

func init() {
	// Every handler returns non-nil slices, so arrays are never null.
	huma.DefaultArrayNullable = false
	// Huma reports request validation failures as 422; the note API spec uses 400.
	base := huma.NewError
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		if status == http.StatusUnprocessableEntity {
			status = http.StatusBadRequest
		}
		return base(status, msg, errs...)
	}
}

// Config returns the Huma configuration: OpenAPI 3.1 at /api/v1/openapi.json,
// no docs UI (it would load third-party scripts), and no $schema links.
func Config() huma.Config {
	cfg := huma.DefaultConfig("Notes API", Version)
	cfg.Info.Description = "Read-only API for the note site. No authentication."
	cfg.OpenAPIPath = "/api/v1/openapi"
	cfg.DocsPath = ""
	cfg.SchemasPath = ""
	cfg.CreateHooks = nil
	return cfg
}

// Register adds the note API to the Gin engine and returns the Huma API.
func Register(engine *gin.Engine, builds Builds, mw ...func(huma.Context, func(huma.Context))) huma.API {
	a := humagin.New(engine, Config())
	h := &handlers{builds: builds}
	h.register(a, mw)
	return a
}

// NewSpecAPI returns a Huma API on a throwaway engine, for printing the OpenAPI document.
func NewSpecAPI() huma.API {
	gin.SetMode(gin.ReleaseMode)
	return Register(gin.New(), nil)
}

type handlers struct{ builds Builds }

// ---- inputs and outputs

// SidebarOutput is the sidebar tree.
type SidebarOutput struct {
	Body []content.Node
}

// ArticleInput selects one article.
type ArticleInput struct {
	ID string `path:"id" maxLength:"512" doc:"Article ID: the content-relative path without .md, URL-encoded (slashes as %2F)."`
}

// ArticleOutput is one article.
type ArticleOutput struct {
	Body *store.ArticleDetail
}

// CalendarInput is an inclusive date range.
type CalendarInput struct {
	From string `query:"from" required:"true" format:"date" doc:"First day, inclusive (YYYY-MM-DD)."`
	To   string `query:"to" required:"true" format:"date" doc:"Last day, inclusive (YYYY-MM-DD). At most 366 days after from."`
}

// CalendarOutput lists the days that have articles.
type CalendarOutput struct {
	Body []store.DayCount
}

// DayInput selects one day.
type DayInput struct {
	Date string `path:"date" format:"date" doc:"Day (YYYY-MM-DD)."`
}

// DayOutput lists one day's articles.
type DayOutput struct {
	Body []store.ArticleSummary
}

// SearchInput is a search query.
type SearchInput struct {
	Q     string `query:"q" required:"true" minLength:"3" maxLength:"200" doc:"Search text. Terms are matched as literal substrings; terms shorter than 3 characters are ignored."`
	Limit int    `query:"limit" minimum:"1" maximum:"20" default:"20"`
}

// SearchOutput lists matching articles, best first.
type SearchOutput struct {
	Body []search.Hit
}

func (h *handlers) register(a huma.API, mw []func(huma.Context, func(huma.Context))) {
	tags := []string{"notes"}
	huma.Register(a, huma.Operation{OperationID: "getSidebar", Method: http.MethodGet, Path: "/api/v1/sidebar",
		Summary: "Get the sidebar tree", Tags: tags}, h.sidebar)
	huma.Register(a, huma.Operation{OperationID: "getArticle", Method: http.MethodGet, Path: "/api/v1/articles/{id}",
		Summary: "Get one article", Tags: tags, Errors: []int{http.StatusNotFound}}, h.article)
	huma.Register(a, huma.Operation{OperationID: "getCalendar", Method: http.MethodGet, Path: "/api/v1/calendar",
		Summary: "Count articles per day", Description: "Returns only days that have at least one article.",
		Tags: tags, Errors: []int{http.StatusBadRequest}}, h.calendar)
	huma.Register(a, huma.Operation{OperationID: "getDay", Method: http.MethodGet, Path: "/api/v1/days/{date}",
		Summary: "List the articles of one day", Description: "Sorted by title.",
		Tags: tags, Errors: []int{http.StatusBadRequest}}, h.day)
	huma.Register(a, huma.Operation{OperationID: "searchArticles", Method: http.MethodGet, Path: "/api/v1/search",
		Summary: "Search articles", Tags: tags, Middlewares: mw,
		Errors: []int{http.StatusBadRequest, http.StatusTooManyRequests}}, h.search)
}

func (h *handlers) reader() (*store.Reader, error) {
	r, err := h.builds.Reader()
	if errors.Is(err, store.ErrNoBuild) {
		return nil, huma.Error503ServiceUnavailable("no content has been built yet")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("cannot open the current build")
	}
	return r, nil
}

func (h *handlers) sidebar(ctx context.Context, _ *struct{}) (*SidebarOutput, error) {
	r, err := h.reader()
	if err != nil {
		return nil, err
	}
	nodes, err := r.Sidebar(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("cannot read the sidebar")
	}
	return &SidebarOutput{Body: nodes}, nil
}

func (h *handlers) article(ctx context.Context, in *ArticleInput) (*ArticleOutput, error) {
	r, err := h.reader()
	if err != nil {
		return nil, err
	}
	a, err := r.Article(ctx, in.ID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, huma.Error404NotFound("article " + in.ID + " does not exist")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("cannot read the article")
	}
	return &ArticleOutput{Body: a}, nil
}

func (h *handlers) calendar(ctx context.Context, in *CalendarInput) (*CalendarOutput, error) {
	from, _ := time.Parse(content.DateLayout, in.From)
	to, _ := time.Parse(content.DateLayout, in.To)
	if to.Before(from) {
		return nil, huma.Error400BadRequest("from must not be after to",
			&huma.ErrorDetail{Location: "query.from", Message: "must not be after to", Value: in.From})
	}
	if to.Sub(from) > MaxCalendarDays*24*time.Hour {
		return nil, huma.Error400BadRequest("the range may span at most 366 days",
			&huma.ErrorDetail{Location: "query.to", Message: "must be at most 366 days after from", Value: in.To})
	}
	r, err := h.reader()
	if err != nil {
		return nil, err
	}
	days, err := r.Calendar(ctx, in.From, in.To)
	if err != nil {
		return nil, huma.Error500InternalServerError("cannot read the calendar")
	}
	return &CalendarOutput{Body: days}, nil
}

func (h *handlers) day(ctx context.Context, in *DayInput) (*DayOutput, error) {
	r, err := h.reader()
	if err != nil {
		return nil, err
	}
	list, err := r.Day(ctx, in.Date)
	if err != nil {
		return nil, huma.Error500InternalServerError("cannot read the day")
	}
	return &DayOutput{Body: list}, nil
}

func (h *handlers) search(ctx context.Context, in *SearchInput) (*SearchOutput, error) {
	r, err := h.reader()
	if err != nil {
		return nil, err
	}
	hits, err := search.Search(ctx, r.DB(), in.Q, in.Limit)
	switch {
	case errors.Is(err, search.ErrTooShort), errors.Is(err, search.ErrTooLong):
		return nil, huma.Error400BadRequest(err.Error(), &huma.ErrorDetail{Location: "query.q", Message: err.Error(), Value: in.Q})
	case err != nil:
		return nil, huma.Error500InternalServerError("search failed")
	}
	return &SearchOutput{Body: hits}, nil
}
