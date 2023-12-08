package router

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

func GetVerificationQRCode(c echo.Context) error {
	// JSON 파일 읽기
	jsonData, err := ioutil.ReadFile("./qrValueProofRequest.json")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading JSON file")
	}

	// JSON 데이터를 클라이언트에게 반환
	return c.JSONBlob(http.StatusOK, jsonData)
}
