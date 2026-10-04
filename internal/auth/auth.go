package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "net/http"
    "strings"
    "sync"
    "time"
)

type TokenRecord struct {
    Email string
    Role  string
    Exp   time.Time
}

var (
    tokenMu sync.RWMutex
    tokens  = map[string]TokenRecord{}
)

func CreateToken(email, role string) (string, error) {
    payload := fmt.Sprintf("%s|%s|%d", email, role, time.Now().Add(1*time.Hour).Unix())
    sig := hmac.New(sha256.New, []byte("private-owner-key"))
    _, err := sig.Write([]byte(payload))
    if err != nil {
        return "", err
    }
    token := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))

    tokenMu.Lock()
    tokens[token] = TokenRecord{Email: email, Role: role, Exp: time.Now().Add(1 * time.Hour)}
    tokenMu.Unlock()
    return token, nil
}

func ValidateToken(token string) (TokenRecord, bool) {
    tokenMu.RLock()
    rec, ok := tokens[token]
    tokenMu.RUnlock()
    if !ok {
        return TokenRecord{}, false
    }
    if time.Now().After(rec.Exp) {
        tokenMu.Lock()
        delete(tokens, token)
        tokenMu.Unlock()
        return TokenRecord{}, false
    }
    return rec, true
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == "" {
            http.Error(w, "missing authorization token", http.StatusUnauthorized)
            return
        }

        rec, ok := ValidateToken(token)
        if !ok {
            http.Error(w, "invalid or expired token", http.StatusUnauthorized)
            return
        }
        if rec.Role != role {
            http.Error(w, "role mismatch", http.StatusForbidden)
            return
        }

        next(w, r)
    }
}
