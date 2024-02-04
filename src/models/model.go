package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"time"
)

// Token JWT claims struct
type Token struct {
	WalletId uuid.UUID
	jwt.StandardClaims
}

type Wallet struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique;not null"` // строковый ID кошелька, генерируется сервером
	Balance float64   `json:"balance" gorm:"not null"`                                                  // дробное число, баланс кошелька
	Token   string    `json:"token" sql:"-"`
}

type Transaction struct {
	ID     uint      `gorm:"primary_key"`                    // уникальный идентификатор транзакции
	From   uuid.UUID `json:"from" gorm:"type:uuid;not null"` // строковый ID исходящего кошелька
	To     uuid.UUID `json:"to" gorm:"type:uuid;not null"`   // строковый ID входящего кошелька
	Amount float64   `json:"amount" gorm:"not null"`         // дробное число, сумма перевода
	Time   time.Time `json:"time" gorm:"default:now()"`      // дата и время перевода в формате RFC 3339

}
