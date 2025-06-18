package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/oauth"
    "github.com/Derek-Roberts/trmnl-markdown-kanban/pkg/web"
)

func main() {
    addr := ":" + os.Getenv("PORT")
    if addr == ":" {
        addr = ":8080"
    }

    // Webhook routes
    http.HandleFunc("/install", web.InstallHandler)
    http.HandleFunc("/markup", web.MarkupHandler)
    // (future) http.HandleFunc("/manage", web.ManageHandler)

    log.Printf("Server starting on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
