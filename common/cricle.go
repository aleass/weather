package common

import (
	"fmt"
	"math"
)

const earthRadius = 6371.0 // 地球半径，单位公里

// 度转弧度
func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}

// 弧度转度
func toDegrees(radian float64) float64 {
	return radian * 180 / math.Pi
}

// 计算7级风范围的经纬度点
func CalculateCoordinates(lon0, lat0, radiusKm float64, numPoints int) []string {
	coords := make([]string, numPoints)

	// 将纬度和经度转换为弧度
	lat0 = toRadians(lat0)
	lon0 = toRadians(lon0)

	for i := 0; i < numPoints; i++ {
		// 角度从0到360度
		theta := 2 * math.Pi * float64(i) / float64(numPoints)

		// 计算新的纬度
		latNew := math.Asin(math.Sin(lat0)*math.Cos(radiusKm/earthRadius) +
			math.Cos(lat0)*math.Sin(radiusKm/earthRadius)*math.Cos(theta))

		// 计算新的经度
		lonNew := lon0 + math.Atan2(math.Sin(theta)*math.Sin(radiusKm/earthRadius)*math.Cos(lat0),
			math.Cos(radiusKm/earthRadius)-math.Sin(lat0)*math.Sin(latNew))

		// 转换回度数并保存
		coords[i] = fmt.Sprintf("%f,%f", toDegrees(lonNew), toDegrees(latNew))
	}
	return coords
}
