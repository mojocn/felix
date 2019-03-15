package internal

import (
	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Login(c *gin.Context) {
	var mdl model.User
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	ip := c.ClientIP()
	data, err := mdl.Login(ip)
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}

func Meta(c *gin.Context) {
	github := viper.GetString("github.client_id")
	githubCallbackUrl := viper.GetString("github.callback_url")
	jsonData(c, gin.H{"github_client_id": github, "github_callback_url": githubCallbackUrl})
}
