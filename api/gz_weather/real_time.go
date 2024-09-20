package gz_weather

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"weather/api/atmp"
	"weather/common"
)

// 仅支持广州
const (
	wdphRain = "https://weixin.tqyb.com.cn/gzweixin//wSituation/weixin_sk_wdph.flow?type=%d&"
	msgTemp  = common.SubStr + "%s  %0.2f\n" //
)

// 记录地址距离
var addrDisance = map[string]float64{}

func GZWeather(district string) string {
	values := url.Values{}
	// 添加查询参数
	values.Add("province", "广东省")
	values.Add("city", "广州市")
	values.Add("district", district)

	// 编码成 URL 查询字符串
	encoded := values.Encode()
	var lonSelf, latSelf = common.LocStr2float(common.MyConfig.Atmp.Loc)

	var (
		msg         string
		weatherData = handler(1, encoded)
		rainInfos   []rainInfo
		msgMaxLen   int
	)

	if weatherData == nil {
		goto qu
	}
	//handler data
	for _, v := range weatherData.ParamBuf {
		//filtration
		if v.Col2 == 0 {
			continue
		}
		//同区计算距离
		var dis float64 = -1
		if _dis, ok := addrDisance[v.Col1]; ok {
			dis = _dis
		} else {
			lon, lat := atmp.SearchAddr(v.Col1)
			if lon > 0 {
				dis = common.Haversine(lon, lat, lonSelf, latSelf)
				addrDisance[v.Col1] = dis
			}
		}
		if msgMaxLen < len(v.Col1) {
			msgMaxLen = len(v.Col1)
		}
		rainInfos = append(rainInfos, rainInfo{dis, v.Col1, v.Col2}) //msgTemp+"   %0.2f km\n"
	}
	if len(rainInfos) > 0 {
		msg += "【当前地区降水检测点】\n"
		sort.Slice(rainInfos, func(i, j int) bool {
			return rainInfos[i].Dis < rainInfos[j].Dis
		})
		msgMaxLen += 2
		for _, str := range rainInfos {
			var formtat = common.SubStr + "%s" + strings.Repeat(" ", int(float64(msgMaxLen-len(str.addr))*1.3)) + "%0.2f   	%0.2f km\n"
			msg += fmt.Sprintf(formtat, str.addr, str.rainfall, str.Dis)
		}
	}

qu:
	//全区
	weatherData = handler(2, encoded)
	if weatherData == nil {
		return msg
	}
	msg += "\n【全市区降水检测点】\n"
	for _, v := range weatherData.ParamBuf {
		if v.Col2 == 0 {
			continue
		}
		msg += fmt.Sprintf(msgTemp, v.Col1, v.Col2)
	}

	return msg + "\n"
}

type rainInfo struct {
	Dis      float64
	addr     string
	rainfall float64
}

func handler(types int, encoded string) *Data {
	//url types 1开始
	urls := fmt.Sprintf(wdphRain, types) + encoded
	xmlData, err := common.HttpRequest(common.OtherType, common.GetType, urls, nil, nil, false, nil)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	var data Data
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	return &data
}

type Data struct {
	//State     string `xml:"state"`
	//ErrorCode int    `xml:"errorcode"`
	//Msg       string `xml:"msg"`
	ParamBuf []Col `xml:"parambuf>col"`
}

type Col struct {
	Col1 string  `xml:"col1"`
	Col2 float64 `xml:"col2"`
}
