package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pajri/personal-backend/domain"
)

var excludedFromAuth = []string{
	"/api/auth/login",
	"/api/auth/signup",
	"/api/auth/verify_email",
	"/api/auth/reset_password/",
	"/api/auth/change_password",
}

func Middleware(authUseCase domain.IAuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := handleAuth(c, authUseCase)
		if !next {
			c.Abort()
			return
		}
		c.Next()
	}
}
