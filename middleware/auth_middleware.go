package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/zeeshanahmad0201/todo-mongo/helpers"
)

func Authorization() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")

		if clientToken == "" {
			http.Error(w, "No authorization header provided", http.StatusInternalServerError)
			return
		}

		claims, err := helpers.ValidateToken(clientToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"email": claims.Email,
			"name":  claims.Name,
			"uid":   claims.UserId,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
