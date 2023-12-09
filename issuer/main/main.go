package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"main/router"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	// "/api" 경로 접두사를 가진 그룹 생성
	v1ClaimGroup := e.Group("/v1/claim")

	v1ClaimGroup.GET("/authentication", router.GetAuthenticationQRCode)
	v1ClaimGroup.POST("/:claimType/:sessionId", router.CreateClaim)
	e.Logger.Fatal(e.Start(":3000"))
}
