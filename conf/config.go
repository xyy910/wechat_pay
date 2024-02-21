package conf

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

var Conf Config

type Config struct {
	ServerPort                 string `toml:"serverPort"`
	MchID                      string `toml:"mchID"`
	MchCertificateSerialNumber string `toml:"mchCertificateSerialNumber"`
	MchAPIv3Key                string `toml:"mchAPIv3Key"`
	Appid                      string `toml:"appid"`
	ClientKeyPath              string `toml:"clientKeyPath"`
	PayNotify                  string `toml:"payNotify"`
}

func InitConf(path string) {
	if _, err := toml.DecodeFile(path+"/config.toml", &Conf); err != nil {
		log.Fatalln("InitConf err: ", err)
	} else {
		Conf.ClientKeyPath = path + "/wechat_pay/business_cert/apiclient_key.pem"
		aa1, _ := json.Marshal(Conf)
		fmt.Println("配置是：", string(aa1))
	}
}
