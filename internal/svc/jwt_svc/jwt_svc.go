package jwt_svc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kilicmu/user-service/internal/database"
)

const (
	ACCESS_TOKEN              = "access_token"
	REFRESH_TOKEN             = "refresh_token"
	COUNT_OF_HOUR             = 72
	ACCESS_TOKEN_Expriration  = time.Hour * COUNT_OF_HOUR
	REFRESH_TOKEN_Expriration = time.Hour * COUNT_OF_HOUR * 3

	ACCESS_TOKEN_Expriration_SECOUND_COUNT  = 72 * 60 * 60
	REFRESH_TOKEN_Expriration_SECOUND_COUNT = 3 * ACCESS_TOKEN_Expriration_SECOUND_COUNT

	UID_TO_ACCESS_AND_REFRESH_PAIR = "UID_TO_ACCESS_AND_REFRESH_PAIR"
)

type UserInfoJWTPayload struct {
	jwt.StandardClaims
	TokenType string `json:"token_type"`
	UId       string `json:"uid"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTService struct{}

func generateJwtAccessTokenByUserInfo(info *database.UserInfo) (string, error) {
	claims := UserInfoJWTPayload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_Expriration).Unix(),
		},
		UId:       info.UId,
		TokenType: ACCESS_TOKEN,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(info.Password))

	if err != nil {
		log.Println("generate jwt error:", err)
		return "", err
	}

	return tokenString, nil
}

func generateJwtRefreshTokenByUserId(info *database.UserInfo) (string, error) {
	claims := UserInfoJWTPayload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_Expriration).Unix(),
		},
		UId:       info.UId,
		TokenType: REFRESH_TOKEN,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(info.Password))

	if err != nil {
		log.Println("generate jwt error:", err)
		return "", err
	}

	return tokenString, nil
}

func (s *JWTService) GenerateJwtTokenPairByUserInfo(info *database.UserInfo) (*TokenPair, error) {
	accessToken, err := generateJwtAccessTokenByUserInfo(info)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJwtRefreshTokenByUserId(info)
	if err != nil {
		return nil, err
	}

	tokenPair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	cacheIdOfUidToTokenPair := fmt.Sprintf("%s/%s", UID_TO_ACCESS_AND_REFRESH_PAIR, info.UId)
	if alreadyExist, err := database.GetRedisUtils().Exists(cacheIdOfUidToTokenPair); err == nil {
		if alreadyExist {
			database.GetRedisUtils().Delete(cacheIdOfUidToTokenPair)
		}
		if json, err := json.Marshal(tokenPair); err == nil {
			database.GetRedisUtils().SetEx(cacheIdOfUidToTokenPair, string(json), ACCESS_TOKEN_Expriration_SECOUND_COUNT)
		}
	}
	return tokenPair, nil
}

func (s *JWTService) GetTokenPairByUserInfo(info *database.UserInfo) (*TokenPair, error) {
	tokenPair := &TokenPair{}
	cacheIdOfUidToTokenPair := fmt.Sprintf("%s/%s", UID_TO_ACCESS_AND_REFRESH_PAIR, info.UId)
	if alreadyExist, err := database.GetRedisUtils().Exists(cacheIdOfUidToTokenPair); err == nil && alreadyExist {
		if tokenPairjson, err := database.GetRedisUtils().Get(cacheIdOfUidToTokenPair); err == nil {
			if err := json.Unmarshal([]byte(tokenPairjson), tokenPair); err == nil {
				return tokenPair, nil
			}
		}

	}
	return tokenPair, fmt.Errorf("token pair not exist")
}

// only decode payload, not varify token
func (s *JWTService) UnsafeGetUserInfoPayload(token string) (*UserInfoJWTPayload, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	encodedTokens := strings.Split(token, ".")
	if len(encodedTokens) != 3 {
		return nil, fmt.Errorf("split token '%s' fail", token)
	}
	encodedPayload := encodedTokens[1]
	if decodedPayload, err := base64.RawStdEncoding.DecodeString(encodedPayload); err == nil {
		payload := UserInfoJWTPayload{}
		if err := json.Unmarshal(decodedPayload, &payload); err == nil {
			return &payload, nil
		}
	}
	return nil, fmt.Errorf("decode payload error")
}

func (s *JWTService) varifyTokenAlive(token string, info *database.UserInfo) (bool, error) {
	if token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, nil
		}
		return []byte(info.Password), nil
	}); err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
				return false, fmt.Errorf("token expired")
			}
			return true, nil
		}
	}
	return false, nil
}

func (s *JWTService) VarifyAccessTokenAlive(token string, info *database.UserInfo) (bool, error) {
	if tokenPair, err := s.GetTokenPairByUserInfo(info); err == nil {
		if tokenPair.AccessToken == token {
			return s.varifyTokenAlive(token, info)
		}
	}
	return false, nil
}

func (s *JWTService) VarifyRefreshTokenAlive(token string, info *database.UserInfo) (bool, error) {
	if tokenPair, err := s.GetTokenPairByUserInfo(info); err == nil {
		if tokenPair.RefreshToken == token {
			return s.varifyTokenAlive(token, info)
		}
	}
	return false, nil
}

func NewJwtService() *JWTService {
	return &JWTService{}
}
