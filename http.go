package stt

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/togo-framework/togo"
)

// Handler exposes the STT service over REST. Mount under /api/ai/stt:
//
//	mux.Handle("/api/ai/stt/", http.StripPrefix("/api/ai/stt", stt.Handler(k)))
//
// POST /  with the raw audio as the body (Content-Type audio/*) and optional
// ?model=&language=&filename= -> JSON {"text": "..."}.
func Handler(k *togo.Kernel) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		svc, ok := FromKernel(k)
		if !ok {
			http.Error(w, "ai-stt not configured (set STT_DRIVER)", http.StatusInternalServerError)
			return
		}
		audio, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<20))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q := r.URL.Query()
		res, err := svc.Transcribe(r.Context(), Request{
			Audio:    audio,
			Filename: q.Get("filename"),
			Model:    q.Get("model"),
			Language: q.Get("language"),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})
	return mux
}
