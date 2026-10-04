package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/v6jdg2nbww-techRhu3e/fiat-verified-wallet-bridge/internal/auth"
    "github.com/v6jdg2nbww-techRhu3e/fiat-verified-wallet-bridge/internal/config"
    "github.com/v6jdg2nbww-techRhu3e/fiat-verified-wallet-bridge/internal/device"
    "github.com/v6jdg2nbww-techRhu3e/fiat-verified-wallet-bridge/internal/store"
)

func main() {
    cfg := config.Load()
    storeSvc := store.NewService()
    deviceSvc := device.NewService()

    adminUser := &store.User{ID: "admin-1", Email: cfg.AdminEmail, Role: "admin"}
    storeSvc.UpsertUser(adminUser)
    deviceSvc.RegisterDevice("admin-device-1", "owner-workstation", cfg.AdminEmail, true)

    mux := http.NewServeMux()

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "mode": cfg.Mode})
    })

    mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            DeviceID string `json:"device_id"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        if req.Email != cfg.AdminEmail || req.Password != cfg.AdminPassword {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }

        if req.DeviceID == "" {
            req.DeviceID = "admin-device-1"
        }

        if !deviceSvc.IsTrusted(req.DeviceID, req.Email) {
            http.Error(w, "device requires trust approval", http.StatusForbidden)
            return
        }

        token, err := auth.CreateToken(req.Email, "admin")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        storeSvc.LogAudit("admin_login", req.Email, fmt.Sprintf("device=%s", req.DeviceID))
        _ = json.NewEncoder(w).Encode(map[string]any{"token": token, "role": "admin"})
    })

    mux.HandleFunc("/api/admin/dashboard", auth.RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
        _ = json.NewEncoder(w).Encode(map[string]any{
            "mode": cfg.Mode,
            "admin_email": cfg.AdminEmail,
            "trusted_devices": deviceSvc.ListTrustedDevices(),
            "status": "private access only",
            "audit_count": len(storeSvc.ListAudit()),
            "time": time.Now().UTC().Format(time.RFC3339),
        })
    }))

    mux.HandleFunc("/api/admin/device/trust", auth.RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req struct {
            DeviceID string `json:"device_id"`
            Label    string `json:"label"`
            Email    string `json:"email"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        deviceSvc.RegisterDevice(req.DeviceID, req.Label, req.Email, true)
        storeSvc.LogAudit("device_trusted", req.Email, fmt.Sprintf("device=%s", req.DeviceID))
        _ = json.NewEncoder(w).Encode(map[string]any{"status": "trusted", "device_id": req.DeviceID})
    }))

    mux.HandleFunc("/api/admin/transfer", auth.RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var payload struct {
            To         string  `json:"to"`
            Amount     float64 `json:"amount"`
            Currency   string  `json:"currency"`
            DeviceID   string  `json:"device_id"`
            Reason     string  `json:"reason"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        if !deviceSvc.IsTrusted(payload.DeviceID, cfg.AdminEmail) {
            http.Error(w, "device not trusted for transfer", http.StatusForbidden)
            return
        }

        if payload.Amount > 50000 {
            http.Error(w, "transfer exceeds threshold; partner review required", http.StatusForbidden)
            return
        }

        storeSvc.LogAudit("transfer_request", cfg.AdminEmail, fmt.Sprintf("to=%s amount=%.2f currency=%s device=%s reason=%s", payload.To, payload.Amount, payload.Currency, payload.DeviceID, payload.Reason))
        _ = json.NewEncoder(w).Encode(map[string]any{"status": "queued", "amount": payload.Amount})
    }))

    mux.HandleFunc("/api/admin/audit", auth.RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
        _ = json.NewEncoder(w).Encode(storeSvc.ListAudit())
    }))

    mux.HandleFunc("/api/monitor/submit", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var payload map[string]any
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        storeSvc.LogAudit("monitor_event", "owner-guard", fmt.Sprintf("payload=%v", payload))
        w.WriteHeader(http.StatusAccepted)
        _ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
    })

    log.Printf("private owner guard listening on %s", cfg.ListenAddr)
    log.Fatal(http.ListenAndServe(cfg.ListenAddr, loggingMiddleware(mux)))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
