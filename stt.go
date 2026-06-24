// Package stt adds Speech-to-Text to togo. Drivers (openai/Whisper, deepgram, …)
// register via init(); select one with STT_DRIVER. Mirrors the ai/tts driver
// pattern. Mount the REST handler under /api/ai/stt.
package stt

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/togo-framework/togo"
)

// Request is a transcription request. Audio holds the raw audio bytes.
type Request struct {
	Audio    []byte `json:"-"`
	Filename string `json:"filename,omitempty"`
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

// Response holds the transcribed text.
type Response struct {
	Text string `json:"text"`
}

// Transcriber is the STT driver interface every driver implements.
type Transcriber interface {
	Transcribe(ctx context.Context, req Request) (Response, error)
}

// DriverFactory builds a Transcriber from the kernel/env.
type DriverFactory func(k *togo.Kernel) (Transcriber, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers an STT driver (called from a driver's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

// Drivers returns the registered driver names.
func Drivers() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}

func init() {
	togo.RegisterProviderFunc("ai-stt", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("STT_DRIVER")
		if name == "" {
			return nil // not configured yet — skip silently
		}
		regMu.RLock()
		f, ok := drivers[name]
		regMu.RUnlock()
		if !ok {
			return fmt.Errorf("ai-stt: unknown driver %q (set STT_DRIVER to one of %v)", name, Drivers())
		}
		t, err := f(k)
		if err != nil {
			return err
		}
		k.Set("ai-stt", &Service{driver: t, name: name})
		return nil
	})
}

// Service is the kernel-bound STT service.
type Service struct {
	driver Transcriber
	name   string
}

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.name }

// Transcribe runs the active driver.
func (s *Service) Transcribe(ctx context.Context, req Request) (Response, error) {
	return s.driver.Transcribe(ctx, req)
}

// FromKernel returns the STT service bound to the kernel.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("ai-stt")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}
