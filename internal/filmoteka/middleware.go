package filmoteka

import (
	"log/slog"
	"net/http"
	"strconv"
)

var (
	accessiblePathsForRegularUser = map[string]struct{}{
		"/movies":                 {},
		"/movies/search":          {},
		"/movies/search_by_actor": {},
		"/actors":                 {},
		"/actors/movies":          {},
	}
)

func (f *Filmoteka) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			f.Logger.Info("Response", slog.String("Unauthorized", strconv.Itoa(http.StatusUnauthorized)))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role := f.GetUserRole(token)

		if role != "admin" {
			if _, ok := accessiblePathsForRegularUser[r.URL.Path]; !ok {
				f.Logger.Info("Response", slog.String("Forbidden", strconv.Itoa(http.StatusForbidden)))
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (f *Filmoteka) GetUserRole(token string) string {
	if token == f.Config.Auth.Admin {
		return "admin"
	}
	return "user"
}
