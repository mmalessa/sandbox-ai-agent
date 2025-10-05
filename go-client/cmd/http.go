package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"go-client/lib/aiclient"
	"go-client/lib/appconfig"
	"go-client/lib/httptools"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run HTTP server",
	Run:   cmd_http,
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func cmd_http(cmd *cobra.Command, args []string) {
	sessionId := uuid.NewString()

	ai := aiclient.New(cfgFile, sessionId, chatName)

	_ = ai

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/ask", func(w http.ResponseWriter, r *http.Request) {
			// Limit request body size
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB

			var req httptools.RequestData
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&req); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid JSON: %v", err)})
				return
			}
			if req.Content == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "content is required"})
				return
			}

			w.Header().Set("Content-Type", "application/json")

			promptBuilder := aiclient.PromptBuilder()
			prompt := promptBuilder.WithTask(req.Content).Get()

			// Run AI call with timeout and return structured errors
			ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
			defer cancel()
			type aiResult struct {
				resp string
				err error
			}
			resultCh := make(chan aiResult, 1)
			go func() {
				resp, err := ai.Ask(prompt)
				resultCh <- aiResult{resp: resp, err: err}
			}()

			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				json.NewEncoder(w).Encode(map[string]string{"error": "AI request timed out"})
				return
			case res := <-resultCh:
				if res.err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": res.err.Error()})
					return
				}
				json.NewEncoder(w).Encode(httptools.ResponseData{Content: res.resp})
			}
		})
	})

	httpPort := appconfig.AppCfg.AiChatCfg[chatName].TmpHttpPort
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	// Start server and handle graceful shutdown
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- srv.ListenAndServe()
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	case <-shutdownCtx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("graceful shutdown error:", err)
		}
	}

}
