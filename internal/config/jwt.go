package config

type JWTConfig struct {
	AccessTokenConfig  AccessTokenConfig
	RefreshTokenConfig RefreshTokenConfig
}

type AccessTokenConfig struct {
	secret                 []byte
	TokenExpirationMinutes int
	Audience               string
}

func (a *AccessTokenConfig) GetSecret() []byte {
	return a.secret
}

type RefreshTokenConfig struct {
	secret                 []byte
	TokenExpirationMinutes int
	Audience               string
}

func (r *RefreshTokenConfig) GetSecret() []byte {
	return r.secret
}
