# Fiat-Verified Wallet Bridge

Private owner-and-partner project.

This project is intentionally private and designed for the owner’s own security workflows:
- private admin-only access
- trusted-device recognition
- browser-side monitoring for authorized environments
- event delivery for account activity
- suspicious-domain and suspicious-account detection
- private audit trail and emergency block actions

This repository is not meant for public release. Keep all secrets and credentials outside the repo in `.env.local` and encrypted offline storage.

## Local run

```bash
go run ./cmd/server
```

## Protected endpoints

- `GET /health`
- `POST /api/auth/login`
- `GET /api/admin/dashboard`
- `POST /api/admin/device/trust`
- `POST /api/admin/transfer`
- `GET /api/admin/audit`

## Root rules

- Only the admin and the trusted partner can access private operations.
- Device trust is required before authenticated access proceeds.
- Any suspicious or unauthorized transfer must be blocked and recorded.
- Browser monitoring is limited to the owner’s own authorized browser profiles and environments.
