package middleware

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mova-backend/internal/database"
)

type AuthUser struct {
	ID      uuid.UUID
	ClerkID string
	Email   string
}

func ClerkAuth(queries *database.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := clerkhttp.RequireHeaderAuthorization()(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, ok := clerk.SessionClaimsFromContext(r.Context())
				if !ok {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
					return
				}

				user, err := queries.GetUserByClerkID(r.Context(), claims.Subject)
				if err != nil {
					Logger.Error("user not found", "clerk_id", claims.Subject, "error", err)
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
					return
				}

				c.Set("auth_user", AuthUser{
					ID:      user.ID,
					ClerkID: user.ClerkID,
					Email:   user.Email,
				})
				c.Request = r
				c.Next()
			}),
		)

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func GetAuthUser(c *gin.Context) (AuthUser, bool) {
	if user, exists := c.Get("auth_user"); exists {
		return user.(AuthUser), true
	}
	return AuthUser{}, false
}
