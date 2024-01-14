package oauth_svc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GOOGLE_UID_PREFIX   = "g"
	GOOGLE_CHANNEL_MARK = "google"
)

type GoogleOAuthValidateRequestBody struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type AccessTokenSuccess struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

type AccessTokenError struct {
	Err              string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AccessTokenResponse struct {
	AccessTokenSuccess
	AccessTokenError
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"` // 用户在Google的唯一标识码
	UId           string
	Canncel       string
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`
}

func (s *GoogleOAuthService) requestAccessTokenByProxy(proxyAddr string, bodyBytes []byte) (*http.Response, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		log.Println(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
	}
	if err != nil {
		log.Println("sequilize error: ", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", os.Getenv("OAUTH_GOOGLE_TRANSFER_TOKEN_ADDRESS"), bytes.NewBuffer(bodyBytes))

	if err != nil {
		log.Println("validate code error: ", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("get google token response error: ", err)
		return nil, err
	}
	return response, nil
}

func (s *GoogleOAuthService) acquireAccessToken(code string) (string, error) {
	body := GoogleOAuthValidateRequestBody{
		Code:         code,
		ClientId:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		RedirectUri:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URI"),
		GrantType:    "authorization_code",
	}
	bodyBytes, err := json.Marshal(body)
	proxyAddr := os.Getenv("PROXY_ADDRESS")

	var response *http.Response
	if proxyAddr == "" {
		http.Post(os.Getenv("OAUTH_GOOGLE_TRANSFER_TOKEN_ADDRESS"), "application/json", bytes.NewBuffer(bodyBytes))
	} else {
		response, _ = s.requestAccessTokenByProxy(proxyAddr, bodyBytes)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("read response body error: ", err)
		return "", err
	}

	data := &AccessTokenResponse{}
	json.Unmarshal(responseBody, data)
	if data.Err != "" {
		err := errors.New(fmt.Sprintf("validate access code error: code: %s, description: %s", data.Err, data.ErrorDescription))
		log.Println("read response body error: ", err)
		return "", err
	}
	log.Println("get access token success: ", data.IdToken)
	return data.IdToken, nil
}

func (s *GoogleOAuthService) transferIdTokenToUserInfo(idToken string) (GoogleUserInfo, error) {
	payload := strings.Split(idToken, ".")[1]
	idTokenJsonString, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		log.Println("decode id token error: ", err)
		return GoogleUserInfo{}, err
	}
	var info GoogleUserInfo
	json.Unmarshal(idTokenJsonString, &info)
	return info, nil
}

func (s *GoogleOAuthService) enforceUserInfo(gUserInfo *GoogleUserInfo) {
	gUserInfo.UId = fmt.Sprintf("%s_%s", GOOGLE_UID_PREFIX, gUserInfo.Sub)
	gUserInfo.Canncel = GOOGLE_CHANNEL_MARK
}

func (s *GoogleOAuthService) GetGoogleUserInfoByAccessCode(code string) (GoogleUserInfo, error) {
	// token, err := s.acquireAccessToken(code)
	// if err != nil {
	// 	return GoogleUserInfo{}, err
	// }
	gUserInfo, err := s.transferIdTokenToUserInfo("eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxNDEzY2Y0ZmEwY2I5MmEzYzNmNWEwNTQ1MDkxMzJjNDc2NjA5MzciLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI5MTU3OTQyMzYyNzQtamF2OW1zYzR0azZjdDE4aWZram05YXNwbmwyM2w2cXIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI5MTU3OTQyMzYyNzQtamF2OW1zYzR0azZjdDE4aWZram05YXNwbmwyM2w2cXIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTQ3NjU3NjgzNzkzMTI4MTIxMTMiLCJlbWFpbCI6InNkZmRmZGFkc2ZAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJCd1NULUFPeVA2cTkwakJybDhnaHd3IiwibmFtZSI6ImhlcmluIG11IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0tXcENJRDNNZzFtanlRN3ZnT2dkYmNKMEh4dnA5RjBrQ3BLNHhPMHpHdXlpVT1zOTYtYyIsImdpdmVuX25hbWUiOiJoZXJpbiIsImZhbWlseV9uYW1lIjoibXUiLCJsb2NhbGUiOiJ6aC1DTiIsImlhdCI6MTcwNDI5NDk0NCwiZXhwIjoxNzA0Mjk4NTQ0fQ.Ji_FwwOclCjxKgZVN6OQ2hVoWoIswgJQyNABAsWXLC68lmDs8Uts6OK7DiIFNByvgoTUMYJT3bjRYUgTtwTOM1N-A6bTOzp1UTF0Ap8NWJ6JvUzOAOz9r4hj7zaR2t-j357FMn5UlIMQdj6kMXeUzSrxBrV6LMCJ8gtjZLRYsKl8scD3kK-wdfMQ_aJ5cEzjU0fY5lyuXg5Kpzi-HLi6M1Upc1ON-k6kntbOTr2x0DPc1dqTomt69xAOTqLbeA_3fFy20kVimCERqtPtaK2acXUUs05jK0gRjqLterlqz5a7uD4IkbZG6PNR89aTino63B25a5A29ph4Vh0bghUmsA")
	// gUserInfo, err := l.svcCtx.GoogleOAuthService.TransferIdTokenToUserInfo(token)
	s.enforceUserInfo(&gUserInfo)
	return gUserInfo, err
}

type GoogleOAuthService struct {
}

func NewGoogleOAuthService() *GoogleOAuthService {
	return &GoogleOAuthService{}
}
