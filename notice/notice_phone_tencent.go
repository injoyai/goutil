package notice

type PhoneTencentConfig struct {
	SecretID  string `json:"secretID"`  //
	SecretKey string `json:"secretKey"` //
	SignName  string `json:"signName"`  //签名
	AppID     string `json:"appID"`     //腾讯云用
}
