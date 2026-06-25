# ai-stt — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Install

```bash
togo install togo-framework/ai-stt
```

A capability plugin — it self-registers on boot; no driver selector needed.

## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `DEEPGRAM_API_KEY` | _see provider docs_ |
| `G` | _see provider docs_ |
| `OPENAI_API_KEY` | _see provider docs_ |
| `OPENAI_BASE_URL` | _see provider docs_ |
| `STT_DRIVER` | _see provider docs_ |

## Usage

```go
text, err := stt.FromKernel(k).Transcribe(ctx, audioBytes, stt.Options{})
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-stt
- README: ../README.md
