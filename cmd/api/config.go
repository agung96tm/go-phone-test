package main

type Config struct {
	Addr         string
	SecretKey    string
	DB           DBConfig
	cors         CorsConfig
	googleOauth2 GoogleOauthConfig
}

type DBConfig struct {
	dsn string
}

type CorsConfig struct {
	trustedOrigins []string
}

type GoogleOauthConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	SendTokenUrl string
}

func DefaultConfig() Config {
	return Config{
		cors: CorsConfig{
			trustedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:4000",
				"http://localhost:5000",
			},
		},
		googleOauth2: GoogleOauthConfig{
			SendTokenUrl: "http://localhost:8000/v1/social/google/",
			RedirectURL:  "http://localhost:3000/auth/google/callback",
			ClientID:     "",
			ClientSecret: "",
		},
	}
}
