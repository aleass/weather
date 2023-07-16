package service

import (
	"encoding/json"
	"io"
	"net/http"
)

var SkyconStatus = map[string]string{
	"CLEAR_DAY":           "白天晴",  //cloudrate < 0.2
	"CLEAR_NIGHT":         "夜间晴",  //cloudrate < 0.2
	"PARTLY_CLOUDY_DAY":   "白天多云", //0.8 >= cloudrate > 0.2
	"PARTLY_CLOUDY_NIGHT": "夜间多云", //0.8 >= cloudrate > 0.2
	"CLOUDY":              "阴",    //cloudrate > 0.8
	"LIGHT_HAZE":          "轻度雾霾", //PM2.5 100~150
	"MODERATE_HAZE":       "中度雾霾", //PM2.5 150~200
	"HEAVY_HAZE":          "重度雾霾", //PM2.5 > 200
	"LIGHT_RAIN":          "小雨",   //见 降水强度
	"MODERATE_RAIN":       "中雨",   //见 降水强度
	"HEAVY_RAIN":          "大雨",   //见 降水强度
	"STORM_RAIN":          "暴雨",   //见 降水强度
	"FOG":                 "雾",    //能见度低，湿度高，风速低，温度低
	"LIGHT_SNOW":          "小雪",   //降水强度
	"MODERATE_SNOW":       "中雪",   //降水强度
	"HEAVY_SNOW":          "大雪",   //降水强度
	"STORM_SNOW":          "暴雪",   //降水强度
	"DUST":                "浮尘",   //AQI > 150, PM10 > 150，湿度 < 30%，风速 < 6 m/s
	"SAND":                "沙尘",   //AQI > 150, PM10> 150，湿度 < 30%，风速 > 6 m/s
	"WIND":                "大风",
}

// 方向
var WindDirection = [...]string{
	"", //北东北 22.5
	"东北",
	"", //东东北 67.5
	"东",
	"", //东东南 112.5
	"东南",
	"", //南东南 157.5
	"南",
	"", //南西南 202.5
	"西南",
	"", //西西南 247.5
	"西",
	"", //西西北 292.5
	"西北",
	"", //北西北 337.5
	"北",
}

// 转移不常见风向
var UnusualWind = [15]float64{
	0: 22.5, 2: 67.5, 4: 112.5, 6: 157.5, 8: 202.5, 10: 247.5, 12: 292.5, 14: 337.5,
}

var WindLevel = [221]*string{}

// 风力描述
var WindLevelStr = [...]string{
	"无风",
	"1级(微风徐徐)",
	"2级(清风)",
	"3级(树叶摇摆)",
	"4级(树枝摇动)",
	"5级(风力强劲)",
	"6级(风力强劲)",
	"7级(风力超强)",
	"8级(狂风大作)",
	"9级(狂风呼啸)",
	"10(级暴风毁树)",
	"11级(暴风毁树)",
	"12级(飓风)",
	"13级(台风)",
	"14级(强台风)",
	"15级(强台风)",
	"16级(超强台风)",
	"17级(超强台风)",
}

func init() {
	WindLevel[0] = &WindLevelStr[0]
	for i := 1; i <= 5; i++ {
		WindLevel[i] = &WindLevelStr[1]
	}
	for i := 6; i <= 11; i++ {
		WindLevel[i] = &WindLevelStr[2]
	}
	for i := 12; i <= 19; i++ {
		WindLevel[i] = &WindLevelStr[3]
	}
	for i := 20; i <= 28; i++ {
		WindLevel[i] = &WindLevelStr[4]
	}
	for i := 29; i <= 38; i++ {
		WindLevel[i] = &WindLevelStr[5]
	}
	for i := 39; i <= 49; i++ {
		WindLevel[i] = &WindLevelStr[6]
	}
	for i := 50; i <= 61; i++ {
		WindLevel[i] = &WindLevelStr[7]
	}
	for i := 62; i <= 74; i++ {
		WindLevel[i] = &WindLevelStr[8]
	}
	for i := 75; i <= 88; i++ {
		WindLevel[i] = &WindLevelStr[9]
	}
	for i := 89; i <= 102; i++ {
		WindLevel[i] = &WindLevelStr[10]
	}
	for i := 103; i <= 117; i++ {
		WindLevel[i] = &WindLevelStr[11]
	}
	for i := 118; i <= 133; i++ {
		WindLevel[i] = &WindLevelStr[12]
	}
	for i := 134; i <= 149; i++ {
		WindLevel[i] = &WindLevelStr[13]
	}
	for i := 150; i <= 166; i++ {
		WindLevel[i] = &WindLevelStr[14]
	}
	for i := 167; i <= 183; i++ {
		WindLevel[i] = &WindLevelStr[15]
	}
	for i := 184; i <= 201; i++ {
		WindLevel[i] = &WindLevelStr[16]
	}
	for i := 202; i <= 220; i++ {
		WindLevel[i] = &WindLevelStr[17]
	}
}

