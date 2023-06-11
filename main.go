package main

func main() {
	//开启携程,分别监控各个地区的天气
	for i, info := range Location {
		go watch_weather(i, info)
		i++
	}
	select {}
}
