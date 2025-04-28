package middleware

import "net/http"

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//origin := r.Header.Get("Origin")
		//// Разрешаем только определённые origin (например, localhost:3000)
		//allowedOrigins := map[string]bool{
		//	"http://localhost:3000": true,
		//	"http://127.0.0.1:3000": true,
		//}
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // если нужны куки

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
