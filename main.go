package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"wechat-pay/wechat_pay"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("fe/*html")
	r.Static("/static", "./fe/dist")
	r.Any("/ping", func(c *gin.Context) {
		a11, _ := io.ReadAll(c.Request.Body)
		fmt.Println("收到回调了！", string(a11))
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Any("/xiadingdan", func(context *gin.Context) {
		refund := context.Query("amount")
		refund1, _ := strconv.Atoi(refund)
		codeurl := wechat_pay.Prepay(refund1)
		context.HTML(http.StatusOK, "html.html", gin.H{"code_url": codeurl})
	})

	r.Any("/tuikuan", func(context *gin.Context) {
		refund := context.Query("amount")
		refund1, _ := strconv.Atoi(refund)
		tradno := context.Query("tradeno")
		res := wechat_pay.Refund(tradno, refund1)
		context.JSON(http.StatusOK, res)
	})

	r.Any("/close", func(context *gin.Context) {
		tradno := context.Query("tradeno")
		res := wechat_pay.Close(tradno)
		context.JSON(http.StatusOK, res)
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "html.html", gin.H{"title": "我是测试", "content": "123456"})
	})
	r.Run("0.0.0.0:8181") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
