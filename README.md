# Fiat-Verified Wallet Bridge

Status: Private/Internal Use Only
Audit Access: Public-Readable (for review and compliance examination)

A private-use wallet and fiat bridge foundation for:
- formally verified cryptographic primitives
- auditable wallet logic
- bridge/conversion pipeline design
- compliance and audit logging
- Ledger-ready integration model

This is not a public payment product and should not be treated as a licensed financial service.

## Purpose

This repository is designed to provide a secure foundation for:
- verified wallet cryptography
- fiat bridge orchestration
- auditable transaction logs
- signature and approval workflows
- governance and compliance artifacts

It is intended for internal/private deployment and review, not mass-market public use.

## Important Notice

- This project is not a fiat currency and not a legal tender system.
- It is not a licensed money transmitter or banking platform.
- It does not connect to a live wallet, Ledger, or real funds by default.
- It should be used only in controlled private/internal deployments.
- Public release requires regulatory review, licensing review, banking/settlement partner review, and independent security audits.

## Architecture Summary

```
Wallet Core
  -> Key management
  -> Signing logic
  -> Multi-signature approval
  -> Audit log

Bridge Layer
  -> Rate oracle
  -> Conversion engine
  -> Settlement orchestration
  -> Fee & policy checks

Compliance Layer
  -> KYC/AML hooks (optional)
  -> Policy engine
  -> Transaction approvals
  -> Logging and immutable records

Verification Layer
  -> formal specification for critical operations
  -> reproducible build artifacts
  -> audit-friendly design
```

## Repo Layout

```
fiat-verified-wallet-bridge/
├── README.md
├── LICENSE
├── PRIVATE_USE_NOTICE.md
├── ARCHITECTURE.md
├── COMPLIANCE.md
├── Makefile
├── go.mod
├── go.sum
├── crypto/
│   ├── wallet/
│   │   ├── key_manager.go
│   │   ├── signing.go
│   │   ├── multisig.go
│   │   └── spec.md
│   ├── verified/
│   │   ├── secp256k1.go
│   │   ├── field_arithmetic.go
│   │   └── spec.md
│   └── tests/
│       └── wallet_test.go
├── bridge/
│   ├── oracle/
│   │   ├── rate_oracle.go
│   │   └── oracle_spec.md
│   ├── conversion/
│   │   ├── converter.go
│   │   ├── fees.go
│   │   └── conversion_spec.md
│   └── tests/
│       └── bridge_test.go
├── compliance/
│   ├── kyc/
│   │   ├── kyc.go
│   │   └── kyc_spec.md
│   ├── aml/
│   │   ├── aml.go
│   │   └── aml_spec.md
│   ├── audit_log/
│   │   ├── ledger.go
│   │   └── ledger_spec.md
│   └── tests/
│       └── compliance_test.go
├── api/
│   ├── server.go
│   ├── handlers/
│   ├── middleware/
│   └── openapi.yaml
├── cli/
│   ├── main.go
│   └── commands/
├── build/
│   ├── build.sh
│   ├── verify.sh
│   └── Dockerfile
├── docs/
│   ├── architecture/
│   ├── compliance/
│   └── formal_specs/
├── .github/
│   └── workflows/
├── .gitignore
└── tests/
    └── integration_test.go
```

## Build and Test

```bash
make build
make test
```

## Security Model

- internal/private deployments only
- no hardcoded live wallet credentials
- no public key material or production secrets committed
- audit-log signatures enabled for review traceability
- reproducible build verification recommended for all release artifacts

## Ledger Readiness

This repository is structured for future Ledger integration, but the current codebase intentionally does not bind to a live main wallet or any production account.

Ledger compatibility should be added only when:
- the architecture is reviewed
- approvals are in place
- testnet flows are validated
- compliance model is reviewed

## Public Audit Access

The repo is readable publicly, but the design intent is private deployment. All relevant artifacts should remain reviewable for auditing and formal verification assessment.

## License

Apache 2.0

## Disclaimer

This repository should be treated as an internal, auditable foundation and not as a live public-facing financial product.
