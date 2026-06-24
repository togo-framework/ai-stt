package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/togo-framework/togo"
)

func init() {
	RegisterDriver("deepgram", func(k *togo.Kernel) (Transcriber, error) {
		key := os.Getenv("DEEPGRAM_API_KEY")
		if key == "" {
			return nil, errors.New("ai-stt deepgram: DEEPGRAM_API_KEY not set")
		}
		return &deepgramDriver{key: key, client: &http.Client{Timeout: 120 * time.Second}}, nil
	})
}

type deepgramDriver struct {
	key    string
	client *http.Client
}

// Transcribe posts the raw audio to Deepgram's listen endpoint.
func (d *deepgramDriver) Transcribe(ctx context.Context, req Request) (Response, error) {
	q := url.Values{}
	model := req.Model
	if model == "" {
		model = "nova-2"
	}
	q.Set("model", model)
	q.Set("smart_format", "true")
	if req.Language != "" {
		q.Set("language", req.Language)
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.deepgram.com/v1/listen?"+q.Encode(), bytes.NewReader(req.Audio))
	if err != nil {
		return Response{}, err
	}
	r.Header.Set("Authorization", "Token "+d.key)
	r.Header.Set("Content-Type", "audio/*")
	resp, err := d.client.Do(r)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("ai-stt deepgram: %s: %s", resp.Status, string(raw))
	}
	var out struct {
		Results struct {
			Channels []struct {
				Alternatives []struct {
					Transcript string `json:"transcript"`
				} `json:"alternatives"`
			} `json:"channels"`
		} `json:"results"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return Response{}, err
	}
	text := ""
	if len(out.Results.Channels) > 0 && len(out.Results.Channels[0].Alternatives) > 0 {
		text = out.Results.Channels[0].Alternatives[0].Transcript
	}
	return Response{Text: text}, nil
}
