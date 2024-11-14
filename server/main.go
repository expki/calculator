package main

import (
	"calculator/config"
	"calculator/game"
	"calculator/public"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

func main() {
	appCtx, stopApp := context.WithCancel(context.Background())

	// Load config
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", os.Args[0], " <config.json>")
	}
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		err = config.CreateSample(os.Args[1])
		if err != nil {
			log.Fatalf("CreateSample: %v", err)
		}
	}
	configRaw, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("ReadFile %q: %v", os.Args[1], err)
	}
	cfg, err := config.ParseConfig(configRaw)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	cfg.TLS, err = cfg.TLS.Configurate()
	if err != nil {
		log.Fatalf("Configurate: %v", err)
	}

	// Logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	defer sugar.Sync()

	// Game
	engine := game.New(logger)

	// Create mux
	mux := http.NewServeMux()

	// HTTP
	server := http.Server{
		Handler: mux,
		Addr:    cfg.Server.HttpAddress,
	}

	// HTTP2
	server2 := http.Server{
		Handler: mux,
		Addr:    cfg.Server.HttpsAddress,
		TLSConfig: &tls.Config{
			GetCertificate: cfg.TLS.GetCertificate,
			ClientAuth:     tls.NoClientCert,
			NextProtos:     []string{"h2", "http/1.1"}, // Enable HTTP/2
		},
	}
	err = http2.ConfigureServer(&server2, &http2.Server{})
	if err != nil {
		sugar.Fatalf("http2.Server: %v", err)
	}

	// Headers middleware
	middlewareHeaders := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// WASM headers
			w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
			w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
			h.ServeHTTP(w, r)
		})
	}

	// Compression middleware
	middlewareCompression := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "zstd") {
				h.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Encoding", "zstd")
			encoder, err := zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedFastest))
			if err != nil {
				sugar.Errorf("Failed to create zstd encoder: %v", err)
				h.ServeHTTP(w, r)
				return
			}
			defer encoder.Close()
			zstrw := &zstdResponseWriter{ResponseWriter: w, Writer: encoder}
			h.ServeHTTP(zstrw, r)
		})
	}

	// Routes
	mux.Handle("/", middlewareHeaders(middlewareCompression(http.FileServerFS(public.Root))))
	mux.Handle("/session", middlewareHeaders(engine.UpgradeHandler()))

	// Start servers
	serverDone := make(chan struct{})
	go func() {
		sugar.Infof("HTTP server starting on %s", cfg.Server.HttpAddress)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			sugar.Errorf("ListenAndServe http: %v", err)
		}
		close(serverDone)
	}()
	server2Done := make(chan struct{})
	go func() {
		sugar.Infof("HTTP2 server starting on %s", cfg.Server.HttpsAddress)
		err := server2.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			sugar.Errorf("ListenAndServe https (http2): %v", err)
		}
		close(server2Done)
	}()

	// Interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for servers to finish
	select {
	case <-interrupt:
		sugar.Info("Interrupt signal received")
	case <-appCtx.Done():
		sugar.Info("App stopped")
	case <-serverDone:
		sugar.Info("HTTP server stopped")
	case <-server2Done:
		sugar.Info("HTTP2 server stopped")
	}
	sugar.Info("Server shutting down")
	shutdownCtx, cancelShutdown := context.WithTimeout(appCtx, 3*time.Second)
	defer cancelShutdown()
	server.Shutdown(shutdownCtx)
	server2.Shutdown(shutdownCtx)
	server.Close()
	server2.Close()
	stopApp()
	sugar.Info("Server stopped")
}

// zstdResponseWriter wraps the http.ResponseWriter to provide zstd compression
type zstdResponseWriter struct {
	http.ResponseWriter
	Writer *zstd.Encoder
}

func (w *zstdResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
