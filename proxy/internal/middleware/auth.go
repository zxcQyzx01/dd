package middleware

import (
    "context"
    "net/http"
    "strings"
    pb_auth "dd/pkg/auth"
)

func AuthMiddleware(authClient pb_auth.AuthServiceClient) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            resp, err := authClient.ValidateToken(context.Background(), &pb_auth.ValidateTokenRequest{
                Token: token,
            })

            if err != nil || !resp.Valid {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Добавляем user_id в контекст запроса
            ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
} 