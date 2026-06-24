<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/ai-stt</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/ai-stt"><img src="https://pkg.go.dev/badge/github.com/togo-framework/ai-stt.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/ai-stt
```

<!-- /togo-header -->

# ai-stt — Speech-to-Text for togo

Adds Speech-to-Text to a togo app on the `ai` plugin family. Drivers register via
`init()`; pick one with `STT_DRIVER`.

```bash
togo install togo-framework/ai-stt
```

## Drivers

| `STT_DRIVER` | Provider | Env |
|---|---|---|
| `openai` | OpenAI Whisper | `OPENAI_API_KEY` (`OPENAI_BASE_URL` optional) |
| `deepgram` | Deepgram | `DEEPGRAM_API_KEY` |

## Use

Mount the REST handler under `/api/ai/stt` and `POST` the raw audio (body) →
`{"text": "..."}`. Or from Go:

```go
svc, _ := stt.FromKernel(k)
res, _ := svc.Transcribe(ctx, stt.Request{Audio: bytes, Model: "whisper-1"})
```

MIT

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
