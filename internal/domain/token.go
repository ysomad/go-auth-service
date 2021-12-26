package domain

// Token represents JWT access token contains account id in payload.
type Token struct {
	AccessToken string `json:"accessToken"`
}
