package middleware

import "github.com/gin-gonic/gin"

func CorsMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
		return
	}
	ctx.Next()
}
