package controllers

import (
	"EWallet/pkg/manager"
	t "EWallet/pkg/tools"

	"github.com/google/uuid"

	"encoding/json"
	"net/http"
)

// CreateTransaction перевод средств с одного кошелька на другой
var CreateTransaction = func(w http.ResponseWriter, r *http.Request) {

	transaction := &manager.Transaction{}
	err := json.NewDecoder(r.Body).Decode(transaction) //декодирует тело запроса (to, amount) в struct
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Respond(w, t.Message(400, "Ошибка в пользовательском запросе или ошибка перевода"))
		return
	}

	// достаём из URL "from"
	id := r.Context().Value("walletId").(uuid.UUID)
	transaction.From = id
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Respond(w, t.Message(400, "Ошибка в пользовательском запросе или ошибка перевода"))
		return
	}

	resp := transaction.Create()
	if resp["status"].(int) == 400 {
		w.WriteHeader(http.StatusBadRequest)
		t.Respond(w, resp)
	}
	if resp["status"].(int) == 404 {
		w.WriteHeader(http.StatusNotFound)
		t.Respond(w, resp)
	}
	w.WriteHeader(http.StatusOK)
	t.Respond(w, resp)
}

// GetTransactionHistory получение историй входящих и исходящих транзакций
var GetTransactionHistory = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("walletId").(uuid.UUID)

	in, out := manager.GetTransactions(id)
	if in == nil && out == nil {
		w.WriteHeader(http.StatusNotFound)
		resp := t.Message(404, "Указанный кошелек не найден")
		t.Respond(w, resp)
	}
	resp := t.Message(200, "История транзакций получена")
	w.WriteHeader(http.StatusOK)
	resp["in"] = in
	resp["out"] = out
	t.Respond(w, resp)
}
