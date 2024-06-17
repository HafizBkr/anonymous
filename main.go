package main

import (
	"anonymous/auth"
	"anonymous/comments"
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
	commentRepo := comments.NewCommentRepo(postgresPool)

	authService := auth.Service(usersRepo, txProvider, logger, jwtProvider)
	userService := users.Service(usersRepo, txProvider, logger)
	postService := posts.NewPostService(postRepo, *authService )
	commentService := comments.NewCommentService(commentRepo, *authService )



	authHandler := auth.NewAuthHandler(authService, logger)
	userHandler := users.Handler(userService, logger)
	
	createPostHandler := posts.CreatePostHandler(postService)
	getAllPostsHandler :=posts.GetAllPostsHandler(postService)
	getPostByUserHAndler :=posts.GetPostsByUserHandler(postService)
	updatePostHAndler := posts.UpdatePostHandler(postService)
	deletePostHandler := posts.DeletePostHandler(postService)
	
	createCommentsHandler := comments.CreateCommentHandler(commentService)
	updateCommentHandler := comments.UpdateCommentHandler(commentService)
	getCommentByPostHandler := comments.GetCommentsByPostIDHandler(commentService)
	getCommentHandler := comments.GetCommentHandler(commentService)
	deleteCommentHandler := comments.DeleteCommentHandler(commentService)
	
	
	
	

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.HandleRegistration)
		r.Post("/login", authHandler.HandleLogin)
		r.Get("/verify-email", authHandler.HandleEmailVerification)
	})
	r.Route("/users", func(r chi.Router) {
		r.Patch("/password", userHandler.HandleChangePassword)
		r.Get("/", userHandler.HandleGetAllUsers)
		r.Patch("/status", userHandler.HandleToggleStatus)
	})
	
	r.Route("/posts", func(r chi.Router) {
			r.Use(authMiddleware.MiddlewareHandler)
			r.Post("/", createPostHandler)
			r.Get("/", getAllPostsHandler)
			r.Get("/user/{userID}", getPostByUserHAndler)
			r.Patch("/{postID}", updatePostHAndler)
			r.Delete("/{postID}", deletePostHandler)
		})
	
	r.Route("/{postID}/comments", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
				r.Post("/", createCommentsHandler)
				r.Get("/", getCommentByPostHandler)
				r.Get("/{commentID}", getCommentHandler)
				r.Patch("/{commentID}", updateCommentHandler)
				r.Delete("/{commentID}", deleteCommentHandler)
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
