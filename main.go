package main

import (
	"anonymous/auth"
	"anonymous/middleware"
	"anonymous/postgres"
	"anonymous/posts"
	"anonymous/provider"
	"anonymous/users"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	usersRepo := users.Repo(postgresPool)
	authMiddleware := middlewares.NewAuthMiddleware(usersRepo, jwtProvider, logger)
	postRepo := posts.NewPostRepo(postgresPool)

	authService := auth.Service(usersRepo, txProvider, logger, jwtProvider)
	userService := users.Service(usersRepo, txProvider, logger)
	postService := posts.NewPostService(postRepo , *authService )



	authHandler := auth.NewAuthHandler(authService, logger)
	userHandler := users.Handler(userService, logger)
	postHandler := posts.CreatePostHandler(postService)


	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.HandleRegistration)
		r.Post("/login", authHandler.HandleLogin)
		r.Get("/verify-email", authHandler.HandleEmailVerification) // Gestionnaire pour le lien de vérification dans l'e-mail/ Ajoutez cette ligne pour gérer la vérification de l'e-mail
	})
	r.Route("/users", func(r chi.Router) {
		r.Patch("/password", userHandler.HandleChangePassword)
		r.Get("/", userHandler.HandleGetAllUsers)
		r.Patch("/status", userHandler.HandleToggleStatus)
	})
	
 r.Route("/posts", func(r chi.Router) {
        r.Use(authMiddleware.MiddlewareHandler)
        r.Post("/", postHandler)
        // Ajoute d'autres routes liées aux posts ici
    })
	
	staticDir := "./static"
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))


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
