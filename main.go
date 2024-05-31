package main

import (
	"anonymous/postgres"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"anonymous/users"
	"anonymous/auth"
	"anonymous/provider"
)

func main() {
	godotenv.Load()
	r := chi.NewRouter()
	port := os.Getenv("PORT")
	databse_url := os.Getenv("DB_URL")
	postgresPool := postgres.GetconnectionPool(databse_url)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	jwtProvider := providers.NewJWTProvider()
	txProvider := providers.NewTransactionProvider(postgresPool)

	usersRepo := users.Repo(postgresPool)

	authService := auth.Service(usersRepo, txProvider, logger, jwtProvider)
	userService := users.Service(usersRepo, txProvider, logger)

	authHandler := auth.NewAuthHandler(authService, logger)
	userHandler := users.Handler(userService, logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.HandleRegistration)
		r.Post("/login", authHandler.HandleLogin)
		r.Get("/verify-email/verify", authHandler.HandleEmailVerification) // Gestionnaire pour le lien de vérification dans l'e-mail/ Ajoutez cette ligne pour gérer la vérification de l'e-mail
	})
	r.Route("/users", func(r chi.Router) {
		r.Patch("/password", userHandler.HandleChangePassword)
		r.Get("/", userHandler.HandleGetAllUsers)
		r.Patch("/status", userHandler.HandleToggleStatus)
	})

	server := http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Println("Server is running on port:", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
