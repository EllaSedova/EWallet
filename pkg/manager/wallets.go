// Создание и аутентификация кошелька

package manager

import (
	t "EWallet/pkg/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"fmt"
	"os"
)

// Create создаёт новый кошелёк и генерирует токен JWT
func (wallet *Wallet) Create() map[string]interface{} {

	wallet.Balance = 100.0

	GetDB().Create(wallet)

	if len(wallet.ID.String()) <= 0 {
		return t.Message(400, "Не удалось создать кошелёк")
	}

	// создаём JWT токен для нового кошелька
	tokenString := GenerateJWT(wallet.ID)
	wallet.Token = tokenString
	fmt.Println("token: ", tokenString)

	response := t.Message(200, "Кошелёк создан")
	response["wallet"] = wallet
	return response
}

// WalletUpdate обновляет балансы двух кошельков для совершения транзакции
func WalletUpdate(transaction *Transaction) bool {

	// проверяем, что в таблице wallets присутствуют id кошельков "from" и "to"
	walletFrom := &Wallet{}
	err := GetDB().Table("wallets").Where("id = ?", transaction.From).First(walletFrom).Error
	if err != nil {
		return false
	}

	walletTo := &Wallet{}
	err = GetDB().Table("wallets").Where("id = ?", transaction.To).First(walletTo).Error
	if err != nil {
		return false
	}

	// обновляем балансы кошельков
	walletFrom.Balance -= transaction.Amount
	walletTo.Balance += transaction.Amount

	db.Model(&walletFrom).Update("balance", walletFrom.Balance)
	db.Model(&walletTo).Update("balance", walletTo.Balance)
	return true
}

// GetWallet получить данные о кошельке по ID
func GetWallet(id uuid.UUID) *Wallet {

	wallet := &Wallet{}
	err := GetDB().Table("wallets").Where("id = ?", id).First(wallet).Error
	if err != nil {
		return nil
	}

	return wallet
}

// Login обновить JWT токен для кошелька
func Login(id uuid.UUID) map[string]interface{} {

	wallet := &Wallet{}
	err := GetDB().Table("wallets").Where("id = ?", id).First(wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return t.Message(400, "Кошелёк не найден")
		}
		return t.Message(404, "Ошибка подключения")
	}

	// перегенерируем JWT токен
	tokenString := GenerateJWT(wallet.ID)
	wallet.Token = tokenString
	fmt.Println("token: ", tokenString)

	response := t.Message(200, "JWT токен обновлён")
	response["wallet"] = wallet
	return response
}

func GenerateJWT(walletID uuid.UUID) string {

	tk := &Token{WalletId: walletID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	return tokenString
}
