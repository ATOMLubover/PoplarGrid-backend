// package main

// import (
// 	_ "poplargrid/docs" // Replace with your project path, ensure 'swag init' has been run to generate docs

// 	"github.com/iris-contrib/swagger/swaggerFiles"
// 	swagger "github.com/iris-contrib/swagger/v12"
// 	"github.com/kataras/iris/v12"
// )

// // HealthCheckResponse Health check response struct
// type HealthCheckResponse struct {
// 	Status  string `json:"status" example:"OK"`
// 	Version string `json:"version" example:"1.0.0"`
// }

// // GreetRequest Example request DTO
// type GreetRequest struct {
// 	Name     string `json:"name" binding:"required" example:"John"`
// 	Language string `json:"language" example:"en" enums:"en,es,fr"`
// }

// // GreetResponse Example response DTO
// type GreetResponse struct {
// 	Message string `json:"message" example:"Hello John!"`
// 	Status  int    `json:"status" example:"200"`
// }

// // @title Iris Swagger API
// // @version 1.0
// // @description A complete runnable Swagger integration example
// // @host localhost:8080
// // @BasePath /api/v1
// func main() {
// 	app := iris.Default()
// 	app.Logger().SetLevel("debug")

// 	// 1. Register Swagger UI and documentation routes
// 	// Explicitly handle /swagger without a trailing slash, redirecting it to /swagger/.
// 	// Use 302 Found (temporary redirect) to prevent browser caching issues.
// 	app.Get("/swagger", func(ctx iris.Context) {
// 		ctx.Redirect("/swagger/", iris.StatusFound)
// 	})

// 	// Core: Use swagger.WrapHandler to handle all requests related to Swagger UI.
// 	// It is responsible for serving index.html and other static files, and fetching doc.json based on configuration.
// 	// When accessing /swagger/, it should return index.html.
// 	// When accessing /swagger, it should theoretically be handled by the redirect above.
// 	// {any:path} will match /swagger/, /swagger/index.html, /swagger/swagger-ui.css, etc.
// 	app.Get("/swagger/{any:path}", swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
// 		// This is the URL Swagger UI uses to fetch the API definition.
// 		// Using a relative path can sometimes resolve environment-specific path resolution issues.
// 		c.URL = "/swagger/doc.json" // Modified to relative path
// 	}))

// 	// 2. Register API routes (to avoid path conflicts)
// 	api := app.Party("/api/v1")
// 	{
// 		api.Get("/health", HealthCheck)
// 		api.Post("/greet", GreetHandler)
// 	}

// 	app.Listen(":8080")
// }

// // HealthCheck godoc
// // @Summary Service health check
// // @Description Checks the service status
// // @Tags System
// // @Produce json
// // @Success 200 {object} HealthCheckResponse
// // @Router /health [get]
// func HealthCheck(ctx iris.Context) {
// 	ctx.JSON(HealthCheckResponse{
// 		Status:  "OK",
// 		Version: "1.0.0",
// 	})
// }

// // GreetHandler godoc
// // @Summary Generate greeting
// // @Description Generates a greeting based on name and language
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param body body GreetRequest true "Request parameters"
// // @Success 200 {object} GreetResponse
// // @Failure 400 {object} map[string]string
// // @Router /greet [post]
// func GreetHandler(ctx iris.Context) {
// 	var req GreetRequest
// 	if err := ctx.ReadJSON(&req); err != nil {
// 		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
// 			Title("Request parameter error").DetailErr(err))
// 		return
// 	}

// 	// Processing logic
// 	greeting := map[string]string{"en": "Hello", "es": "Hola", "fr": "Bonjour"}
// 	msg, ok := greeting[req.Language]
// 	if !ok {
// 		msg = greeting["en"] // Default language
// 	}

// 	ctx.JSON(GreetResponse{
// 		Message: msg + " " + req.Name + "!",
// 		Status:  iris.StatusOK,
// 	})
// }
