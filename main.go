package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/ofonimefrancis/safepromo/config"
	"github.com/ofonimefrancis/safepromo/event"
)

// Routes returns a chi router instance which includes all routes needed for the application to run
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
		middleware.Timeout(60*time.Second),            // Timeout requests after 60 seconds
	)
	chiCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	})
	router.Use(chiCors.Handler)

	router.Mount("/api/events", event.Routes()) // Mount Golang Program debug/profiling route

	return router
}

func main() {
	config.Init()

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("👉 %s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("⚠️  Logging err: %s\n", err.Error())
	}

	// This block of code will allow graceful shutdown of our server, using the server Shutdown method which is a part of the standard library
	PORT := ":" + config.Get().Port
	server := http.Server{
		Addr:    PORT,
		Handler: router,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("😔 Shutting down. Goodbye..")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("⚠️  HTTP server Shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()
	log.Printf("Serving at 🔥 %s \n", PORT)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("⚠️  HTTP server ListenAndServe error: %v", err)
	}

	<-idleConnsClosed
}
