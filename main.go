package main

import (
	_ "Puzzle/internal/logic"
	_ "Puzzle/internal/packed"

	"Puzzle/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {

	cmd.Main.Run(gctx.GetInitCtx())
	//s := g.Server()
	//s.BindHandler("/", func(r *ghttp.Request) {
	//	r.Response.Write("Hello World!",
	//		r.Get("name").String(),
	//		r.Get("age").Int())
	//})
	//s.SetPort(8000)
	//s.Run()
}
