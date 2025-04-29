package request

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authHeader := r.Header.Get("Authorization")
		//if authHeader == "" {
		//	http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
		//	return
		//}
		//
		//tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		//if tokenStr == authHeader {
		//	http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
		//	return
		//}
		//
		//token, err := jwt.Parse([]byte(tokenStr), jwt.WithKeySet(CachedJWK))
		//if err != nil {
		//	http.Error(w, "Неверный токен", http.StatusUnauthorized)
		//	return
		//}
		//
		//sub, ok := token.Get("sub")
		//if !ok {
		//	http.Error(w, "Отсутствует sub в токене", http.StatusUnauthorized)
		//	return
		//}
		//
		//ctx := context.WithValue(r.Context(), "userID", sub)
		//next.ServeHTTP(w, r.WithContext(ctx))
	})
}
