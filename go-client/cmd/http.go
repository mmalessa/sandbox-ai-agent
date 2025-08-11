package cmd

import (
	"encoding/json"
	"fmt"
	"go-client/lib/aiclient"
	"go-client/lib/httptools"
	"net/http"

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

	ai := aiclient.New(cfgFile, sessionId)

	_ = ai

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/ask", func(w http.ResponseWriter, r *http.Request) {
			var req httptools.RequestData
			var resp httptools.ResponseData

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, fmt.Sprintf("Invalid JSON error: %#v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			response, err := ai.Ask(req.Content)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp.Content = response

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, "Encoding error", http.StatusInternalServerError)
			}
		})
	})
	http.ListenAndServe(":3000", r)

}
