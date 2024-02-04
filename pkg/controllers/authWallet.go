package controllers

import (
	"EWallet/pkg/manager"
	t "EWallet/pkg/tools"

	"github.com/google/uuid"

	"encoding/json"
	"net/http"
)

// CreateWallet создать кошелёк
var CreateWallet = func(w http.ResponseWriter, r *http.Request) {
	wallet := &manager.Wallet{}
	response := wallet.Create() //Создать кошелёк

	// проверяем response
	if response["status"].(int) == 400 {
		w.WriteHeader(http.StatusBadRequest)
		t.Respond(w, response)
	}
	w.WriteHeader(http.StatusOK)
	t.Respond(w, response)
}

// GetWalletBalance получение текущего состояния кошелька
var GetWalletBalance = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("walletId").(uuid.UUID)

	data := manager.GetWallet(id)
	if data == nil {
		resp := t.Message(404, "Указанный кошелек не найден")
		w.WriteHeader(http.StatusNotFound)
		t.Respond(w, resp)
	}
	resp := t.Message(200, "Данные кошелька получены")
	w.WriteHeader(http.StatusOK)
	resp["data"] = data
	t.Respond(w, resp)
}

// WalletLogin перегенерировать JWT для кошелька
var WalletLogin = func(w http.ResponseWriter, r *http.Request) {
	wallet := &manager.Wallet{}
	err := json.NewDecoder(r.Body).Decode(wallet) //декодирует тело запроса в struct
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Respond(w, t.Message(400, "Invalid request"))
		return
	}
	resp := manager.Login(wallet.ID)
	t.Respond(w, resp)
}
