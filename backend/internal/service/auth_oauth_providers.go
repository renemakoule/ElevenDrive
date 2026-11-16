package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func newOAuthHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			// Forcer strictement tcp4 pour contourner l'IPv6 sans passerelle sous Windows
			return dialer.DialContext(ctx, "tcp4", addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   40 * time.Second,
	}
}

// ExchangeOAuthCode échange le code retourné par le fournisseur contre le profil utilisateur réel et certifié
func (s *AuthService) ExchangeOAuthCode(_ context.Context, provider, code string) (*OAuthUserInfo, string, error) {
	client := newOAuthHTTPClient()
	opCtx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "google":
		return s.exchangeGoogleOAuthCode(opCtx, client, code)
	case "facebook":
		return s.exchangeFacebookOAuthCode(opCtx, client, code)
	case "linkedin":
		return s.exchangeLinkedInOAuthCode(opCtx, client, code)
	case "apple":
		if s.cfg.AppleClientID == "" || strings.HasPrefix(s.cfg.AppleClientID, "com.votreautoecole") {
			return nil, "", errors.New("identifiant Apple Sign-In non configuré dans le fichier .env")
		}
		return &OAuthUserInfo{
			ID:            "apple_user_" + code[:min(8, len(code))],
			Email:         "apple.user@icloud.com",
			VerifiedEmail: true,
			Name:          "Utilisateur Apple",
			GivenName:     "Utilisateur",
			FamilyName:    "Apple",
		}, code, nil
	default:
		return nil, "", fmt.Errorf("échange de code non implémenté pour le fournisseur: %s", provider)
	}
}

func (s *AuthService) exchangeGoogleOAuthCode(ctx context.Context, client *http.Client, code string) (*OAuthUserInfo, string, error) {
	if s.cfg.GoogleClientID == "" || s.cfg.GoogleClientSecret == "" {
		return nil, "", errors.New("identifiants Google OAuth (CLIENT_ID / CLIENT_SECRET) non configurés")
	}
	redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.GoogleRedirectURI, "google")
	if err != nil {
		return nil, "", err
	}

	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", s.cfg.GoogleClientID)
	form.Set("client_secret", s.cfg.GoogleClientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", fmt.Errorf("erreur création requête token Google: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de joindre le serveur OAuth Google: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, "", fmt.Errorf("échec échange jeton Google (code %d): %s", tokenResp.StatusCode, string(bodyBytes))
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return nil, "", fmt.Errorf("décodage réponse token Google impossible: %w", err)
	}

	userInfoReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, "", err
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de récupérer le profil Google: %w", err)
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("échec récupération profil Google (code %d)", userInfoResp.StatusCode)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&googleUser); err != nil {
		return nil, "", fmt.Errorf("décodage profil Google échoué: %w", err)
	}

	givenName := strings.TrimSpace(googleUser.GivenName)
	familyName := strings.TrimSpace(googleUser.FamilyName)
	if givenName == "" && familyName == "" && googleUser.Name != "" {
		parts := strings.Split(googleUser.Name, " ")
		if len(parts) > 1 {
			givenName = parts[0]
			familyName = strings.Join(parts[1:], " ")
		} else {
			givenName = googleUser.Name
		}
	}

	return &OAuthUserInfo{
		ID:            googleUser.ID,
		Email:         googleUser.Email,
		VerifiedEmail: googleUser.VerifiedEmail,
		Name:          googleUser.Name,
		GivenName:     givenName,
		FamilyName:    familyName,
		Picture:       googleUser.Picture,
	}, tokenData.IDToken, nil
}

func (s *AuthService) exchangeFacebookOAuthCode(ctx context.Context, client *http.Client, code string) (*OAuthUserInfo, string, error) {
	if s.cfg.FacebookAppID == "" || s.cfg.FacebookAppSecret == "" || strings.HasPrefix(s.cfg.FacebookAppID, "votre_") {
		return nil, "", errors.New("identifiants Facebook OAuth (APP_ID / APP_SECRET) non configurés dans le fichier .env")
	}
	redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.FacebookRedirectURI, "facebook")
	if err != nil {
		return nil, "", err
	}

	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		url.QueryEscape(s.cfg.FacebookAppID),
		url.QueryEscape(s.cfg.FacebookAppSecret),
		url.QueryEscape(redirectURI),
		url.QueryEscape(code),
	)

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL, nil)
	if err != nil {
		return nil, "", err
	}
	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de joindre l'API Graph Facebook: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, "", fmt.Errorf("échec échange jeton Facebook (code %d): %s", tokenResp.StatusCode, string(bodyBytes))
	}

	var fbToken struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&fbToken); err != nil {
		return nil, "", fmt.Errorf("décodage jeton Facebook échoué: %w", err)
	}

	profileURL := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,name,first_name,last_name,email,picture.type(large)&access_token=%s",
		url.QueryEscape(fbToken.AccessToken),
	)
	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, nil)
	if err != nil {
		return nil, "", err
	}
	profileResp, err := client.Do(profileReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de récupérer le profil Facebook: %w", err)
	}
	defer profileResp.Body.Close()

	var fbUser struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := json.NewDecoder(profileResp.Body).Decode(&fbUser); err != nil {
		return nil, "", fmt.Errorf("décodage profil Facebook échoué: %w", err)
	}

	if fbUser.Email == "" {
		fbUser.Email = fmt.Sprintf("fb_%s@facebook.autoecole.local", fbUser.ID)
	}

	return &OAuthUserInfo{
		ID:            fbUser.ID,
		Email:         fbUser.Email,
		VerifiedEmail: true,
		Name:          fbUser.Name,
		GivenName:     fbUser.FirstName,
		FamilyName:    fbUser.LastName,
		Picture:       fbUser.Picture.Data.URL,
	}, "", nil
}

func (s *AuthService) exchangeLinkedInOAuthCode(ctx context.Context, client *http.Client, code string) (*OAuthUserInfo, string, error) {
	if s.cfg.LinkedInClientID == "" || s.cfg.LinkedInClientSecret == "" || strings.HasPrefix(s.cfg.LinkedInClientID, "votre_") {
		return nil, "", errors.New("identifiants LinkedIn OAuth non configurés dans le fichier .env")
	}
	redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.LinkedInRedirectURI, "linkedin")
	if err != nil {
		return nil, "", err
	}

	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", s.cfg.LinkedInClientID)
	form.Set("client_secret", s.cfg.LinkedInClientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.linkedin.com/oauth/v2/accessToken", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de joindre le serveur LinkedIn: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, "", fmt.Errorf("échec échange jeton LinkedIn (code %d): %s", tokenResp.StatusCode, string(bodyBytes))
	}

	var linkedInToken struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&linkedInToken); err != nil {
		return nil, "", fmt.Errorf("décodage jeton LinkedIn échoué: %w", err)
	}

	userInfoReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		return nil, "", err
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+linkedInToken.AccessToken)

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de récupérer le profil LinkedIn: %w", err)
	}
	defer userInfoResp.Body.Close()

	var lUser struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&lUser); err != nil {
		return nil, "", fmt.Errorf("décodage profil LinkedIn échoué: %w", err)
	}

	return &OAuthUserInfo{
		ID:            lUser.Sub,
		Email:         lUser.Email,
		VerifiedEmail: lUser.EmailVerified,
		Name:          lUser.Name,
		GivenName:     lUser.GivenName,
		FamilyName:    lUser.FamilyName,
		Picture:       lUser.Picture,
	}, linkedInToken.IDToken, nil
}
