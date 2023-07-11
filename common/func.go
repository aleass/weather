package common

import (
	"fmt"
	"strconv"
	"strings"
)

func Str2Int64(val string) int64 {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return int64(num)
}

func Str2Int(val string) int {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return num
}

func Int642Str(val int64) string {
	if val == 0 {
		return ""
	}
	return strconv.Itoa(int(val))
}

func DefaultVal(val string) string {
	if val == "" {
		return "0"
	}
	return val
}

func AdjustData(list []DaysPastTimeRank) string {
	str := "[购买情况] code    name"

	//计算长度
	var nameLen, len1Month, len3Month, len6Month, len1year int
	for _, info := range list {
		if len(info.Name) > nameLen {
			nameLen = len(info.Name)
		}
		if len(info.Past1Month) > len1Month {
			len1Month = len(info.Past1Month)
		}
		if len(info.Past3Months) > len3Month {
			len3Month = len(info.Past3Months)
		}
		if len(info.Past6Months) > len6Month {
			len6Month = len(info.Past6Months)
		}
		if len(info.Past1Year) > len1year {
			len1year = len(info.Past1Year)
		}
	}
	str += strings.Repeat(" ", 2*(nameLen-4)/3) + "近一月"
	str += strings.Repeat(" ", 2*(len1Month-3)/3) + "近三月"
	str += strings.Repeat(" ", 2*(len3Month-3)/3) + "近六月"
	str += strings.Repeat(" ", 2*(len6Month-3)/3) + "近一年"
	str += strings.Repeat(" ", 2*(len1year-3)/3) + "成立至今\n"

	for _, info := range list {
		l := (nameLen - len(info.Name)) / 3 * 2
		if l > 0 {
			info.Name += strings.Repeat(" ", l+1)
		}

		info.Past1Month += strings.Repeat(" ", len1Month-len(info.Past1Month))
		info.Past3Months += strings.Repeat(" ", len3Month-len(info.Past3Months))
		info.Past6Months += strings.Repeat(" ", len6Month-len(info.Past6Months))
		info.Past1Year += strings.Repeat(" ", len1year-len(info.Past1Year))

		str += fmt.Sprintf("[%s] %s %s %s %s %s %s %s \n", info.Buy, info.Code, info.Name, info.Past1Month,
			info.Past3Months, info.Past6Months, info.Past1Year, info.SinceInception)
	}
	return str
}
