// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error *string `json:"error,omitempty"`
}

// Task defines model for Task.
type Task struct {
	ErrMsg  *string `json:"ErrMsg,omitempty"`
	ID      *string `json:"ID,omitempty"`
	LastErr *string `json:"LastErr,omitempty"`
	Payload *string `json:"Payload,omitempty"`
	Queue   *string `json:"Queue,omitempty"`
	Retries *int    `json:"Retries,omitempty"`
	Timeout *string `json:"Timeout,omitempty"`
	Type    *string `json:"Type,omitempty"`
}

// DefaultError defines model for DefaultError.
type DefaultError = ErrorResponse

// PostMatchlinksJSONBody defines parameters for PostMatchlinks.
type PostMatchlinksJSONBody struct {
	CrawlSiteID *string `json:"CrawlSiteID,omitempty"`
	Debug       *bool   `json:"Debug,omitempty"`
	MaxDepth    *int    `json:"MaxDepth,omitempty"`
	SearchTerms *string `json:"SearchTerms,omitempty"`
	URL         *string `json:"URL,omitempty"`
}

// PostMatchlinksJSONRequestBody defines body for PostMatchlinks for application/json ContentType.
type PostMatchlinksJSONRequestBody PostMatchlinksJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all matching tasks
	// (GET /matchlinks)
	GetMatchlinks(ctx echo.Context) error
	// Create a new matching task
	// (POST /matchlinks)
	PostMatchlinks(ctx echo.Context) error
	// Delete the matching task
	// (DELETE /matchlinks/{id})
	DeleteMatchlinksId(ctx echo.Context, id string) error
	// Get the details of a matching task
	// (GET /matchlinks/{id})
	GetMatchlinksId(ctx echo.Context, id string) error
	// Ping the server
	// (GET /ping)
	GetPing(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetMatchlinks converts echo context to params.
func (w *ServerInterfaceWrapper) GetMatchlinks(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMatchlinks(ctx)
	return err
}

// PostMatchlinks converts echo context to params.
func (w *ServerInterfaceWrapper) PostMatchlinks(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostMatchlinks(ctx)
	return err
}

// DeleteMatchlinksId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteMatchlinksId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteMatchlinksId(ctx, id)
	return err
}

// GetMatchlinksId converts echo context to params.
func (w *ServerInterfaceWrapper) GetMatchlinksId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMatchlinksId(ctx, id)
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

	router.GET(baseURL+"/matchlinks", wrapper.GetMatchlinks)
	router.POST(baseURL+"/matchlinks", wrapper.PostMatchlinks)
	router.DELETE(baseURL+"/matchlinks/:id", wrapper.DeleteMatchlinksId)
	router.GET(baseURL+"/matchlinks/:id", wrapper.GetMatchlinksId)
	router.GET(baseURL+"/ping", wrapper.GetPing)

}
