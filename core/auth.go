package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

var (
	tokenStrategy auth.Strategy
)

func AuthMiddleware() gin.HandlerFunc {
	strategy := setupGoGuardian()
	return func(c *gin.Context) {
		_, user, err := strategy.AuthenticateRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "User Unauthorised",
				"error":   err.Error(),
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func setupGoGuardian() union.Union {
	cache := NewCache(time.Minute * 30)
	basicStrategy := basic.NewCached(validateUser, cache)
	tokenStrategy = token.New(token.NoOpAuthenticate, cache)
	strategy := union.New(basicStrategy, tokenStrategy)

	return strategy
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "admin" && password == "admin" {
		return auth.NewUserInfo("admin", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}
