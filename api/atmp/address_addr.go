package atmp

import (
	"fmt"
	"services/common"
)

const (
	searchUrl = "https://restapi.amap.com/v3/geocode/geo?address=%s&key=%s"
)

// 地址查询
// 地址
// 经纬度
func SearchAddr(addr string) (float64, float64) {
	key := common.MyConfig.Atmp.Key
	url := fmt.Sprintf(searchUrl, addr, key)
	var resp SearchResp
	_, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return 0, 0
	}
	if len(resp.Geocodes) == 0 {
		return 0, 0
	}
	return common.LocStr2float(resp.Geocodes[0].Location)
}

type SearchResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string        `json:"formatted_address"`
		Country          string        `json:"country"`
		Province         string        `json:"province"`
		Citycode         string        `json:"citycode"`
		City             string        `json:"city"`
		District         string        `json:"district"`
		Township         []interface{} `json:"township"`
		Neighborhood     struct {
			Name []interface{} `json:"name"`
			Type []interface{} `json:"type"`
		} `json:"neighborhood"`
		Building struct {
			Name []interface{} `json:"name"`
			Type []interface{} `json:"type"`
		} `json:"building"`
		Adcode   string `json:"adcode"`
		Street   string `json:"street"`
		Number   string `json:"number"`
		Location string `json:"location"`
		Level    string `json:"level"`
	} `json:"geocodes"`
}
