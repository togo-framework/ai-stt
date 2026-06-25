# ai-stt — documentation

Speech-to-Text for togo — Whisper + Deepgram drivers

## Install

```bash
togo install togo-framework/ai-stt
```

A capability plugin — it self-registers on boot; no driver selector needed.

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `DEEPGRAM_API_KEY"` |
| `OPENAI_API_KEY"` |
| `OPENAI_BASE_URL"` |
| `STT_DRIVER"` |

## Usage

```go
text, err := stt.FromKernel(k).Transcribe(ctx, audioBytes, stt.Options{})
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-stt
- Full README: ../README.md
