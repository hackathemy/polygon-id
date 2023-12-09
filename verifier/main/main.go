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
	e.GET("/v1/votes/:claimType/:contractAddress", router.GetVerificationQRCode)

	e.Logger.Fatal(e.Start(":3001"))
}
