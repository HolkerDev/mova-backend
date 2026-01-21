package middleware

import (
	"errors"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mova-backend/internal/service"
)

type AuthUser struct {
	ID      uuid.UUID
	ClerkID string
	Email   string
}

func ClerkAuth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := clerkhttp.RequireHeaderAuthorization()(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, ok := clerk.SessionClaimsFromContext(r.Context())
				if !ok {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
					return
				}

				user, err := userService.GetUserByClerkID(r.Context(), claims.Subject)
				if err != nil {
					if errors.Is(err, service.ErrUserNotFound) {
						Logger.Error("user not found", "clerk_id", claims.Subject)
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
						return
					}
					Logger.Error("failed to get user", "clerk_id", claims.Subject, "error", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
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
