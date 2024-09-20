package atmp

import (
	"bytes"
	"fmt"
	"weather/common"
)

const (
	createPhotoUrl = "https://restapi.amap.com/v3/staticmap?paths=%s&location=%s&key=%s&zoom=%d&size=%s"
)

// 地址查询
// 地址
// 经纬度
func CreatePhoto(loc []string, name string, radius7, radius10, radius12 float64) *bytes.Reader {
	//生产参数
	//台风路径
	paths := "2,0x0000FF,0.2,,:"
	for _, l := range loc {
		paths += l + ";"
	}
	paths = paths[:len(paths)-1]
	l1, l2 := common.LocStr2float(loc[len(loc)-1])

	//台风7级风圈
	loc2 := common.CalculateCoordinates(l1, l2, radius7, 360)
	paths += "|2,0xb2b7be,0.5,0xb2b7be,0.5:"
	for _, l := range loc2 {
		paths += l + ";"
	}
	paths = paths[:len(paths)-1]

	//台风10级风圈
	loc2 = common.CalculateCoordinates(l1, l2, radius10, 360)
	paths += "|2,0xd0a774,0.5,0xd0a774,0.5:"
	for _, l := range loc2 {
		paths += l + ";"
	}
	paths = paths[:len(paths)-1]

	//台风12级风圈
	loc2 = common.CalculateCoordinates(l1, l2, radius12, 360)
	paths += "|2,0xe69c41,0.5,0xe69c41,0.5:"
	for _, l := range loc2 {
		paths += l + ";"
	}
	paths = paths[:len(paths)-1]

	var (
		key = common.MyConfig.Atmp.Key
		//labels  = fmt.Sprintf("人,,,,,:%s", common.MyConfig.Atmp.Loc)
		//markers = fmt.Sprintf(",,:%s", loc[len(loc)-1])
		location = loc[len(loc)-1]
		size     = "500*500"
		zoom     = 5
	)

	url := fmt.Sprintf(createPhotoUrl, paths, location, key, zoom, size)

	//加标签
	if radius7+radius10+radius12 == 0 {
		var markers = fmt.Sprintf(",,:%s", loc[len(loc)-1])
		url += "&markers=" + markers
	}

	//labels
	//url += "&labels=" + labels

	resp, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, nil)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	return bytes.NewReader(resp)
}
