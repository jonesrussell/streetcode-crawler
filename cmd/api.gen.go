// Package cmd provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package cmd

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// PostArticlesStartJSONBody defines parameters for PostArticlesStart.
type PostArticlesStartJSONBody struct {
	CrawlSiteID *string `json:"CrawlSiteID,omitempty"`
	Debug       *bool   `json:"Debug,omitempty"`
	MaxDepth    *int    `json:"MaxDepth,omitempty"`
	SearchTerms *string `json:"SearchTerms,omitempty"`
	URL         *string `json:"URL,omitempty"`
}

// PostCrawlingStartJSONBody defines parameters for PostCrawlingStart.
type PostCrawlingStartJSONBody struct {
	CrawlSiteID *string `json:"CrawlSiteID,omitempty"`
	Debug       *bool   `json:"Debug,omitempty"`
	MaxDepth    *int    `json:"MaxDepth,omitempty"`
	SearchTerms *string `json:"SearchTerms,omitempty"`
	URL         *string `json:"URL,omitempty"`
}

// PostArticlesStartJSONRequestBody defines body for PostArticlesStart for application/json ContentType.
type PostArticlesStartJSONRequestBody PostArticlesStartJSONBody

// PostCrawlingStartJSONRequestBody defines body for PostCrawlingStart for application/json ContentType.
type PostCrawlingStartJSONRequestBody PostCrawlingStartJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Retrieve the status and results of a match job
	// (GET /articles/info/{id})
	GetArticlesInfoId(ctx echo.Context, id string) error
	// Retrieve an index of match jobs
	// (GET /articles/jobs)
	GetArticlesJobs(ctx echo.Context) error
	// Start the matching process
	// (POST /articles/start)
	PostArticlesStart(ctx echo.Context) error
	// Stop the matching process
	// (POST /articles/stop/{id})
	PostArticlesStopId(ctx echo.Context, id string) error
	// Retrieve the status and results of a crawl job
	// (GET /crawling/info/{id})
	GetCrawlingInfoId(ctx echo.Context, id string) error
	// Retrieve an index of crawl jobs
	// (GET /crawling/jobs)
	GetCrawlingJobs(ctx echo.Context) error
	// Start the crawling process
	// (POST /crawling/start)
	PostCrawlingStart(ctx echo.Context) error
	// Stop the crawling process
	// (POST /crawling/stop/{id})
	PostCrawlingStopId(ctx echo.Context, id string) error
	// Ping the server
	// (GET /ping)
	GetPing(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetArticlesInfoId converts echo context to params.
func (w *ServerInterfaceWrapper) GetArticlesInfoId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetArticlesInfoId(ctx, id)
	return err
}

// GetArticlesJobs converts echo context to params.
func (w *ServerInterfaceWrapper) GetArticlesJobs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetArticlesJobs(ctx)
	return err
}

// PostArticlesStart converts echo context to params.
func (w *ServerInterfaceWrapper) PostArticlesStart(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostArticlesStart(ctx)
	return err
}

// PostArticlesStopId converts echo context to params.
func (w *ServerInterfaceWrapper) PostArticlesStopId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostArticlesStopId(ctx, id)
	return err
}

// GetCrawlingInfoId converts echo context to params.
func (w *ServerInterfaceWrapper) GetCrawlingInfoId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCrawlingInfoId(ctx, id)
	return err
}

// GetCrawlingJobs converts echo context to params.
func (w *ServerInterfaceWrapper) GetCrawlingJobs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCrawlingJobs(ctx)
	return err
}

// PostCrawlingStart converts echo context to params.
func (w *ServerInterfaceWrapper) PostCrawlingStart(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostCrawlingStart(ctx)
	return err
}

// PostCrawlingStopId converts echo context to params.
func (w *ServerInterfaceWrapper) PostCrawlingStopId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostCrawlingStopId(ctx, id)
	return err
}

// GetPing converts echo context to params.
func (w *ServerInterfaceWrapper) GetPing(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPing(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/articles/info/:id", wrapper.GetArticlesInfoId)
	router.GET(baseURL+"/articles/jobs", wrapper.GetArticlesJobs)
	router.POST(baseURL+"/articles/start", wrapper.PostArticlesStart)
	router.POST(baseURL+"/articles/stop/:id", wrapper.PostArticlesStopId)
	router.GET(baseURL+"/crawling/info/:id", wrapper.GetCrawlingInfoId)
	router.GET(baseURL+"/crawling/jobs", wrapper.GetCrawlingJobs)
	router.POST(baseURL+"/crawling/start", wrapper.PostCrawlingStart)
	router.POST(baseURL+"/crawling/stop/:id", wrapper.PostCrawlingStopId)
	router.GET(baseURL+"/ping", wrapper.GetPing)

}
