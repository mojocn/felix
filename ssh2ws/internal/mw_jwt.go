package internal

import (
	"net/http"
	"strings"

	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
)

const contextKeyUserObj = "authedUserObj"
const bearerLength = len("Bearer ")

func ctxTokenToUser(c *gin.Context, roleId uint) {
	token, ok := c.GetQuery("_t")
	if !ok {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength:])
	}
	usr, err := model.JwtParseUser(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": err.Error()})
		return
	}
	if (usr.RoleId & roleId) != roleId {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "roleId 没有权限"})
		return
	}

	//store the user Model in the context
	c.Set(contextKeyUserObj, *usr)
	c.Next()
	// after request
}

func MwUserAdmin(c *gin.Context) {
	ctxTokenToUser(c, 2)
}

func MwUserComment(c *gin.Context) {
	ctxTokenToUser(c, 8)
}
