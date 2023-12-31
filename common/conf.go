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
			Addr        string `json:"addr"`
			WechatNotes string `json:"wechatNotes"`
			Coordinate  string `json:"coordinate"`
			Switch      bool   `json:"switch" desc:"开关"`
			AllowWeek   string `json:"allowWeek"`
		} `json:"addres"`
	} `json:"caiyun"`

	UrlConfigPass []struct {
		Name  string `json:"name"`
		Notes string `mapstructure:"note"`
	} `json:"urlConfigPass"`

	GeoMapToken string `json:"geoMapToken"`
	QqMapToken  string `json:"qqMapToken"`

	DB struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		DbName   string `json:"dbName"`
		Port     string `json:"port"`
	} `json:"db"`

	Fund struct {
		Host  string   `json:"host"`
		Notes []string `mapstructure:"notes"`
	} `json:"fund"`
}
