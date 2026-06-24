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
