package main

import (
	"github.com/labstack/echo/v4"
	"main/router"
)

func main() {
	e := echo.New()

	// "/api" 경로 접두사를 가진 그룹 생성
	v1ClaimGroup := e.Group("/v1/claim")

	v1ClaimGroup.GET("/authentication", router.GetAuthenticationQRCode)
	v1ClaimGroup.POST("/:sessionId", router.CreateClaim)
	e.Logger.Fatal(e.Start(":3000"))
}
