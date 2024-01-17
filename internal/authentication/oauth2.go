package authentication

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"time"
)

type contextKey string

const AccessTokenKey = "accessToken"
const IsAuthenticatedKey = contextKey("isAuthenticated")

type UserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type GoogleOauth2 struct {
	Config            *oauth2.Config
	OauthGoogleUrlAPI string
	Oauthstate        string
	SendTokenUrl      string
}

func NewGoogleOauth2(redirectUrl, sendTokenUrl, clientID, clientSecret string) *GoogleOauth2 {
	return &GoogleOauth2{
		OauthGoogleUrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
		Config: &oauth2.Config{
			RedirectURL:  redirectUrl,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		SendTokenUrl: sendTokenUrl,
		Oauthstate:   "oauthstate",
	}
}

func (g GoogleOauth2) GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: g.Oauthstate, Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (g GoogleOauth2) StateValid(r *http.Request) bool {
	oauthState, _ := r.Cookie(g.Oauthstate)
	return r.FormValue("state") == oauthState.Value
}

func (g GoogleOauth2) GetGoogleToken(r *http.Request) (string, error) {
	code := r.FormValue("code")
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	return token.AccessToken, nil
}

func (g GoogleOauth2) GetUserDataByToken(token string) (*UserData, error) {
	response, err := http.Get(g.OauthGoogleUrlAPI + token)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	var userData UserData
	if err := json.Unmarshal(contents, &userData); err != nil {
		return nil, err
	}
	return &userData, nil
}

type Token struct {
	Token string `json:"token"`
}

type TokenResp struct {
	Access string `json:"access"`
}

func (g GoogleOauth2) SendLoginGoogle(token string) (*TokenResp, error) {
	reqToken := Token{Token: token}
	body, _ := json.Marshal(reqToken)
	resp, err := http.Post(g.SendTokenUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var result TokenResp
	if err := json.Unmarshal(responseBody, &result); err != nil {
		fmt.Println(string(responseBody))
		return nil, err
	}

	return &result, nil
}
