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
	"github.com/wastewater-intelligence-network/win-api/utils"
)

var (
	tokenStrategy auth.Strategy
)

func AuthMiddleware(policy *utils.Policy) gin.HandlerFunc {
	strategy := setupGoGuardian(policy)
	return func(c *gin.Context) {
		if policy.IsOpen(c.Request.RequestURI) {
			c.Next()
			return
		}
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

func setupGoGuardian(policy *utils.Policy) union.Union {
	cache := NewCache(time.Minute * 30)
	basicStrategy := basic.NewCached(getValidationUserFunc(policy), cache)
	tokenStrategy = token.New(token.NoOpAuthenticate, cache)
	strategy := union.New(basicStrategy, tokenStrategy)

	return strategy
}

func getValidationUserFunc(policy *utils.Policy) basic.AuthenticateFunc {
	return func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
		// here connect to db or any other service to fetch user and validate it.
		if userName == "admin" && password == "admin" {
			return auth.NewUserInfo(
				"admin",
				"1",
				[]string{"transporter"},
				nil,
			), nil
		}

		return nil, fmt.Errorf("Invalid credentials")
	}
}

func PolicyMiddleware(policy *utils.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := c.Get("user")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user := u.(auth.Info)
		fmt.Println(user)
		if !policy.Check(c.Request.RequestURI, user) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
