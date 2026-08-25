package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

const maxBodyBytes = 8 * 1024

var e164Pattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

type messageRequest struct {
	To       string `json:"to"`
	Body     string `json:"body"`
	SenderID string `json:"sender_id,omitempty"`
}

type validationResponse struct {
	Valid            bool   `json:"valid"`
	To               string `json:"to"`
	BodyRunes        int    `json:"body_runes"`
	BodyBytes        int    `json:"body_bytes"`
	ProviderDispatch bool   `json:"provider_dispatch"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func validateMessage(input messageRequest) (validationResponse, error) {
	input.To = strings.TrimSpace(input.To)
	input.SenderID = strings.TrimSpace(input.SenderID)
	if !e164Pattern.MatchString(input.To) {
		return validationResponse{}, errors.New("to must use basic E.164 format")
	}
	if !utf8.ValidString(input.Body) {
		return validationResponse{}, errors.New("body must be valid UTF-8")
	}
	if strings.TrimSpace(input.Body) == "" {
		return validationResponse{}, errors.New("body must not be empty")
	}
	bodyRunes := utf8.RuneCountInString(input.Body)
	if bodyRunes > 1600 || len(input.Body) > maxBodyBytes {
		return validationResponse{}, errors.New("body exceeds validation limits")
	}
	if len(input.SenderID) > 32 {
		return validationResponse{}, errors.New("sender_id exceeds 32 bytes")
	}
	return validationResponse{
		Valid:            true,
		To:               input.To,
		BodyRunes:        bodyRunes,
		BodyBytes:        len(input.Body),
		ProviderDispatch: false,
	}, nil
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes+1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input messageRequest
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only one complete JSON object is allowed"})
		return
	}

	response, err := validateMessage(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "sky-sms-envelope"})
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", healthHandler)
	mux.HandleFunc("/v1/messages/validate", validateHandler)
	return mux
}

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Sky SMS Envelope listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
