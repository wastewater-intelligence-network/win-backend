package core

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/wastewater-intelligence-network/win-api/db"
	"github.com/wastewater-intelligence-network/win-api/model"
	"github.com/wastewater-intelligence-network/win-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	TokenExpiryHours = time.Hour * 168
)

var (
	tokenStrategy auth.Strategy
)

func AuthMiddleware(policy *utils.Policy, conn *db.DBConnection) gin.HandlerFunc {
	strategy := setupGoGuardian(policy, conn)
	return func(c *gin.Context) {
		if policy.IsOpen(c.Request.RequestURI) {
			fmt.Println("Open")
			c.Next()
			return
		}
		_, user, err := strategy.AuthenticateRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "User Unauthorised",
				"error":   err.Error(),
				"status":  401,
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func setupGoGuardian(policy *utils.Policy, conn *db.DBConnection) union.Union {
	cache := NewCache(TokenExpiryHours)
	basicStrategy := basic.NewCached(getValidationUserFunc(policy, conn), cache)
	tokenStrategy = token.New(token.NoOpAuthenticate, cache)
	strategy := union.New(basicStrategy, tokenStrategy)

	return strategy
}

func getValidationUserFunc(policy *utils.Policy, conn *db.DBConnection) basic.AuthenticateFunc {
	return func(ctx context.Context, r *http.Request, username, password string) (auth.Info, error) {
		res := conn.FindOne(WIN_COLLECTION_USERS, bson.M{
			"username": username,
		})

		var user model.User
		err := res.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("Invalid credentials")
		}

		passwordHash := sha1.Sum([]byte(password))
		hashString := fmt.Sprintf("%x", passwordHash)

		fmt.Println(user.Hash)
		fmt.Println(hashString)
		fmt.Println(user)

		if user.Hash == hashString {
			return auth.NewUserInfo(
				user.Username,
				user.ID.String(),
				user.Roles,
				nil,
			), nil
		}

		return nil, fmt.Errorf("Invalid credentials")
	}
}

func PolicyMiddleware(policy *utils.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		if policy.IsOpen(c.Request.RequestURI) {
			c.Next()
			return
		}

		u, ok := c.Get("user")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user := u.(auth.Info)
		if !policy.Check(c.Request.RequestURI, user) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
