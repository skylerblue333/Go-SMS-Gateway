# Security

Sky SMS Envelope is an engineering-beta validation service and does not contact an SMS provider.

Implemented controls include bounded HTTP bodies, strict JSON decoding, basic E.164-compatible destination validation, UTF-8/body-size limits, HTTP server timeouts, graceful shutdown, race-tested Go code, vulnerability scanning, and a non-root distroless image.

This repository does not implement authentication, authorization, tenant isolation, rate limiting, provider credentials, webhook verification, consent/opt-out enforcement, regional telecom compliance, durable audit trails, abuse prevention, or DDoS protection. Do not expose it as a public messaging gateway without adding and verifying those controls.

A structurally valid phone number is not evidence of ownership, consent, reachability, or legal permission to contact that number.
