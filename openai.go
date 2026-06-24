package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/togo-framework/togo"
)

func init() {
	RegisterDriver("openai", func(k *togo.Kernel) (Transcriber, error) {
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			return nil, errors.New("ai-stt openai: OPENAI_API_KEY not set")
		}
		base := os.Getenv("OPENAI_BASE_URL")
		if base == "" {
			base = "https://api.openai.com/v1"
		}
		return &openaiDriver{key: key, base: base, client: &http.Client{Timeout: 120 * time.Second}}, nil
	})
}

type openaiDriver struct {
	key    string
	base   string
	client *http.Client
}

// Transcribe posts the audio to OpenAI's Whisper transcription endpoint.
func (d *openaiDriver) Transcribe(ctx context.Context, req Request) (Response, error) {
	model := req.Model
	if model == "" {
		model = "whisper-1"
	}
	fname := req.Filename
	if fname == "" {
		fname = "audio.mp3"
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("file", fname)
	if err != nil {
		return Response{}, err
	}
	if _, err := fw.Write(req.Audio); err != nil {
		return Response{}, err
	}
	_ = w.WriteField("model", model)
	if req.Language != "" {
		_ = w.WriteField("language", req.Language)
	}
	w.Close()

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, d.base+"/audio/transcriptions", &body)
	if err != nil {
		return Response{}, err
	}
	r.Header.Set("Authorization", "Bearer "+d.key)
	r.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := d.client.Do(r)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("ai-stt openai: %s: %s", resp.Status, string(raw))
	}
	var out struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return Response{}, err
	}
	return Response{Text: out.Text}, nil
}