//<0.031	<0.08	无雨／雪
//0.031 ～ 0.25	0.08~3.44	小雨／雪
//0.25 ～ 0.35	3.44~11.33	中雨／雪
//0.35 ～ 0.48	11.33~51.30	大雨／雪
//>=0.48	>=51.30	暴雨／雪
/*
 "content": [
	{
		"province": "广东省",
		"status": "预警中",
		"code": "1602",
		"description": "广州市气象台14日06时18分发布暴雨黄色和雷雨大风黄色预警信号:受佛山方向移近的雷雨云团影响，预计未来1~3小时广州市越秀区、天河区降水明显，累积雨量30~50毫米，并伴有6级左右短时大风和雷电。从14日06时18分起，广州市越秀区、天河区暴雨和雷雨大风黄色预警信号生效，正值上班上学高峰期，请注意做好防御工作。广州市气象台06月14日06时18分发布。",
		"regionId": "",
		"county": "广州市",
		"pubtimestamp": 1686694680,
		"latlon": [
			23.130061,
			113.264499
		],
		"city": "广东省",
		"alertId": "44010041600000_20230614062133",
		"title": "广州市气象台发布雷雨大风黄色预警[III级/较重]",
		"adcode": "440100",
		"source": "国家预警信息发布中心",
		"location": "广东省广州市",
		"request_status": "ok"
	}
],



标题:广州市气象台发布雷雨大风黄色预警[III级/较重]:
内容:广州市气象台14日06时18分发布暴雨黄色和雷雨大风黄色预警信号:受佛山方向移近的雷雨云团影响，预计未来1~3小时广州市越秀区、天河区降水明显，累积雨量30~50毫米，并伴有6级左右短时大风和雷电。从14日06时18分起，广州市越秀区、天河区暴雨和雷雨大风黄色预警信号生效，正值上班上学高峰期，请注意做好防御工作。广州市气象台06月14日06时18分发布。
状态:预警中
来源:国家预警信息发布中心


*/
type Weather struct {
	Status     string    `json:"status" desc:""`
	ApiVersion string    `json:"api_version" desc:""`
	ApiStatus  string    `json:"api_status" desc:""`
	Lang       string    `json:"lang" desc:""`
	Unit       string    `json:"unit" desc:""`
	Tzshift    int       `json:"tzshift" desc:""`
	Timezone   string    `json:"timezone" desc:""`
	ServerTime int       `json:"server_time" desc:""`
	Location   []float64 `json:"location" desc:""`
	Result     struct {
		Alert struct {
			Status  string `json:"status" desc:""`
			Content []struct {
				//Province     string `json:"province" desc:"省份"`
				Status string `json:"status" desc:"状态"`
				//Code         string `json:"code" desc:""`
				Description string `json:"description" desc:"详情"`
				//RegionId     string `json:"regionId" desc:""`
				//County       string `json:"county" desc:"县"`
				Pubtimestamp int64 `json:"pubtimestamp" desc:"发布时间"`
				//Latlon       []float64 `json:"latlon" desc:"经纬度"`
				//City         string    `json:"city" desc:"城市"`
				//AlertId       string    `json:"alertId" desc:"预警id"`
				Title string `json:"title" desc:"标题"`
				//Adcode        string    `json:"adcode" desc:"代码"`
				Source        string `json:"source" desc:"来源"`
				Location      string `json:"location" desc:"地区"`
				RequestStatus string `json:"request_status" desc:"ok"`
			} `json:"content" desc:""`
			Adcodes []struct {
				//Adcode int    `json:"adcode" desc:""`
				Name string `json:"name" desc:"地点"` //上海 desc:""市
			} `json:"adcodes" desc:""`
		} `json:"alert" desc:"地址"`
		//
		Realtime Realtime `json:"realtime" desc:"实时级别预报"`
		Minutely struct {
			//Status          string    `json:"status" desc:""`
			//Datasource string `json:"datasource" desc:"数据源"`
			//Precipitation2H []float64 `json:"precipitation_2h" desc:"表示未来2小时每一分钟的雷达降水强度"`
			//Precipitation []float64 `json:"precipitation" desc:"降水强度:表示未来1小时每一分钟的雷达降水强度"`
			//Probability     []float64 `json:"probability" desc:"降水概率:未来两小时每半小时的降水概率"`
			Description string `json:"description" desc:"未来2小时天气描述"`
		} `json:"minutely" desc:"分钟级别预报"`
		Hourly struct {
			////Status        string `json:"status" desc:""`
			Description string `json:"description" desc:"未来24小时天气描述"`
			//Precipitation []struct {
			//	//Datetime    string  `json:"datetime" desc:"时间"`
			//	Value       float64 `json:"value" desc:""`
			//	Probability int     `json:"probability" desc:"降水概率(%)"`
			//} `json:"precipitation" desc:"降水强度"`
			//Temperature []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"temperature" desc:"温度"`
			//ApparentTemperature []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"apparent_temperature" desc:""`
			//Wind []struct {
			//	//Datetime  string  `json:"datetime" desc:""`
			//	Speed     float64 `json:"speed" desc:""`
			//	Direction float64 `json:"direction" desc:""`
			//} `json:"wind" desc:""`
			//Humidity []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"humidity" desc:""`
			//Cloudrate []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"cloudrate" desc:""`
			//Skycon []struct {
			//	//Datetime string `json:"datetime" desc:""`
			//	Value string `json:"value" desc:""`
			//} `json:"skycon" desc:""`
			//Pressure []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"pressure" desc:""`
			//Visibility []struct {
			//	//Datetime string  `json:"datetime" desc:"时间"`
			//	Value float64 `json:"value" desc:""`
			//} `json:"visibility" desc:""`
			//Dswrf []struct {
			//	Datetime string  `json:"datetime" desc:"时间"`
			//	Value    float64 `json:"value" desc:""`
			//} `json:"dswrf" desc:""`
			//AirQuality struct {
			//	Aqi []struct {
			//		Datetime string `json:"datetime" desc:""`
			//		Value    struct {
			//			Chn int `json:"chn" desc:""`
			//			Usa int `json:"usa" desc:""`
			//		} `json:"value" desc:""`
			//	} `json:"aqi" desc:""`
			//	Pm25 []struct {
			//		Datetime string `json:"datetime" desc:""`
			//		Value    int    `json:"value" desc:""`
			//	} `json:"pm25" desc:""`
			//} `json:"air_quality" desc:""`
		} `json:"hourly" desc:"小时级别预报"`
		Daily struct {
			//Status string `json:"status" desc:""`
			Astro []struct {
				//	Date    string `json:"date" desc:""`
				Sunrise struct {
					Time string `json:"time" desc:""`
				} `json:"sunrise" desc:"日出"`
				Sunset struct {
					Time string `json:"time" desc:""`
				} `json:"sunset" desc:"日落"`
			} `json:"astro" desc:"日出日落时间"`
			/*Precipitation08H20H []struct {
				//Date        string  `json:"date" desc:""`
				Max         float64 `json:"max" desc:""`
				Min         float64 `json:"min" desc:""`
				Avg         float64 `json:"avg" desc:""`
				Probability int     `json:"probability" desc:""`
			} `json:"precipitation_08h_20h" desc:"白天降水数据"`
			Precipitation20H32H []struct {
				//Date        string  `json:"date" desc:""`
				Max         float64 `json:"max" desc:""`
				Min         float64 `json:"min" desc:""`
				Avg         float64 `json:"avg" desc:""`
				Probability int     `json:"probability" desc:""`
			} `json:"precipitation_20h_32h" desc:"夜晚降水数据"`
			Precipitation []struct {
				//Date        string  `json:"date" desc:""`
				Max         float64 `json:"max" desc:""`
				Min         float64 `json:"min" desc:""`
				Avg         float64 `json:"avg" desc:""`
				Probability int     `json:"probability" desc:""`
			} `json:"precipitation" desc:"降水"`
			*/
			Temperature []struct {
				//Date string  `json:"date" desc:""`
				Max float64 `json:"max" desc:""`
				Min float64 `json:"min" desc:""`
				Avg float64 `json:"avg" desc:""`
			} `json:"temperature" desc:"全天地表 2 米气温"`
			/*
				Temperature08H20H []struct {
					//Date string  `json:"date" desc:""`
					Max float64 `json:"max" desc:""`
					Min float64 `json:"min" desc:""`
					Avg float64 `json:"avg" desc:""`
				} `json:"temperature_08h_20h" desc:"白天地表 2 米气温"`
				Temperature20H32H []struct {
					//Date string  `json:"date" desc:""`
					Max float64 `json:"max" desc:""`
					Min float64 `json:"min" desc:""`
					Avg float64 `json:"avg" desc:""`
				} `json:"temperature_20h_32h" desc:"夜晚地表 2 米气温"`
				//Wind []struct {
				//	Date string `json:"date" desc:""`
				//	Max  struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"max" desc:""`
				//	Min struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"min" desc:""`
				//	Avg struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"avg" desc:""`
				//} `json:"wind" desc:"风速"`
				//Wind08H20H []struct {
				//	Date string `json:"date" desc:""`
				//	Max  struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"max" desc:""`
				//	Min struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"min" desc:""`
				//	Avg struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"avg" desc:""`
				//} `json:"wind_08h_20h" desc:"夜晚地表 10 米风速"`
				//Wind20H32H []struct {
				//	Date string `json:"date" desc:""`
				//	Max  struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"max" desc:""`
				//	Min struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"min" desc:""`
				//	Avg struct {
				//		Speed     float64 `json:"speed" desc:""`
				//		Direction float64 `json:"direction" desc:""`
				//	} `json:"avg" desc:""`
				//} `json:"wind_20h_32h" desc:"夜晚地表 10 米风速"`
				Humidity []struct {
					//Date string  `json:"date" desc:""`
					//Max  float64 `json:"max" desc:""`
					//Min  float64 `json:"min" desc:""`
					Avg float64 `json:"avg" desc:""`
				} `json:"humidity" desc:"地表 2 米相对湿度(%)"`
				//Cloudrate []struct {
				//	Date string  `json:"date" desc:""`
				//	Max  float64 `json:"max" desc:""`
				//	Min  float64 `json:"min" desc:""`
				//	Avg  float64 `json:"avg" desc:""`
				//} `json:"cloudrate" desc:"云量(0.0-1.0)"`
				//Pressure []struct {
				//	Date string  `json:"date" desc:""`
				//	Max  float64 `json:"max" desc:""`
				//	Min  float64 `json:"min" desc:""`
				//	Avg  float64 `json:"avg" desc:""`
				//} `json:"pressure" desc:"地面气压"`
				//Visibility []struct {
				//	Date string  `json:"date" desc:""`
				//	Max  float64 `json:"max" desc:""`
				//	Min  float64 `json:"min" desc:""`
				//	Avg  float64 `json:"avg" desc:""`
				//} `json:"visibility" desc:"地表水平能见度"`
				//Dswrf []struct {
				//	Date string  `json:"date" desc:""`
				//	Max  float64 `json:"max" desc:""`
				//	Min  float64 `json:"min" desc:""`
				//	Avg  float64 `json:"avg" desc:""`
				//} `json:"dswrf" desc:"向下短波辐射通量(W/M2)"`
				//AirQuality struct {
				//	Aqi []struct {
				//		Date string `json:"date" desc:""`
				//		Max  struct {
				//			Chn int `json:"chn" desc:""`
				//			Usa int `json:"usa" desc:""`
				//		} `json:"max" desc:""`
				//		Avg struct {
				//			Chn int `json:"chn" desc:""`
				//			Usa int `json:"usa" desc:""`
				//		} `json:"avg" desc:""`
				//		Min struct {
				//			Chn int `json:"chn" desc:""`
				//			Usa int `json:"usa" desc:""`
				//		} `json:"min" desc:""`
				//	} `json:"aqi" desc:""`
				//	Pm25 []struct {
				//		Date string `json:"date" desc:""`
				//		Max  int    `json:"max" desc:""`
				//		Avg  int    `json:"avg" desc:""`
				//		Min  int    `json:"min" desc:""`
				//	} `json:"pm25" desc:""`
				//} `json:"air_quality" desc:"国标 AQI"`
				Skycon []struct {
					//Date  string `json:"date" desc:""`
					Value string `json:"value" desc:""`
				} `json:"skycon" desc:"全天主要 天气现象"`
				Skycon08H20H []struct {
					//Date  string `json:"date" desc:""`
					Value string `json:"value" desc:""`
				} `json:"skycon_08h_20h" desc:"白天主要 天气现象"`
				Skycon20H32H []struct {
					//Date  string `json:"date" desc:""`
					Value string `json:"value" desc:""`
				} `json:"skycon_20h_32h" desc:"夜晚主要 天气现象"`
				LifeIndex struct {
					Ultraviolet []struct {
						//Date  string `json:"date" desc:""`
						//Index string `json:"index" desc:""`
						Desc string `json:"desc" desc:""`
					} `json:"ultraviolet" desc:"紫外线指数自然语言"`
					//CarWashing []struct {
					//Date  string `json:"date" desc:""`
					//Index string `json:"index" desc:""`
					//Desc string `json:"desc" desc:""`
					//} `json:"carWashing" desc:"洗车指数自然语言"`
					Dressing []struct {
						//Date  string `json:"date" desc:""`
						//Index string `json:"index" desc:""`
						Desc string `json:"desc" desc:""`
					} `json:"dressing" desc:"穿衣指数自然语言"`
					Comfort []struct {
						//Date  string `json:"date" desc:""`
						//Index string `json:"index" desc:""`
						Desc string `json:"desc" desc:""`
					} `json:"comfort" desc:"舒适度指数自然语言"`
					ColdRisk []struct {
						//Date  string `json:"date" desc:""`
						//Index string `json:"index" desc:"等级描述"`
						Desc string `json:"desc" desc:""`
					} `json:"coldRisk" desc:"感冒指数自然语言"`
				} `json:"life_index" desc:"生活指数"`
			*/
		} `json:"daily" desc:"天级别预报"`
		ForecastKeypoint string `json:"forecast_keypoint" desc:"自然语言描述当前情况"`
	} `json:"result" desc:""`
}
type Realtime struct {
	//Status  string        `json:"status" desc:""`
	Temperature float64 `json:"temperature"  desc:"温度"`
	Humidity    float64 `json:"humidity" desc:"地表 2 米湿度相对湿度(%)"`
	//Cloudrate   float64 `json:"cloudrate" desc:"总云量(0.0-1.0)"`
	Skycon string `json:"skycon" desc:"天气现象"`
	//Visibility  float64 `json:"visibility" desc:"地表水平能见度"`
	//Dswrf       float64 `json:"dswrf" desc:"向下短波辐射通量(W/M2)"`
	Wind struct {
		Speed     float64 `json:"speed" desc:"风速 米/秒"`
		Direction float64 `json:"direction" desc:"风向"`
	} `json:"wind" desc:"风速"`

	//Pressure            float64 `json:"pressure" desc:"地面气压"`
	ApparentTemperature float64 `json:"apparent_temperature" desc:"体感温度"`
	Precipitation       struct {
		Local struct {
			//Status     string  `json:"status" desc:""`
			//Datasource string  `json:"datasource" desc:"数据源"`
			Intensity float64 `json:"intensity" desc:"降水强度"`
		} `json:"local" desc:""`
		Nearest struct {
			//Status    string  `json:"status" desc:""`
			Distance  float64 `json:"distance" desc:"最近降水带与本地的距离"`
			Intensity float64 `json:"intensity" desc:"最近降水处的降水强度"`
		} `json:"nearest" desc:""`
	} `json:"precipitation" desc:"降水"`
	AirQuality struct {
		//Pm25 int     `json:"pm25" desc:"PM25 浓度(μg/m3)"`
		//Pm10 int     `json:"pm10" desc:"PM10 浓度(μg/m3)"`
		//O3   int     `json:"o3" desc:"臭氧浓度(μg/m3)"`
		//So2  int     `json:"so2" desc:"二氧化硫浓度(μg/m3)"`
		//No2  int     `json:"no2" desc:"二氧化氮浓度(μg/m3)"`
		//Co   float64 `json:"co" desc:"一氧化碳浓度(mg/m3)"`
		Aqi struct {
			Chn int `json:"chn" desc:"国标"`
			//Usa int `json:"usa" desc:"美标"`
		} `json:"aqi" desc:"空氣品質指標"`
		Description struct {
			Chn string `json:"chn" desc:"国标"`
			//Usa string `json:"usa" desc:"美标"`
		} `json:"description" desc:"空氣品質指標描述"`
	} `json:"air_quality" desc:"空氣品質"`
	LifeIndex struct {
		Ultraviolet struct {
			//Index float64 `json:"index" desc:"级别"`
			Desc string `json:"desc" desc:"描述"`
		} `json:"ultraviolet" desc:"紫外线"`
		Comfort struct {
			//Index int    `json:"index" desc:"级别"`
			Desc string `json:"desc" desc:"描述"`
		} `json:"comfort" desc:"体感舒适度"`
	} `json:"life_index" desc:"生活指数"`
}

// 获取数据
func GetWeatherRawData(url string) (*Weather, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var wegther Weather
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &wegther)
	response.Body.Close()
	return &wegther, err
}
