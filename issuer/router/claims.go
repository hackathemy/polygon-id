package router

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

const (
	IssuerUIHost = "https://issuer-admin.polygonid.me"
	IssuerHost   = "https://self-hosted-platform.polygonid.me"
	IssuerDID    = "did:polygonid:polygon:mumbai:2qNzSKEuYnHwN7NgdmVM8DMYgpWVnCtnup1esfeCJ1"
)

type DAOVerificationToken struct {
	CredentialSchema  string `json:"credentialSchema"`
	Type              string `json:"type"`
	CredentialSubject struct {
		ID    string `json:"id"`
		Token int    `json:"token"`
	} `json:"credentialSubject"`
	SignatureProof bool `json:"signatureProof"`
	MtProof        bool `json:"mtProof"`
}

type ReFreshBuilder struct {
	CredentialSchema  string `json:"credentialSchema"`
	Type              string `json:"type"`
	CredentialSubject struct {
		ID        string `json:"id"`
		IsBuilder int    `json:"isBuilder"`
	} `json:"credentialSubject"`
	SignatureProof bool `json:"signatureProof"`
	MtProof        bool `json:"mtProof"`
}

type Connection struct {
	CreatedAt  string `json:"createdAt"`
	ID         string `json:"id"`
	IssuerID   string `json:"issuerID"`
	ModifiedAt string `json:"modifiedAt"`
	UserID     string `json:"userID"`
}

type CreateClaimRequest struct {
	TokenNumber int `json:"tokenNumber"`
}

func CreateClaim(c echo.Context) error {
	sessionId := c.Param("sessionId")
	claimType := c.Param("claimType")
	var createClaimRequest CreateClaimRequest
	if err := c.Bind(&createClaimRequest); err != nil {
		return err
	}
	tokenNumber := createClaimRequest.TokenNumber

	checkApiURL := fmt.Sprintf("%s/v1/authentication/sessions/%s", IssuerUIHost, sessionId)
	response, err := sendRequest(checkApiURL, "GET", nil, "issuer-ui")

	if response.StatusCode != http.StatusOK {
		return c.String(http.StatusNotFound, "connection does not exist")
	}

	// 응답 데이터 읽기
	responseData, err := ioutil.ReadAll(response.Body)
	responseDataMap := make(map[string]Connection)
	err = json.Unmarshal(responseData, &responseDataMap)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}

	userID := responseDataMap["connection"].UserID

	createApiURL := fmt.Sprintf("%s/v1/%s/claims", IssuerHost, IssuerDID)
	// JSON으로 변환
	var jsonData []byte
	if claimType == "builder" {
		jsonData, err = json.Marshal(makeClaimData(userID, tokenNumber))
	} else if claimType == "funder" {
		jsonData, err = json.Marshal(makeClaimData2(userID))
	} else {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error marshaling JSON")
	}

	response, err = sendRequest(createApiURL, "POST", jsonData, "issuer")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending request")
	}

	// 응답 데이터 읽기
	responseData, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}

	responseClaimDataMap := make(map[string]string)

	err = json.Unmarshal(responseData, &responseClaimDataMap)
	id := responseClaimDataMap["id"]

	createQRCODEApiURL := fmt.Sprintf("%s/v1/%s/claims/%s/qrcode", IssuerHost, IssuerDID, id)

	response, err = sendRequest(createQRCODEApiURL, "GET", nil, "issuer")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending request")
	}

	// 응답 데이터 읽기
	responseData, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}

	// 응답 데이터 출력
	return c.String(http.StatusOK, string(responseData))
}

func GetAuthenticationQRCode(c echo.Context) error {
	apiURL := fmt.Sprintf("%s/v1/authentication/qrcode", IssuerUIHost)
	// JSON으로 변환

	response, err := sendRequest(apiURL, "GET", nil, "issuer-ui")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending request")
	}

	// 응답 데이터 읽기
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}

	// 응답 데이터 출력
	return c.String(http.StatusOK, string(responseData))
}

func makeClaimData(userId string, tokenNumber int) DAOVerificationToken {
	credential := DAOVerificationToken{
		CredentialSchema: "ipfs://QmPX6sCNSzVDvrdz8fbaFfjckRqqXL4KPjpkhrLefVt4QN",
		Type:             "DAOVerificationToken",
		CredentialSubject: struct {
			ID    string `json:"id"`
			Token int    `json:"token"`
		}{
			ID:    userId,
			Token: tokenNumber,
		},
		SignatureProof: true,
		MtProof:        true,
	}
	return credential
}

func makeClaimData2(userId string) ReFreshBuilder {
	credential := ReFreshBuilder{
		CredentialSchema: "ipfs://QmTcVnJn1u8aM5NBYs3wBhtZhSAQqFtABM3tik765XzPgE",
		Type:             "ReFreshBuilder",
		CredentialSubject: struct {
			ID        string `json:"id"`
			IsBuilder int    `json:"isBuilder"`
		}{
			ID:        userId,
			IsBuilder: 1,
		},
		SignatureProof: true,
		MtProof:        true,
	}
	return credential
}

func sendRequest(url, method string, body []byte, destination string) (*http.Response, error) {
	client := &http.Client{}
	// 요청 생성
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// 요청 헤더 설정 (필요에 따라)
	req.Header.Set("Content-Type", "application/json")

	// Basic Auth 문자열 생성
	authString := ""

	if destination == "issuer" {
		username := "user-issuer"
		password := "password-issuer"
		// Basic Auth 문자열 생성
		authString = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	} else if destination == "issuer-ui" {
		authString = "dXNlci1hcGk6cGFzc3dvcmQtYXBp"
	}
	req.Header.Set("Authorization", "Basic "+authString)
	// Basic Auth 헤더 설정

	// 요청 보내기
	return client.Do(req)
}
