package wechat_pay

import "time"

type CreateOrderReq struct {
	Description   string     `json:"description"`
	OutTradeNo    string     `json:"out_trade_no"`
	TimeExpire    time.Time  `json:"time_expire"`
	Attach        string     `json:"attach"`
	NotifyURL     string     `json:"notify_url"`
	GoodsTag      string     `json:"goods_tag"`
	SupportFapiao bool       `json:"support_fapiao"`
	Amount        Amount     `json:"amount"`
	Detail        Detail     `json:"detail"`
	SceneInfo     SceneInfo  `json:"scene_info"`
	SettleInfo    SettleInfo `json:"settle_info"`
}
type Amount struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
}
type GoodsDetail struct {
	MerchantGoodsID  string `json:"merchant_goods_id"`
	WechatpayGoodsID string `json:"wechatpay_goods_id"`
	GoodsName        string `json:"goods_name"`
	Quantity         int64  `json:"quantity"`
	UnitPrice        int64  `json:"unit_price"`
}
type Detail struct {
	CostPrice   int64         `json:"cost_price"`
	InvoiceID   string        `json:"invoice_id"`
	GoodsDetail []GoodsDetail `json:"goods_detail"`
}
type StoreInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AreaCode string `json:"area_code"`
	Address  string `json:"address"`
}
type SceneInfo struct {
	PayerClientIP string    `json:"payer_client_ip"`
	DeviceID      string    `json:"device_id"`
	StoreInfo     StoreInfo `json:"store_info"`
}
type SettleInfo struct {
	ProfitSharing bool `json:"profit_sharing"`
}
