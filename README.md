# polygon-id
해당 문서는 POLYGON-ID 사용방법에 관한 문서 입니다.

https://schema-builder.polygonid.me/builder
해당 사이트에 접속하여서 VC스키마를 정의 해줍니다.
<img width="1421" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/deff8577-6259-43bc-9594-3c1d680d2f0e">
<img width="1421" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/251be27d-c2dc-45f7-a08e-19859e27c512">


https://github.com/pjhnocegood/polygon-id/blob/main/issuer/router/claims.go
makeClaimData 함수와 사용된 구조체를 생성한 스키마에 맞게 변경 합니다.
```
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
```

VC를 생성하기 위해서는 사용자의 DID가 필요합니다.
사용자의 DID값을 얻기 위하여 http://127.0.0.1:3000/v1/claim/authentication를 호출하여 authentication QR코드를 만들어서 인증을 받습니다.

<img width="1000" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/4e03ecc6-4a39-458a-9266-54c563fa0ee9">

For Android: [Polygon ID on Google Play](https://play.google.com/store/apps/details?id=com.polygonid.wallet)

For iOS: [For iOS: Polygon ID on the App Store](https://apps.apple.com/us/app/polygon-id/id1629870183)

<img width="1295" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/3c44560c-92aa-4677-9e1d-a6337a63202d">

QRLINK에 있는 값으로 QRCODE를 생성 후 PoyglonID APP을 통하여 스캔 후 authentication을 진행 해줍니다.
![image](https://github.com/pjhnocegood/polygon-id/assets/36693435/6ceaec99-211f-47bd-9435-1908bbaae5db)


http://127.0.0.1:3000/v1/claim/{{sessionID}}위의 API를 통해서 받은 sessionID와 보유 토큰수를 입력하여 VC를 생성하는 QR코드 값을 받습니다.
<img width="994" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/c88968d0-cf11-4060-a962-46779da48aca">

위의 API를 통해 받은 JSON데이터로  QRCODE를 생성 후 PoyglonID APP을 통하여 스캔하면 POLYGON-ID APP에 VC가 생성 됩니다.
<img width="1154" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/666378cb-71d4-4381-83a6-abce33aa4621">

http://127.0.0.1:3003/v1/contract/deploy를 호출하여 생성한 VC를 검증할 스마트컨트랙트를 배포하고 ZK로 어떤 내용을 검증할지 정의합니다.


아래의 예제는 votesTheshhold를 넘는 투표를 받으면 해당 컨트랙트가 가지고 있는 tokenAddress에 해당하는 토큰을 builder에게 전달하는 컨트랙트 예제 입니다.
<img width="1000" alt="image" src="https://github.com/pjhnocegood/polygon-id/assets/36693435/3797d7f8-d7d6-406e-a229-d7b35fb1e60e">

https://github.com/pjhnocegood/polygon-id/blob/main/on-chain-verification/contracts/TokenTransferContract.sol














