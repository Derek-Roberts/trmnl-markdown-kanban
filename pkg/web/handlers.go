package web

import (
    "bytes"
    "context"
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "os"

    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/kanban"
    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/oauth"
)

// installRequest is TRMNL’s payload to /install
type installRequest struct {
    Code                   string `json:"code"`
    InstallationCallbackURL string `json:"installation_callback_url"`
}

// InstallHandler handles the OAuth installation webhook.
func InstallHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
        return
    }

    var req installRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Exchange code for token
    tok, err := oauth.ExchangeCode(r.Context(), req.Code)
    if err != nil {
        log.Printf("token exchange error: %v", err)
        http.Error(w, "token exchange failed", http.StatusInternalServerError)
        return
    }
    // TODO: persist tok.AccessToken, tok.RefreshToken, tok.Expiry, etc.

    // Call installation callback
    cbBody := map[string]string{"installation_id": tok.AccessToken} // or your own ID
    bodyBytes, _ := json.Marshal(cbBody)
    cbReq, _ := http.NewRequestWithContext(
        context.Background(),
        "POST",
        req.InstallationCallbackURL,
        bytes.NewBuffer(bodyBytes),
    )
    cbReq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
    cbReq.Header.Set("Content-Type", "application/json")

    if resp, err := http.DefaultClient.Do(cbReq); err != nil || resp.StatusCode >= 300 {
        log.Printf("callback failed: %v, status: %v", err, resp.StatusCode)
        http.Error(w, "callback failed", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// Pre-parse the template
var boardTmpl = template.Must(template.ParseFiles("templates/board.liquid"))

// MarkupHandler renders the Kanban board as HTML/Liquid
func MarkupHandler(w http.ResponseWriter, r *http.Request) {
    board, err := kanban.LoadBoard("/data/kanban.md")
    if err != nil {
        http.Error(w, "failed to load board: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if err := boardTmpl.Execute(w, board); err != nil {
        http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
    }
}
