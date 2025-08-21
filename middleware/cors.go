package middleware

import (
    "net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // origin := r.Header.Get("Origin")
        // allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
        
        // originAllowed := false
        // for _, allowed := range allowedOrigins {
        //     if origin == allowed {
        //         w.Header().Set("Access-Control-Allow-Origin", origin)
        //         originAllowed = true
        //         break
        //     }
        // }

        // if !originAllowed {
        //     http.Error(w, "Origin not allowed", http.StatusForbidden)
        //     return
        // }
        
        w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}