package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/libragen/felix/model"
)

const contextKeyWslogHookId = "wslogHookID"

func JwtMiddlewareWslog(c *gin.Context) {
	//ip := c.ClientIP()
	//log.Printf("hook api IP:%s", ip)
	//if ip == "::1" || ip == "127.0.0.1" {
	//	c.Set(contextKeyWslogHookId, uint(1))
	//	c.Next()
	//	return
	//}
	token, ok := c.GetQuery("_t")
	if !ok {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength:])
	}

	m, err := model.WslogHookCheckToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": fmt.Sprintf("_t token check error: %s", err)})
		return
	}

	//store the user Model in the context
	c.Set(contextKeyWslogHookId, m.Id)

	c.Next()
	// after request
}
