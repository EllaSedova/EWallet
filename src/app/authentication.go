package app

import (
	"EWallet/src/models"
	t "EWallet/src/tools"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"

	"context"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/v1/wallet", "/api/v1/wallet/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                     //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = t.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = t.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = t.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = t.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "walletId", tk.WalletId)
		r = r.WithContext(ctx)

		// проверяем совпадает ли id из URl и из JWT токена
		params := mux.Vars(r)
		idFromURL, _ := uuid.Parse(params["walletId"])
		if idFromURL != tk.WalletId {
			response = t.Message(false, "Токен не соответствует кошельку")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
