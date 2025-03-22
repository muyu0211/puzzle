package cmd

import (
	"Puzzle/internal/controller/image"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gmode"
	"time"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 设置Session过期时间为1min
			s.SetSessionMaxAge(time.Minute * 30)

			// 开发阶段禁用浏览器缓存
			if gmode.IsDevelop() {
				s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					r.Response.Header().Set("Cache-Control", "no-store")
				})
			}

			// 实例化出Hello控制器对象
			//h := hello.NewHello()
			//s.BindHandler("/second", h.SayHello)

			// 批量绑定：将user控制器对象绑定之/user路由中
			//s.BindObject("/user", user.NewUser(), "AddUser, GetUser")

			/**
			// 分组路由
			s.Group("/user", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					user.NewUser(),
				)
			})

			s.Group("/hello", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewHello(),
				)
			})
			*/

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(ghttp.MiddlewareCORS)
				group.Bind(
					image.NewImage(),
				)

				//i := image.NewImage()
				//s.BindHandler("/puzzle", i.Index)
				//s.BindHandler("/start", i.StartPuzzle)
				//// 为 /upload 路由单独绑定自定义中间件
				////group.Middleware(middleware.ImageSplitValid)
				//group.POST("/upload", i.UploadImage)
			})

			s.SetServerRoot("resource/public")
			s.Run()
			return nil
		},
	}
)
