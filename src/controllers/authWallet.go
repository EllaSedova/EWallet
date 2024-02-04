package controllers

import (
	"EWallet/src/models"
	t "EWallet/src/tools"

	"github.com/google/uuid"

	"encoding/json"
	"fmt"
	"net/http"
)

// CreateWallet создать кошелёк
var CreateWallet = func(w http.ResponseWriter, r *http.Request) {
	wallet := &models.Wallet{}
	resp := wallet.Create() //Создать кошелёк
	t.Respond(w, resp)
}

// GetWalletBalance получение текущего состояния кошелька
var GetWalletBalance = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("walletId").(uuid.UUID)

	fmt.Println("ID", id)
	data := models.GetWallet(id)
	if data == nil {
		resp := t.Message(false, "no data")
		t.Respond(w, resp)
	}
	resp := t.Message(true, "success")
	resp["data"] = data
	t.Respond(w, resp)
}

// WalletLogin перегенерировать JWT для кошелька
var WalletLogin = func(w http.ResponseWriter, r *http.Request) {
	wallet := &models.Wallet{}
	err := json.NewDecoder(r.Body).Decode(wallet) //декодирует тело запроса в struct
	if err != nil {
		t.Respond(w, t.Message(false, "Invalid request"))
		return
	}
	resp := models.Login(wallet.ID)
	t.Respond(w, resp)
}

//var PutMoney = func(w http.ResponseWriter, r *http.Request) {
//
//	wallet := &models.Wallet{}
//	err := json.NewDecoder(r.Body).Decode(wallet) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
//	if err != nil {
//		t.Respond(w, t.Message(false, "Invalid request"))
//		return
//	}
//	// todo: сделать метод для изменения кошелька
//	//resp := wallet.Create() //Создать кошелёк
//	t.Respond(w, resp)
//}
