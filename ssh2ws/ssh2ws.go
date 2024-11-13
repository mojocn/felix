package ssh2ws

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/felix/felixbin"
	"github.com/mojocn/felix/model"
	"github.com/mojocn/felix/ssh2ws/internal"
)

func RunSsh2ws(bindAddress, user, password, secret string, expire time.Duration, verbose bool) error {
	err := model.CreateGodUser(user, password)
	if err != nil {
		return err
	}
	//config jwt variables
	model.AppSecret = secret
	model.ExpireTime = expire
	model.AppIss = "felix.mojotv.cn"
	if !verbose {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.MaxMultipartMemory = 32 << 20
	//sever static file in http's root path
	binStaticMiddleware, err := felixbin.NewGinStaticBinMiddleware("/")
	if err != nil {
		return err
	}

	mwCORS := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 2400 * time.Hour,
	})
	r.Use(binStaticMiddleware, mwCORS)

	{
		r.POST("comment-login", internal.LoginCommenter)       //评论用户登陆
		r.POST("comment-register", internal.RegisterCommenter) //评论用户注册
	}

	api := r.Group("api")
	api.POST("admin-login", internal.LoginAdmin) //管理后台登陆
	api.GET("meta", internal.Meta)

	authG := api.Use(internal.MwUserAdmin)
	{

		//create wslog hook

		authG.GET("ssh", internal.SshAll)
		authG.POST("ssh", internal.SshCreate)
		authG.GET("ssh/:id", internal.SshOne)
		authG.PATCH("ssh", internal.SshUpdate)
		authG.DELETE("ssh/:id", internal.SshDelete)

		authG.GET("sftp/:id", internal.SftpLs)
		authG.GET("sftp/:id/dl", internal.SftpDl)
		authG.GET("sftp/:id/cat", internal.SftpCat)
		authG.GET("sftp/:id/rm", internal.SftpRm)
		authG.GET("sftp/:id/rename", internal.SftpRename)
		authG.GET("sftp/:id/mkdir", internal.SftpMkdir)
		authG.POST("sftp/:id/up", internal.SftpUp)

		authG.GET("user", internal.UserAll)
		authG.POST("user", internal.RegisterCommenter)
		authG.DELETE("user/:id", internal.UserDelete)
		authG.PATCH("user", internal.UserUpdate)

	}

	if err := r.Run(bindAddress); err != nil {
		return err
	}
	return nil
}
