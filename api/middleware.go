package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/token"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorisationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorisationHeader) == 0 {
			err := errors.New("missing authorisation header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorisationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorisation header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorisationType := strings.ToLower(fields[0])
		if authorisationType != authorizationTypeBearer {
			err := errors.New("invalid authorisation type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}

}

func (server *Server) mustBeAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		userRoles, err := server.store.GetUserRoles(c, userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if len(userRoles) == 0 {
			c.Next()
		}
		c.Next()
	}

}

func (server *Server) currentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			//c.Redirect(http.StatusTemporaryRedirect, "/users/login")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("missing authorization header")))
			return
		}

		// Assuming your token format is "Bearer <token>"
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		claims, err := server.tokenMaker.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		c.Set("user_name", claims.UserName)
		args := db.GetUserByUserNameOrEmailParams{
			Username: claims.UserName,
		}

		// get all user attributes from db and set them in the context
		user, err := server.store.GetUserByUserNameOrEmail(c, args)
		if err != nil {
			msg := fmt.Sprintf("error getting user from db with username %s", claims.UserName)
			log.Error().Msg(msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		userRoles, err := server.store.GetUserRoles(c, user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		roleIDs := make([]int64, len(userRoles))
		for i, role := range userRoles {
			roleIDs[i] = role.RoleID
		}

		fmt.Println("user roles are ", userRoles)
		c.Set("roles", roleIDs)

		c.Set("user_id", user.ID)

		// log.Info().Msg("User %s is logged in  as %s his id is %d", claims.Username, user_roles, user.ID)

		c.Next()
	}
}
