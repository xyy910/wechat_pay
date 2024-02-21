package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"wechat-pay/conf"
	"wechat-pay/wechat_pay"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("GET current path err:", err)
	}
	conf.InitConf(dir)

	r := gin.Default()
	r.LoadHTMLGlob("fe/*html")
	r.Static("/static", "./fe/dist")

	r.Any("/prepay/demo", func(context *gin.Context) {
		refund := context.Query("amount")
		refund1, _ := strconv.Atoi(refund)
		res := wechat_pay.PrepayTest(refund1)
		context.HTML(http.StatusOK, "html.html",
			gin.H{"code_url": res.CodeUrl,
				"trade_no": res.TradeNo})
	})

	r.POST("/prepay", func(context *gin.Context) {
		var req wechat_pay.CreateOrderReq
		err := context.ShouldBind(&req)
		if err != nil {
			context.JSON(http.StatusBadRequest, err)
		} else {
			res := wechat_pay.Prepay(&req)
			context.JSON(http.StatusOK, res)
		}
	})

	r.Any("/order/query", func(context *gin.Context) {
		tradeno := context.Query("tradeno")
		res := wechat_pay.Query(tradeno)
		context.JSON(http.StatusOK, res)
	})

	r.Any("/refund", func(context *gin.Context) {
		refund := context.Query("amount")
		refund1, _ := strconv.Atoi(refund)
		tradno := context.Query("tradeno")
		res := wechat_pay.Refund(tradno, refund1)
		context.JSON(http.StatusOK, res)
	})

	r.Any("/order/close", func(context *gin.Context) {
		tradno := context.Query("tradeno")
		res := wechat_pay.Close(tradno)
		context.JSON(http.StatusOK, res)
	})

	r.Any("/certificates", func(context *gin.Context) {
		res := wechat_pay.DownloadCerts()
		context.JSON(http.StatusOK, res)
	})

	r.Run("0.0.0.0:" + conf.Conf.ServerPort)
}
