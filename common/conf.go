package common

type Config struct {
	// 结构映射
	Wechat []struct {
		Token string `mapstructure:"token"`
		Notes string `mapstructure:"note"`
	} `mapstructure:"wechat"`
	CaiYun struct {
		Token  string `json:"token"`
		Addres []struct {
			Name        string `json:"name"`
			WechatNotes string `json:"wechatNotes"`
			Coordinate  string `json:"coordinate"`
			Switch      bool   `json:"switch" desc:"开关"`
		} `json:"addres"`
	} `json:"caiyun"`
}
