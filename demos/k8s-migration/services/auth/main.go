package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

type TokenResponse struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
}

type AuthService struct {
    jwtSecret []byte
    port      string
}

func NewAuthService() *AuthService {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        panic("JWT_SECRET not set")
    }
    return &AuthService{
        jwtSecret: []byte(secret),
        port:      os.Getenv("PORT"),
    }
}

func (s *AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
    // Simplified — real impl calls user-service gRPC
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(TokenResponse{
        Token:     "jwt.token.here",
        ExpiresAt: time.Now().Add(24 * time.Hour),
    })
}

func main() {
    svc := NewAuthService()
    http.HandleFunc("/login", svc.HandleLogin)
    fmt.Printf("auth-service listening on :%s\n", svc.port)
    http.ListenAndServe(":"+svc.port, nil)
}
