package main

import (
	"anonymous/auth"
	"anonymous/chat"
	"anonymous/comments"
	"anonymous/communitychats"
	"anonymous/comunauter"
	middlewares "anonymous/middleware"
	"anonymous/points"
	"anonymous/postgres"
	"anonymous/posts"
	providers "anonymous/provider"
	"anonymous/replies"
	searchalgorithm "anonymous/search_algorithm"
	"anonymous/users"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	repliesRepo := replies.NewCommentReplyRepo(postgresPool)
	comunityRepo := comunauter.NewCommunityRepo(postgresPool)
	pointsRepo := points.NewPointsRepo(postgresPool)
	communityChatRepo := communitychats.NewCommunityChatRepo(postgresPool)

	authService := auth.Service(usersRepo, txProvider, logger, jwtProvider)
	userService := users.Service(usersRepo, txProvider, logger)
	postService := posts.NewPostService(postRepo, *authService)
	commentService := comments.NewCommentService(commentRepo, *authService)
	repliesService := replies.NewCommentReplyService(repliesRepo, *authService)
	comunityService := comunauter.NewCommunityService(comunityRepo, *authService)
	pointService := points.NewPointsService(pointsRepo, logger, jwtProvider)
	communityChatService := communitychats.NewCommunityChatService(communityChatRepo, authService)

	authHandler := auth.NewAuthHandler(authService, logger)
	userHandler := users.Handler(userService, logger)
	pointHandler := points.NewPointsHandler(pointService, logger)

	createPostHandler := posts.CreatePostHandler(postService)
	getAllPostsHandler := posts.GetAllPostsHandler(postService)
	getPostByUserHAndler := posts.GetPostsByUserHandler(postService)
	updatePostHAndler := posts.UpdatePostHandler(postService)
	deletePostHandler := posts.DeletePostHandler(postService)
	LikePostHandler := posts.LikePostHandler(postService)
	UnlikePostHandler := posts.UnlikePostHandler(postService)
	AddReactionHandler := posts.AddReactionHandler(postService)
	RemoveReactionHandler := posts.RemoveReactionHandler(postService)

	createCommentsHandler := comments.CreateCommentHandler(commentService)
	updateCommentHandler := comments.UpdateCommentHandler(commentService)
	getCommentByPostHandler := comments.GetCommentsByPostIDHandler(commentService)
	getCommentHandler := comments.GetCommentHandler(commentService)
	deleteCommentHandler := comments.DeleteCommentHandler(commentService)

	createCommentReplyHandler := replies.CreateCommentReplyHandler(repliesService)
	getCommentRepliesHandler := replies.GetCommentReplyHandler(repliesService)
	getCommentRepliesByCommentIDHandler := replies.GetCommentRepliesByCommentIDHandler(repliesService)
	updateCommentReplyHandler := replies.UpdateCommentReplyHandler(repliesService)
	deleteCommentReplyHandler := replies.DeleteCommentReplyHandler(repliesService)

	createComunityHandler := comunauter.CreateCommunityHandler(comunityService)
	joinComunityHandler := comunauter.JoinCommunityHandler(comunityService)
	getComunityHandler := comunauter.GetCommunityHandler(comunityService)
	getallComunityHandler := comunauter.GetAllCommunitiesHandler(comunityService)
	getAllUserComunity := comunauter.GetCommunityMembersHandler(comunityService)
	communityChatHandler := communitychats.NewCommunityChatHandler(*communityChatService)
	getconversationHandler := chat.GetConversationsHandler(postgresPool)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.HandleRegistration)
		r.Post("/login", authHandler.HandleLogin)
		r.Get("/verify-email", authHandler.HandleEmailVerification)
		r.Post("/forgot-password", authHandler.HandleForgotPassword)  
		r.Post("/reset-password", authHandler.HandleResetPassword)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Get("/", userHandler.HandleGetAllUsers)
		r.Patch("/status", userHandler.HandleToggleStatus)
		r.Patch("/password", userHandler.HandleChangePassword)
		r.Get("/{userID}", userHandler.HandleGetUserByID)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/", createPostHandler)
		r.Get("/", getAllPostsHandler)
		r.Get("/user/{userID}", getPostByUserHAndler)
		r.Patch("/{postID}", updatePostHAndler)
		r.Delete("/{postID}", deletePostHandler)
		r.Post("/{postID}/like", LikePostHandler)
		r.Delete("/{postID}/like", UnlikePostHandler)
		r.Post("/{postID}/reaction", AddReactionHandler)
		r.Delete("/{postID}/reaction", RemoveReactionHandler)
	})

	r.Route("/{postID}/comments", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/", createCommentsHandler)
		r.Get("/", getCommentByPostHandler)
		r.Get("/{commentID}", getCommentHandler)
		r.Patch("/{commentID}", updateCommentHandler)
		r.Delete("/{commentID}", deleteCommentHandler)
		r.Put("/{commentID}/reactions", comments.AddOrUpdateReactionHandler(commentService))
		r.Get("/{commentID}/reaction_counts", comments.GetReactionCountsHandler(commentService))

	})

	r.Route("/{commentID}/replies", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/", createCommentReplyHandler)
		r.Get("/", getCommentRepliesByCommentIDHandler)
		r.Get("/{replyID}", getCommentRepliesHandler)
		r.Patch("/{replyID}", updateCommentReplyHandler)
		r.Delete("/{replyID}", deleteCommentReplyHandler)
	})

	r.Route("/comunity", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/", createComunityHandler)
		r.Get("/", getallComunityHandler)
		r.Get("/{communityID}", getComunityHandler)
		r.Get("/u/{communityID}", getComunityHandler)
		r.Post("/{communityID}", joinComunityHandler)
		r.Get("/menbers/{communityID}", getAllUserComunity)

	})

	r.Get("/posts/{postID}/likes/count", posts.GetLikesCountHandler(postgresPool))
	r.Get("/posts/{postID}/reactions/count", posts.GetReactionsCountHandler(postgresPool))
	r.Get("/posts/{postID}/comments/count", comments.GetCommentsCountByPostIDHandler(commentService))

	r.With(authMiddleware.MiddlewareHandler).Get("/search", searchalgorithm.SearchHandler(searchalgorithm.NewSearchService(postgresPool)))

	r.Route("/chat", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Get("/conversations", getconversationHandler)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			chat.HandleHTTPMessage(postgresPool, w, r)
		})
		r.Handle("/ws", authMiddleware.MiddlewareHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chat.HandleWebSocket(postgresPool, w, r)
		})))
		r.Get("/messages/{user1ID}/{user2ID}", chat.GetMessagesBetweenUsersHandler(postgresPool))
		r.Get("/messages/owner", chat.GetMessagesByOwnerHandler(postgresPool))
		r.Patch("/messages/{messageID}", func(w http.ResponseWriter, r *http.Request) {
			chat.UpdateMessageHandler(postgresPool, w, r)
		})
		r.Delete("/messages/{messageID}", func(w http.ResponseWriter, r *http.Request) {
			chat.DeleteMessageHandler(postgresPool, w, r)
		})

	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	r.Route("/points", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/", pointHandler.HandleLikeUserProfile)
		r.Get("/{userID}", pointHandler.HandleGetUserProfileLikes)
	})
	r.Route("/community_chats", func(r chi.Router) {
		r.Use(authMiddleware.MiddlewareHandler)
		r.Post("/{communityID}/messages", communityChatHandler.SendMessage)
		r.Get("/{communityID}/messages", communityChatHandler.GetMessages)
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
