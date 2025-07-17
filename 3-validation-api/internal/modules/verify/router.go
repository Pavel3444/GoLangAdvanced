package verify

import (
	"net/http"
	"strings"
)

func SetupRoutes() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
		SendEmailHandler(w, r)
	})

	http.HandleFunc("/verify/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/verify/")
		if path == "" {
			http.Error(w, "Отсутствует хеш", http.StatusBadRequest)
			return
		}

		hash := strings.Split(path, "/")[0]

		VerifyHandler(w, r, hash)
	})
}
