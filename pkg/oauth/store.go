package oauth

import (
    "encoding/json"
    "os"
    "path/filepath"

    "golang.org/x/oauth2"
)

// TokenFile is the path where the token will be stored.
// It defaults to DATA_DIR/token.json, where DATA_DIR is set via env var.
var TokenFile = filepath.Join(os.Getenv("DATA_DIR"), "token.json")

// SaveToken serializes the OAuth2 token to disk.
func SaveToken(tok *oauth2.Token) error {
    // Ensure directory exists with restricted permissions
    if err := os.MkdirAll(filepath.Dir(TokenFile), 0o700); err != nil {
        return err
    }
    f, err := os.OpenFile(TokenFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
    if err != nil {
        return err
    }
    defer f.Close()
    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    return enc.Encode(tok)
}

// LoadToken deserializes the OAuth2 token from disk.
func LoadToken() (*oauth2.Token, error) {
    f, err := os.Open(TokenFile)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var tok oauth2.Token
    if err := json.NewDecoder(f).Decode(&tok); err != nil {
        return nil, err
    }
    return &tok, nil
}
