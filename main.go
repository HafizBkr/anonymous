package main

import(
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"os"
	"net/http"
	"time"
	"net"
	"log"
)
func main (){ 
	godotenv.Load()
		r := chi.NewRouter()
		port := os.Getenv("PORT")
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
