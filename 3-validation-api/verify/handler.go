package verify

import (
	//"github.com/jordan-wright/email"
	"net/http"
)

type SendRequest struct {
	To string `json:"to"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Парсим email
	// 2. Генерим хэш
	// 3. Кладём в хранилище
	// 4. Отправляем email через jordan-wright/email
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем hash из URL
	// 2. Проверяем в storage
	// 3. Если есть — подтверждён, иначе — ошибка
}
