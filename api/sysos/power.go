package sysos

import (
	"services/common"
)

func handlerPower(line string) {
	if len(line) < 14 {
		return
	}
	if line[:14] != "Combined Power" {
		return
	}
	//Combined Power (CPU + GPU + ANE): 289 mW
	OSPower = line[34:]
	common.Logger.Info(OSPower)
}
