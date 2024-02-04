package middleware

import (
	"EWallet/pkg/manager"
	t "EWallet/pkg/tools"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"context"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/v1/wallet", "/api/v1/wallet/login"} // список эндпоинтов, не требующих аутентификации
		requestPath := r.URL.Path

		// проверяем нужна ли аутентификация
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // достаём токен

		if tokenHeader == "" { // нет токена
			response = t.Message(403, "Отсутствует токен аутентификации")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // токен приходит в формате `Bearer {token-body}` и его надо обработать
		if len(splitted) != 2 {
			response = t.Message(403, "Недействительный/неправильно сформированный токен аутентификации")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		tokenPart := splitted[1] // берём нужную часть токена
		tk := &manager.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = t.Message(403, "Неправильно сформированный токен аутентификации")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		if !token.Valid {
			response = t.Message(403, "Токен недействителен.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "walletId", tk.WalletId)
		r = r.WithContext(ctx)

		// проверяем совпадает ли id из URl и из JWT токена
		params := mux.Vars(r)
		idFromURL, _ := uuid.Parse(params["walletId"])
		if idFromURL != tk.WalletId {
			response = t.Message(403, "Токен не соответствует кошельку")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			t.Respond(w, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
