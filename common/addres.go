package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"net/url"
)

var (
	qqMapUrl  = "https://apis.map.qq.com/jsapi?qt=geoc&addr=%s&key="
	geoMapUrl = "https://restapi.amap.com/v3/geocode/geo?address=%s&key="
	qqMapLen  = len(qqMapUrl)
	geoMapLen = len(geoMapUrl)
)

type GeoMap struct {
	Status   string `json:"status"` //返回值为 0 或 1，0 表示请求失败；1 表示请求成功。
	Info     string `json:"info"`   //当 status 为 0 时，info 会返回具体错误原因，否则返回“OK”。
	Count    string `json:"count"`  //返回结果的个数。
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		//Country          string        `json:"country"`
		//Province         string        `json:"province"`
		//Citycode         string        `json:"citycode"`
		//City             string        `json:"city"`
		//District         []interface{} `json:"district"`
		//Township         []interface{} `json:"township"`
		//Neighborhood     struct {
		//	Name []interface{} `json:"name"`
		//	Type []interface{} `json:"type"`
		//} `json:"neighborhood"`
		//Building struct {
		//	Name []interface{} `json:"name"`
		//	Type []interface{} `json:"type"`
		//} `json:"building"`
		//Adcode   string        `json:"adcode"`
		//Street   []interface{} `json:"street"`
		//Number   []interface{} `json:"number"`
		Location string `json:"location"`
		//Level    string        `json:"level"`
	} `json:"geocodes"`
}

type qqMap struct {
	Detail struct {
		//Name            string `json:"name"`
		//City            string `json:"city"`
		//District        string `json:"district"`
		//Adcode          string `json:"adcode"`
		Pointx string `json:"pointx"`
		Pointy string `json:"pointy"`
		//GpsType         string `json:"gps_type"`
		//Reliability     string `json:"reliability"`
		//Province        string `json:"province"`
		//Deviation       string `json:"deviation"`
		//PcdConflictFlag string `json:"pcd_conflict_flag"`
		//QueryStatus     string `json:"query_status"`
		//ServerRetcode   string `json:"server_retcode"`
		//Similarity      string `json:"similarity"`
		//SplitAddr       string `json:"split_addr"`
		//Street          string `json:"street"`
		//StreetNumber    string `json:"street_number"`
		//KeyPoi          string `json:"key_poi"`
		//CategoryCode    string `json:"category_code"`
		//AddressType     string `json:"address_type"`
		//PoiId           string `json:"poi_id"`
		//Town            string `json:"town"`
		//TownCode        string `json:"town_code"`
		//TownLevel       int    `json:"town_level"`
		//KeyRole         string `json:"key_role"`
		//ShortAddress    string `json:"short_address"`
		AnalysisAddress string `json:"analysis_address"`
		//FormatAddress   string `json:"format_address"`
		//PersonName      string `json:"person_name"`
		//Tel             string `json:"tel"`
	} `json:"detail"`
}

func SetToken(qqToken, geoToken string) {
	if qqToken != "" {
		qqMapUrl += qqToken
	}
	if geoToken != "" {
		geoMapUrl += geoToken
	}

}

func GetKeyWordAddr(word string) (string, string, error) {
	if len(qqMapUrl) == qqMapLen || len(geoMapUrl) == geoMapLen {
		LogSend("token 不存在", ErrType)
		return "", "", nil
	}

	work := url.QueryEscape(word)
	if len(qqMapUrl) != qqMapLen {
		var res, err = http.Get(fmt.Sprintf(qqMapUrl, work))
		if err != nil {
			LogSend(err.Error(), ErrType)
			return "", "", nil
		}
		raw, _ := io.ReadAll(res.Body)
		decoder := simplifiedchinese.GBK.NewDecoder()
		utf8Bytes, err := decoder.Bytes(raw)
		if err != nil {
			LogSend(err.Error(), ErrType)
			return "", "", nil
		}
		data := &qqMap{}
		json.Unmarshal(utf8Bytes, data)
		if data.Detail.AnalysisAddress != "" {
			return data.Detail.Pointx + "," + data.Detail.Pointy, data.Detail.AnalysisAddress, nil
		}
	}
	if len(geoMapUrl) != geoMapLen {
		var res, err = http.Get(fmt.Sprintf(geoMapUrl, work))
		if err != nil {
			LogSend(err.Error(), ErrType)
			return "", "", nil
		}
		raw, _ := io.ReadAll(res.Body)
		data := &GeoMap{}
		json.Unmarshal(raw, data)
		if data.Status == "0" {
			LogSend(data.Info, ErrType)
			return "", "", nil
		}
		if len(data.Count) > 1 {
			var addr string
			for _, v := range data.Geocodes {
				addr += v.FormattedAddress + "\r\n"
			}
			return "", "", errors.New("搜索结果过多:")
		}
		switch len(data.Count) {
		case 1:
			return data.Geocodes[0].Location, data.Geocodes[0].FormattedAddress, nil
		case 0:
			return "", "", errors.New("搜索不到")
		}
	}
	return "", "", nil

}
