package comunauter

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	 "github.com/google/uuid"
)

func JoinCommunityHandler(service CommunityService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        communityID := chi.URLParam(r, "communityID")
        if _, err := uuid.Parse(communityID); err != nil {
            http.Error(w, "Invalid community ID", http.StatusBadRequest)
            return
        }

        err := service.JoinCommunity(token, communityID)
        if err != nil {
            if err.Error() == "unauthorized" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}


func CreateCommunityHandler(service CommunityService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var payload CommunityPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        community, err := service.CreateCommunity(&payload, token)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        jsonResponse(w, http.StatusCreated, community)
    }
}


func GetCommunityHandler(service CommunityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID := chi.URLParam(r, "communityID")
		uuid, err := uuid.Parse(communityID)
		if err != nil {
			http.Error(w, "Invalid community ID", http.StatusBadRequest)
			return
		}

		community, err := service.GetCommunity(uuid.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		jsonResponse(w, http.StatusOK, community)
	}
}




func GetAllCommunitiesHandler(service CommunityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communities, err := service.GetAllCommunities()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, communities)
	}
}

func GetCommunityMembersHandler(service CommunityService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        communityID := chi.URLParam(r, "communityID")
        if _, err := uuid.Parse(communityID); err != nil {
            http.Error(w, "Invalid community ID", http.StatusBadRequest)
            return
        }
        members, err := service.GetCommunityUsers(communityID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        jsonResponse(w, http.StatusOK, members)
    }
}


func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}



