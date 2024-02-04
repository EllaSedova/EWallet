package main

import (
	"EWallet/pkg/controllers"
	"EWallet/pkg/middleware"
	"EWallet/pkg/tools"

	"github.com/gorilla/mux"

	"fmt"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/wallet", controllers.CreateWallet).Methods("POST")
	router.HandleFunc("/api/v1/wallet/login", controllers.WalletLogin).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/send", controllers.CreateTransaction).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/history", controllers.GetTransactionHistory).Methods("GET")
	router.HandleFunc("/api/v1/wallet/{walletId}", controllers.GetWalletBalance).Methods("GET")

	router.Use(middleware.JwtAuthentication) //attach JWT auth middleware

	router.NotFoundHandler = http.HandlerFunc(tools.NotFoundHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the middleware, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
