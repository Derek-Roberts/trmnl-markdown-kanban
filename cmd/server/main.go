package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/oauth"
    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/web"
)

func main() {
    // Ensure DATA_DIR is set (default to /data)
    if os.Getenv("DATA_DIR") == "" {
        os.Setenv("DATA_DIR", "/data")
    }

    // Attempt to load existing token
    if tok, err := oauth.LoadToken(); err == nil {
        web.SetTokenProvider(oauth.NewTokenProvider(tok))
        log.Println("Loaded persisted OAuth token")
    } else {
        log.Println("No persisted token found, fresh install required")
    }

    addr := ":" + os.Getenv("PORT")
    if addr == ":" {
        addr = ":8080"
    }

    http.HandleFunc("/install", web.InstallHandler)
    http.HandleFunc("/markup", web.MarkupHandler)
    // (future) http.HandleFunc("/manage", web.ManageHandler)

    log.Printf("Server starting on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
