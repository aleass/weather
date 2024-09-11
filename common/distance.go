package common

import "math"

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
