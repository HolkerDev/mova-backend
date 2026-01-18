package middleware

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

func ClerkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := clerkhttp.RequireHeaderAuthorization()(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, ok := clerk.SessionClaimsFromContext(r.Context())
				if !ok {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
					return
				}

				c.Set("clerk_user_id", claims.Subject)
				c.Request = r
				c.Next()
			}),
		)

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func GetClerkUserID(c *gin.Context) string {
	if userID, exists := c.Get("clerk_user_id"); exists {
		return userID.(string)
	}
	return ""
}
