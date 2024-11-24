package atmp

import (
	"bytes"
	"fmt"
	"services/common"
	"strings"
)

const (
	createPhotoUrl = "https://restapi.amap.com/v3/staticmap?paths=%s&location=%s&key=%s&zoom=%d&size=%s"
)

// 地址查询
// 地址
// 经纬度
func CreatePhoto(loc string, locArr, forecasts []string, forecastsName, forecastsDate [][2]string, dis, radius7, radius10, radius12 float64) *bytes.Reader {
	//生产参数
	//台风路径
	paths := "2,0x0000FF,0.2,,:"
	for _, l := range locArr {
		paths += l + ";"
	}
	//预测
	for _, l := range forecasts {
		paths += l + ";"
	}
	paths = paths[:len(paths)-1]
	l1, l2 := common.LocStr2float(locArr[len(locArr)-1])

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

	//预测名称
	var labels string
	//for _, info := range forecastsName {
	//	labels += fmt.Sprintf("%s,0,0,14,0x000000,0xffffff:%s|", info[0], info[1]) //日本,0,0,14,0x000000,0xffffff:118.713386,20.727620
	//}

	for _, info := range forecastsDate {
		labels += fmt.Sprintf("%s,0,0,14,0x000000,0xffffff:%s|", strings.Replace(info[0], " ", "_", 1), info[1]) //日本,0,0,14,0x000000,0xffffff:118.713386,20.727620
	}

	var (
		key      = common.MyConfig.Atmp.Key
		location = common.CalculateMidpoint(locArr[len(locArr)-1], loc)
		size     = "1000*1000"
		zoom     = zoomHandler(dis)
	)

	url := fmt.Sprintf(createPhotoUrl, paths, location, key, zoom, size)

	//标签
	if labels != "" {
		labels = labels[:len(labels)-1]
		url += "&labels=" + labels
	}

	if radius7 == 0 && radius10 == 0 && radius12 == 0 {
		url += fmt.Sprintf("&markers=small,0x000000,:%s", locArr[len(locArr)-1])
	}

	resp, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, nil)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	return bytes.NewReader(resp)
}

// 500 缩放
// 3 1000
// 4 500
// 5 200
// 6 100
// 7 50
// 根据距离 计算缩放
func zoomHandler(dis float64) int {
	//1080 下 缩放
	switch {
	case dis <= 0 || dis > 4000:
		return 3
	case dis > 3000:
		return 4
	case dis > 1000:
		return 5
	default:
		return 6
	}
}
