package wechat_pay

import (
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/consts"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/certificates"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	nethttp "net/http"
	neturl "net/url"
	"time"
	"wechat-pay/conf"
)

type PrepayRes struct {
	CodeUrl string `json:"code_url"`
	TradeNo string `json:"trade_no"`
}

func GetTradeNo() string {
	return fmt.Sprint(time.Now().UnixMilli()) + time.Now().Local().Format("060102150405")
}

func PrepayTest(amount int) *PrepayRes {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	tradeno := GetTradeNo()
	svc := native.NativeApiService{Client: client}
	resp, result, err := svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:         core.String(conf.Conf.Appid),
			Mchid:         core.String(conf.Conf.MchID),
			Description:   core.String("小臭臭"),
			OutTradeNo:    core.String(tradeno),
			TimeExpire:    core.Time(time.Now()),
			Attach:        core.String("小肉肉"),
			NotifyUrl:     core.String(conf.Conf.PayNotify),
			GoodsTag:      core.String("WXG"),
			SupportFapiao: core.Bool(false),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(int64(amount)),
			},
			Detail: &native.Detail{
				CostPrice: core.Int64(608800),
				GoodsDetail: []native.GoodsDetail{native.GoodsDetail{
					GoodsName:        core.String("奥特曼"),
					MerchantGoodsId:  core.String("ABC"),
					Quantity:         core.Int64(1),
					UnitPrice:        core.Int64(828800),
					WechatpayGoodsId: core.String("1001"),
				}},
				InvoiceId: core.String("wx123"),
			},
			SettleInfo: &native.SettleInfo{
				ProfitSharing: core.Bool(false),
			},
			SceneInfo: &native.SceneInfo{
				DeviceId:      core.String("013467007045764"),
				PayerClientIp: core.String("14.23.150.211"),
				StoreInfo: &native.StoreInfo{
					Address:  core.String("北京市朝阳区北苑家园"),
					AreaCode: core.String("440305"),
					Id:       core.String("0001"),
					Name:     core.String("小肉肉"),
				},
			},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call Prepay err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
	return &PrepayRes{
		CodeUrl: *resp.CodeUrl,
		TradeNo: tradeno,
	}
}

func Prepay(order *CreateOrderReq) *PrepayRes {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	tradeno := order.OutTradeNo
	if tradeno == "" {
		tradeno = GetTradeNo()
	}
	svc := native.NativeApiService{Client: client}
	nativePrepayReq := native.PrepayRequest{
		Appid:         core.String(conf.Conf.Appid),
		Mchid:         core.String(conf.Conf.MchID),
		Description:   core.String(order.Description),
		OutTradeNo:    core.String(tradeno),
		TimeExpire:    core.Time(time.Now()),
		Attach:        core.String(order.Attach),
		NotifyUrl:     core.String(conf.Conf.PayNotify),
		GoodsTag:      core.String(order.GoodsTag),
		SupportFapiao: core.Bool(false),
		Amount: &native.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(order.Amount.Total)),
		},
		SettleInfo: &native.SettleInfo{
			ProfitSharing: core.Bool(order.SettleInfo.ProfitSharing),
		},
		SceneInfo: &native.SceneInfo{
			DeviceId:      core.String(order.SceneInfo.DeviceID),
			PayerClientIp: core.String(order.SceneInfo.PayerClientIP),
			StoreInfo: &native.StoreInfo{
				Address:  core.String(order.SceneInfo.StoreInfo.Address),
				AreaCode: core.String(order.SceneInfo.StoreInfo.AreaCode),
				Id:       core.String(order.SceneInfo.StoreInfo.ID),
				Name:     core.String(order.SceneInfo.StoreInfo.Name),
			},
		},
	}
	nativePrepayReq.Detail = &native.Detail{
		GoodsDetail: make([]native.GoodsDetail, 0),
	}
	for _, v := range order.Detail.GoodsDetail {
		nativePrepayReq.Detail.GoodsDetail = append(nativePrepayReq.Detail.GoodsDetail, native.GoodsDetail{
			GoodsName:        core.String(v.GoodsName),
			MerchantGoodsId:  core.String(v.MerchantGoodsID),
			Quantity:         core.Int64(v.Quantity),
			UnitPrice:        core.Int64(v.UnitPrice),
			WechatpayGoodsId: core.String(v.WechatpayGoodsID),
		})
	}
	resp, result, err := svc.Prepay(ctx, nativePrepayReq)

	if err != nil {
		// 处理错误
		log.Printf("call Prepay err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
	return &PrepayRes{
		CodeUrl: *resp.CodeUrl,
		TradeNo: tradeno,
	}
}

func Refund(tradNo string, amount int) *RefundResp {
	resp, _, err := RefundPay(&RefundReq{
		OutRefundNo: tradNo,
		OutTradeNo:  tradNo,
		Amount: Amount1{
			Refund:   amount,
			Total:    amount,
			Currency: "CNY",
		},
	})
	if err != nil {
		fmt.Println("退款出错啦！", err)
	}
	return resp
}

func Close(tradeNo string) string {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := native.NativeApiService{Client: client}
	result, err := svc.CloseOrder(ctx,
		native.CloseOrderRequest{
			OutTradeNo: core.String(tradeNo),
			Mchid:      core.String(conf.Conf.MchID),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call CloseOrder err:%s", err)
		return "关闭订单失败了！" + err.Error()
	} else {
		// 处理返回结果
		log.Printf("status=%d", result.Response.StatusCode)
		return "关闭订单成功了！"
	}
}

func Query(tradeNo string) *payments.Transaction {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := native.NativeApiService{Client: client}
	resp, _, err := svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(tradeNo),
		Mchid:      core.String(conf.Conf.MchID),
	})
	if err != nil {
		fmt.Println("查询订单出错啦！", err)
	}
	return resp
}

