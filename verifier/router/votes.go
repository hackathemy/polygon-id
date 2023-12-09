package router

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetVerificationQRCode(c echo.Context) error {
	claimType := c.Param("claimType")
	contractAddress := c.Param("contractAddress")
	var response map[string]interface{}
	if claimType == "builder" {
		response = map[string]interface{}{
			"id":   "7f38a193-0918-4a48-9fac-36adfdb8b542",
			"typ":  "application/iden3comm-plain-json",
			"type": "https://iden3-communication.io/proofs/1.0/contract-invoke-request",
			"thid": "7f38a193-0918-4a48-9fac-36adfdb8b542",
			"body": map[string]interface{}{
				"reason": "airdrop participation",
				"transaction_data": map[string]interface{}{
					"contract_address": contractAddress,
					"method_id":        "b68967e2",
					"chain_id":         80001,
					"network":          "polygon-mumbai",
				},
				"scope": []map[string]interface{}{
					{
						"id":        1702148496,
						"circuitId": "credentialAtomicQuerySigV2OnChain",
						"query": map[string]interface{}{
							"allowedIssuers": []string{
								"did:polygonid:polygon:mumbai:2qNzSKEuYnHwN7NgdmVM8DMYgpWVnCtnup1esfeCJ1",
							},
							"context": "ipfs://QmWamN8tGpJjKdr1NDEduedP16CH85W4Ua8oBi4NE3MAqA",
							"credentialSubject": map[string]interface{}{
								"isBuilder": map[string]interface{}{
									"$eq": 1,
								},
							},
							"type": "ReFreshBuilder",
						},
					},
				},
			},
		}

	} else if claimType == "funder" {
		response = map[string]interface{}{
			"id":   "7f38a193-0918-4a48-9fac-36adfdb8b542",
			"typ":  "application/iden3comm-plain-json",
			"type": "https://iden3-communication.io/proofs/1.0/contract-invoke-request",
			"thid": "7f38a193-0918-4a48-9fac-36adfdb8b542",
			"body": map[string]interface{}{
				"reason": "airdrop participation",
				"transaction_data": map[string]interface{}{
					"contract_address": contractAddress,
					"method_id":        "b68967e2",
					"chain_id":         80001,
					"network":          "polygon-mumbai",
				},
				"scope": []map[string]interface{}{
					{
						"id":        1701840378,
						"circuitId": "credentialAtomicQuerySigV2OnChain",
						"query": map[string]interface{}{
							"allowedIssuers": []string{
								"did:polygonid:polygon:mumbai:2qNzSKEuYnHwN7NgdmVM8DMYgpWVnCtnup1esfeCJ1",
							},
							"context": "ipfs://QmQb3pfSfmFZNQapcXk3zdnDnmmpqmiZ6YWFcCwyq14ajM",
							"credentialSubject": map[string]interface{}{
								"token": map[string]interface{}{
									"$gt": 0,
								},
							},
							"type": "DAOVerificationToken",
						},
					},
				},
			},
		}

	}
	// JSON 응답 반환
	return c.JSON(http.StatusOK, response)

}
