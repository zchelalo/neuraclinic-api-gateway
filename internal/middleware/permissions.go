package middleware

import (
	"fmt"
	"net/http"

	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

func PermissionsMiddleware(required ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, err := RequireAuth(r)
			if err != nil {
				response.WriteError(w, r, err)
				return
			}
			if len(required) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			allowed := make(map[string]struct{}, len(auth.PermissionsKeys))
			for _, key := range auth.PermissionsKeys {
				allowed[key] = struct{}{}
			}
			for _, permission := range required {
				if _, ok := allowed[permission]; !ok {
					response.WriteError(w, r, response.Forbidden("permission_denied", fmt.Sprintf("missing required permission %s", permission), nil))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
