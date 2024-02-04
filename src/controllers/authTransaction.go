package controllers

import (
	"EWallet/src/models"
	t "EWallet/src/tools"

	"github.com/google/uuid"

	"encoding/json"
	"net/http"
)

// CreateTransaction перевод средств с одного кошелька на другой
var CreateTransaction = func(w http.ResponseWriter, r *http.Request) {
	transaction := &models.Transaction{}
	err := json.NewDecoder(r.Body).Decode(transaction) //декодирует тело запроса (to, amount) в struct
	if err != nil {
		t.Respond(w, t.Message(false, "Invalid request"))
		return
	}

	// достаём из URL "from"
	id := r.Context().Value("walletId").(uuid.UUID)
	transaction.From = id
	if err != nil {
		t.Respond(w, t.Message(false, "Ошибка парсинга"))
		return
	}

	resp := transaction.Create()
	t.Respond(w, resp)
}

// GetTransactionHistory получение историй входящих и исходящих транзакций
var GetTransactionHistory = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("walletId").(uuid.UUID)

	in, out := models.GetTransactions(id)
	if in == nil && out == nil {
		resp := t.Message(false, "no data")
		t.Respond(w, resp)
	}
	resp := t.Message(true, "success")
	resp["in"] = in
	resp["out"] = out
	t.Respond(w, resp)
}
