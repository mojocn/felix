package ssh2ws

import (
	"time"

	"github.com/dejavuzhou/felix/felixbin"
	"github.com/dejavuzhou/felix/model"
	"github.com/dejavuzhou/felix/ssh2ws/internal"
	"github.com/dejavuzhou/felix/wslog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	mwJwt := internal.JwtMiddleware

	api := r.Group("api")
	api.POST("login", internal.Login)
	api.POST("register", internal.UserCreate)
	api.GET("meta", internal.Meta)

	//terminal log
	hub := wslog.NewHub()
	go hub.Run()

	{
		//websocket
		r.GET("ws/hook", mwJwt, internal.Wslog(hub))
		r.GET("ws/ssh/:id", mwJwt, internal.WsSsh)
	}
	//给外部调用
	{
		api.POST("wslog/hook-api", internal.JwtMiddlewareWslog, internal.WsLogHookApi(hub))
		api.GET("wslog/hook", mwJwt, internal.WslogHookAll)
		api.POST("wslog/hook", mwJwt, internal.WslogHookCreate)
		api.PATCH("wslog/hook", mwJwt, internal.WslogHookUpdate)
		api.DELETE("wslog/hook/:id", mwJwt, internal.WslogHookDelete)

		api.GET("wslog/msg", mwJwt, internal.WslogMsgAll)
		api.DELETE("wslog/msg/:id", mwJwt, internal.WslogMsgDelete)
	}

	//评论
	{
		api.GET("comment", internal.CommentAll)
		api.GET("comment/:id/:action", mwJwt, internal.CommentAction)
		api.POST("comment", mwJwt, internal.CommentCreate)
		api.DELETE("comment", mwJwt, internal.CommentDelete)
	}

	authG := api.Use(mwJwt)
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

		authG.POST("ginbro/gen", internal.GinbroGen)
		authG.POST("ginbro/db", internal.GinbroDb)
		authG.GET("ginbro/dl", internal.GinbroDownload)

		authG.GET("ssh-log", internal.SshLogAll)
		authG.DELETE("ssh-log/:id", internal.SshLogDelete)
		authG.PATCH("ssh-log", internal.SshLogUpdate)

		authG.GET("user", internal.UserAll)
		authG.POST("user", internal.UserCreate)
		//api.GET("user/:id", internal.SshAll)
		authG.DELETE("user/:id", internal.UserDelete)
		authG.PATCH("user", internal.UserUpdate)

	}

	if err := r.Run(bindAddress); err != nil {
		return err
	}
	return nil
}
