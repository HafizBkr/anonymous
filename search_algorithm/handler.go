package searchalgorithm

import (
	"net/http"

	"github.com/go-chi/render"
)

func SearchHandler(service SearchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Query parameter is required", http.StatusBadRequest)
			return
		}

		results, err := service.Search(query, 5, 0)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, results)
	}
}
