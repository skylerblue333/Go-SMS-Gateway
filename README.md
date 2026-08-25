# Sky SMS Envelope

**Status: engineering beta.** A bounded Go HTTP service for validating SMS message envelopes before a separately configured provider integration.

## Implemented

- Basic E.164-compatible destination validation (`+` plus 8–15 digits, non-zero country prefix).
- UTF-8 message body validation with a 1,600-rune / 8 KiB ceiling.
- Optional sender identifier capped at 32 bytes.
- Strict single-object JSON parsing with unknown-field rejection.
- `/healthz` and `/readyz` operational endpoints.
- `POST /v1/messages/validate` response includes body rune/byte counts and always reports `provider_dispatch: false`.
- HTTP server timeouts and graceful shutdown.
- Go 1.26 format/vet/test/race/govulncheck/build gates plus non-root container smoke testing.

Example:

```bash
curl -sS http://127.0.0.1:8080/v1/messages/validate \
  -H 'content-type: application/json' \
  --data '{"to":"+14155550123","body":"Hello from Sky"}'
```

## Scope limitations

This repository does **not send SMS messages**. It has no Twilio/carrier/provider integration, delivery receipts, retries, durable queue, sender registration, opt-in/opt-out compliance workflows, regional telecom policy, tenant isolation, billing, rate limiting, HA, or production deployment.

Basic E.164 formatting is only structural validation; it does not prove that a number exists, is reachable, or is legally permitted to receive a message.

## SKYCOIN4444 integration

Use this service as a pre-provider validation boundary. Provider credentials, consent/compliance, durable delivery state, retries, observability, and billing must remain in separately verified integrations.
