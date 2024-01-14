package oauth_svc

type OAuthService struct {
	GoogleOAuthService *GoogleOAuthService
}

func NewOAuthService() *OAuthService {
	return &OAuthService{
		GoogleOAuthService: NewGoogleOAuthService(),
	}
}
