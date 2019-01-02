package middleware

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		bearer := r.Header.Get("Authorization")

		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			token = bearer[7:]
		}

		if user, err := findUserByToken(token); err == nil {
			//If the token is valid, save the user profile to the request context
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		fmt.Fprintf(w, `{"error" : "Invalid token", "code" : %d } `, http.StatusUnauthorized)
		return
	})
}

func findUserByToken(token string) (*User, error) {

	var activeUser *User

	//has to fetch use details from database
	for _, v := range users {
		if v.Token == token {
			return v, nil
		}
	}

	return activeUser, fmt.Errorf("User with token, %s not found", token)
}

*/
