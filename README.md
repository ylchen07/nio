# Nio v0.1 — Personal AI representative (Go)

Nio is a **working minimal gateway**, not yet a complete Agent2Agent (A2A) implementation. No A2A conformance is claimed; the discovery endpoint is explicitly provisional. There is no `systemd` configuration.

## Build

Requires **Go 1.27.4 or newer**. No third-party dependencies.

> Note: The source targets Go 1.27.4. This release was previously smoke-tested with Go 1.23.2, but a Go 1.27.4 compiler is unavailable in the current build environment, so the updated source has not yet been rebuilt or tested with 1.27.4. Do not treat the previously distributed Linux executable as rebuilt with Go 1.27.4.

```sh
go test ./...
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o nio .
```

The binary is self-contained; it writes a local `nio-data.json` file to preserve messages, agent identities and approval decisions. **The v0.1 persistence format is JSON, not SQLite.** Back up the file and do not run concurrent instances writing the same path.

## Run

```sh
export NIO_ADMIN_TOKEN="$(openssl rand -hex 32)" # Save this securely!
export OPENAI_API_KEY="your-openai-api-key"        # Optional for offline testing
export NIO_LISTEN=127.0.0.1:8080
export NIO_PUBLIC_URL=https://agent.example.com
./nio
```

Defaults: working hours 09:00–17:00, Monday-Friday in `America/Toronto`. Current MVP working-hour values are hard-coded but the timezone is configurable via `NIO_TIMEZONE`. With no `OPENAI_API_KEY`, Nio returns a safe fixed acknowledgment so you can test locally. With an API key, it uses the configured `NIO_MODEL` (`gpt-4.1-mini` by default) via OpenAI Chat Completions. Adjust to an available model in your account.

## Try the API

```sh
curl -s http://127.0.0.1:8080/healthz
curl -s http://127.0.0.1:8080/SKILL.md
curl -s -X POST http://127.0.0.1:8080/v1/introductions \
  -H 'Content-Type: application/json' \
  -d '{"name":"AliceAgent","purpose":"Request a technical meeting"}'
# Save the returned introduction ID
curl -s http://127.0.0.1:8080/admin/queue -H "Authorization: Bearer $NIO_ADMIN_TOKEN"
curl -s -X POST http://127.0.0.1:8080/admin/introductions/approve \
  -H "Authorization: Bearer $NIO_ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"id":"INTRO_ID"}'
# Save the returned agent token and deliver it securely after verifying the agent independently.
curl -s -X POST http://127.0.0.1:8080/v1/messages \
  -H "Authorization: Bearer AGENT_TOKEN" -H 'Content-Type: application/json' \
  -d '{"text":"Can we discuss a meeting tomorrow?","request_approval":true}'
```

Use `GET /admin/queue` to see introductions and approvals, and `POST /admin/approvals/decide` with `{"id":"APPROVAL_ID","decision":"approved"}` or `"rejected"` to decide. These decisions are stored; outbound notifications and proactive follow-up messages are **not implemented** yet.

## Security and limitations

- Bind to localhost unless placed behind an HTTPS reverse proxy with rate limiting. Never expose the admin API over plain HTTP on the public internet.
- Unknown agents can only submit rate-limited introductions; a human must independently verify claimed identity before issuing a token. The server does not perform proof-of-control challenges yet.
- Authenticated tokens are SHA-256 hashed at rest, but agent tokens are handed out once in the administrative approval response. Protect that response.
- Single-process in-memory rate limiting; remote IP comes from the network connection, not untrusted proxy headers.
- A single on-disk JSON snapshot is updated after mutations. It has no transaction engine or multi-process locking.
- Requesting an approval creates a pending record, **not a real-world commitment**; the AI does not autonomously book meetings, initiate payments, or send email.
- Authentication and model safety prompts are foundations, not a substitute for strict structured authorization of future real-world tools.
- Full A2A 1.0 interoperability, SQLite, robust durable queues, identity proof-of-control, real human notifications and configurable working hour policies are next milestones.
