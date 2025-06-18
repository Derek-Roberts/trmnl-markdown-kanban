package oauth

import (
    "context"
    "errors"
    "os"
    "time"

    "golang.org/x/oauth2"
)

// Config holds the OAuth2 settings for TRMNL.
var Config = &oauth2.Config{
    ClientID:     os.Getenv("TRMNL_CLIENT_ID"),
    ClientSecret: os.Getenv("TRMNL_CLIENT_SECRET"),
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://usetrmnl.com/oauth/authorize",
        TokenURL: "https://usetrmnl.com/oauth/token",
    },
    RedirectURL: os.Getenv("REDIRECT_URL"),
    Scopes:      []string{},
}

// ExchangeCode swaps the one-time code for an access token.
func ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
    if code == "" {
        return nil, errors.New("empty OAuth2 code")
    }
    tok, err := Config.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }
    // refresh a minute before expiry
    tok.Expiry = tok.Expiry.Add(-1 * time.Minute)
    return tok, nil
}

// TokenProvider refreshes tokens as needed.
type TokenProvider struct {
    ts oauth2.TokenSource
}

// NewTokenProvider wraps an initial token.
func NewTokenProvider(tok *oauth2.Token) *TokenProvider {
    return &TokenProvider{
        ts: Config.TokenSource(context.Background(), tok),
    }
}

// Token returns a fresh token, refreshing if expired.
func (p *TokenProvider) Token() (*oauth2.Token, error) {
    tok, err := p.ts.Token()
    if err != nil {
        return nil, err
    }
    if !tok.Valid() {
        return nil, errors.New("invalid OAuth2 token")
    }
    return tok, nil
}
