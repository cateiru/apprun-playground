package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルート
	e.GET("/", handleHealth)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func handleHealth(c echo.Context) error {
	response := HealthResponse{
		Status: "ok",
	}
	return c.JSON(http.StatusOK, response)
}
