package helperFunc

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func ModifyDate(needTime string) (string, error) {
	var needDateFormat string = ""
	needTimeSplit := strings.Split(needTime, "T")
	needTimeDateSplit := strings.Split(needTimeSplit[0], "-")
	needTimeTimeSplit := strings.Split(needTimeSplit[1], ".")
	needDateFormat = fmt.Sprintf("%s.%s.%s %s", needTimeDateSplit[2], needTimeDateSplit[1], needTimeDateSplit[0], needTimeTimeSplit[0])
	return needDateFormat, nil
}

func RandomRangeInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}