func DownloadCerts() *certificates.DownloadCertificatesResponse {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}

	// 发送请求，以下载微信支付平台证书为例
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
	svc := certificates.CertificatesApiService{Client: client}
	resp, result, err := svc.DownloadCertificates(ctx)
	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	return resp
}

type RefundReq struct {
	OutRefundNo string  `json:"out_refund_no"`
	OutTradeNo  string  `json:"out_trade_no"`
	Amount      Amount1 `json:"amount"`
}

type Amount1 struct {
	Refund   int    `json:"refund"`
	Total    int    `json:"total"`
	Currency string `json:"currency"`
}

type RefundResp struct {
	RefundId            string       `json:"refund_id"`
	OutRefundNo         string       `json:"out_refund_no"`
	TransactionId       string       `json:"transaction_id"`
	OutTradeNo          string       `json:"out_trade_no"`
	Channel             string       `json:"channel"`
	UserReceivedAccount string       `json:"user_received_account"`
	CreateTime          string       `json:"create_time"`
	Status              string       `json:"status"`
	Amount              RefundAmount `json:"amount"`
}

type RefundAmount struct {
	Refund           int    `json:"refund"`
	Total            int    `json:"total"`
	PayerTotal       int    `json:"payer_total"`
	PayerRefund      int    `json:"payer_refund"`
	SettlementRefund int    `json:"settlement_refund"`
	DiscountRefund   int    `json:"discount_refund"`
	Currency         string `json:"currency"`
}

func RefundPay(req *RefundReq) (resp *RefundResp, result *core.APIResult, err error) {
	var (
		localVarHTTPMethod   = nethttp.MethodPost
		localVarPostBody     interface{}
		localVarQueryParams  neturl.Values
		localVarHeaderParams = nethttp.Header{}
	)

	ctx := context.Background()

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.Conf.ClientKeyPath)
	if err != nil {
		log.Print("load merchant private key error")
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.Conf.MchID, conf.Conf.MchCertificateSerialNumber, mchPrivateKey, conf.Conf.MchAPIv3Key),
	}

	localVarPath := consts.WechatPayAPIServer + "/v3/refund/domestic/refunds"
	// Make sure All Required Params are properly set

	// Setup Body Params
	localVarPostBody = req

	// Determine the Content-Type Header
	localVarHTTPContentTypes := []string{"application/json"}
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	// Perform Http Request
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	result, err = client.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams, localVarQueryParams, localVarPostBody, localVarHTTPContentType)
	if err != nil {
		return nil, result, err
	}

	// Extract PrepayResponse from Http Response
	resp = new(RefundResp)
	err = core.UnMarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}
