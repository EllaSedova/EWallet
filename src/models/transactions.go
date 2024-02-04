package models

import (
	t "EWallet/src/tools"

	uuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Validate проверка транзакции
func (transaction *Transaction) Validate() (map[string]interface{}, bool) {
	// проверяем сумму перевода
	if transaction.Amount <= 0 {
		return t.Message(false, "Неверная сумма перевода"), false
	}

	// проверяем, что в таблице wallets присутствуют id кошельков "from" и "to"
	walletFrom := &Wallet{}
	GetDB()
	err := GetDB().Table("wallets").Where("id = ?", transaction.From).First(walletFrom).Error
	if err == gorm.ErrRecordNotFound {
		return t.Message(false, "Отправителя не существует"), false
	}

	walletTo := &Wallet{}
	err = GetDB().Table("wallets").Where("id = ?", transaction.To).First(walletTo).Error
	if err == gorm.ErrRecordNotFound {
		return t.Message(false, "Получателя не существует"), false
	}

	// проверяем баланс отправителя (balance >= amount)
	if walletFrom.Balance < transaction.Amount {
		return t.Message(false, "Не достаточно средств на счёте"), false
	}

	return t.Message(false, "Requirement passed"), true
}

// Create создаёт новую транзакцию
func (transaction *Transaction) Create() map[string]interface{} {
	if resp, ok := transaction.Validate(); !ok {
		return resp
	}

	GetDB().Create(transaction)

	if transaction.ID <= 0 {
		return t.Message(false, "Failed to create transaction, connection error.")
	}

	flag := WalletUpdate(transaction)
	if !flag {
		// удаляем запись о транзакции из таблицы transactions
		db.Delete(&transaction, transaction.ID)
		return t.Message(false, "Failed to create transaction. Try again.")
	}
	response := t.Message(true, "Transaction has been created")
	response["transaction"] = transaction
	return response
}

// GetTransactions получить историю транзакций
func GetTransactions(id uuid.UUID) ([]*Transaction, []*Transaction) {
	out := make([]*Transaction, 0)
	in := make([]*Transaction, 0)
	err1 := GetDB().Table("transactions").Where(&Transaction{From: id}).Find(&out).Error
	err2 := GetDB().Table("transactions").Where(&Transaction{To: id}).Find(&in).Error
	if err1 != nil && err2 != nil {
		return nil, nil
	}

	return in, out
}