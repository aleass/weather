package common

import (
	"fmt"
	"math"
	"strings"
)

// Haversine 公式计算两点间的距离
func Haversine(lon1, lat1, lon2, lat2 float64) float64 {
	const R = 6371 // 地球半径，单位为公里

	// 将角度转为弧度
	lon1, lat1, lon2, lat2 = degToRad(lon1), degToRad(lat1), degToRad(lon2), degToRad(lat2)

	// 计算经度和纬度差值
	dlon := lon2 - lon1
	dlat := lat2 - lat1

	// Haversine 公式
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 计算距离
	distance := R * c
	return distance
}

// 角度转弧度
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// 计算中间点的函数
func CalculateMidpoint(loc1, loc2 string) string {
	locs1 := strings.Split(loc1, ",")
	locs2 := strings.Split(loc2, ",")
	lat1, lon1 := Str2Float64(locs1[1]), Str2Float64(locs1[0])
	lat2, lon2 := Str2Float64(locs2[1]), Str2Float64(locs2[0])

	// 将经纬度转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// 计算经度差
	dLon := lon2Rad - lon1Rad

	// 计算中点坐标
	Bx := math.Cos(lat2Rad) * math.Cos(dLon)
	By := math.Cos(lat2Rad) * math.Sin(dLon)

	latMid := math.Atan2(math.Sin(lat1Rad)+math.Sin(lat2Rad), math.Sqrt(math.Pow(math.Cos(lat1Rad)+Bx, 2)+math.Pow(By, 2)))
	lonMid := lon1Rad + math.Atan2(By, math.Cos(lat1Rad)+Bx)

	// 将结果转换为度
	latMid = latMid * 180 / math.Pi
	lonMid = lonMid * 180 / math.Pi

	return fmt.Sprintf("%f,%f", lonMid, latMid)
}
