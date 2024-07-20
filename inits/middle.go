package inits

import "github.com/kataras/iris/v12"

func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Expose-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Api, Accept, Authorization, Version, Token")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
