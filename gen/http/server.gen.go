// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gen

import (
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Change account password
	// (PATCH /api/v1/account/password)
	ChangePassword(c *gin.Context)
	// User login
	// (POST /api/v1/login)
	LoginAccount(c *gin.Context)
	// Logout user
	// (POST /api/v1/logout)
	LogoutAccount(c *gin.Context)
	// Create a new account
	// (POST /api/v1/register)
	CreateAccount(c *gin.Context)
	// User login
	// (POST /api/v2/login)
	LoginWithJWT(c *gin.Context)
	// Logout Account
	// (POST /api/v2/logout)
	LogoutJWT(c *gin.Context)
	// Refresh tokens
	// (POST /api/v2/token/refresh)
	RefreshToken(c *gin.Context)
	// Liveness probe endpoint
	// (GET /livez)
	CheckLiveness(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ChangePassword operation middleware
func (siw *ServerInterfaceWrapper) ChangePassword(c *gin.Context) {

	c.Set(CookieAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ChangePassword(c)
}

// LoginAccount operation middleware
func (siw *ServerInterfaceWrapper) LoginAccount(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.LoginAccount(c)
}

// LogoutAccount operation middleware
func (siw *ServerInterfaceWrapper) LogoutAccount(c *gin.Context) {

	c.Set(CookieAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.LogoutAccount(c)
}

// CreateAccount operation middleware
func (siw *ServerInterfaceWrapper) CreateAccount(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateAccount(c)
}

// LoginWithJWT operation middleware
func (siw *ServerInterfaceWrapper) LoginWithJWT(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.LoginWithJWT(c)
}

// LogoutJWT operation middleware
func (siw *ServerInterfaceWrapper) LogoutJWT(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.LogoutJWT(c)
}

// RefreshToken operation middleware
func (siw *ServerInterfaceWrapper) RefreshToken(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RefreshToken(c)
}

// CheckLiveness operation middleware
func (siw *ServerInterfaceWrapper) CheckLiveness(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CheckLiveness(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.PATCH(options.BaseURL+"/api/v1/account/password", wrapper.ChangePassword)
	router.POST(options.BaseURL+"/api/v1/login", wrapper.LoginAccount)
	router.POST(options.BaseURL+"/api/v1/logout", wrapper.LogoutAccount)
	router.POST(options.BaseURL+"/api/v1/register", wrapper.CreateAccount)
	router.POST(options.BaseURL+"/api/v2/login", wrapper.LoginWithJWT)
	router.POST(options.BaseURL+"/api/v2/logout", wrapper.LogoutJWT)
	router.POST(options.BaseURL+"/api/v2/token/refresh", wrapper.RefreshToken)
	router.GET(options.BaseURL+"/livez", wrapper.CheckLiveness)
}
