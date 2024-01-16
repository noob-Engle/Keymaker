package pkg

import (
	"strconv"
	"strings"
	"time"
)

type LTypeTransfrom struct {
}

func (LTypeTransfrom) DtoText(anydata any) (returnVal string) {
	return any_to_doc(anydata)
}

func (LTypeTransfrom) DToInt(anydata any) (returnVal int) {
	text := any_to_doc(anydata)
	pos := strings.Index(text, ".")
	if pos != -1 {
		text = text[:pos]
	}
	if text == "true" {
		return 1
	}
	res, err := strconv.ParseInt(text, 0, 64)
	if err != nil {
		return 0
	}
	returnVal = int(res)
	return
}

func (LTypeTransfrom) DToDubble(anydata any) (returnVal float64) {
	text := any_to_doc(anydata)
	res, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	returnVal = res
	return
}

func (LTypeTransfrom) DToBoolean(anydata any) (returnVal bool) {
	text := any_to_doc(anydata)
	if text == "true" || text == "1" {
		return true
	}
	if text == "false" || text == "0" {
		return false
	}
	res, err := strconv.ParseBool(text)
	if err != nil {
		return false
	}
	returnVal = res
	return
}

func (LTypeTransfrom) DToString(anydata any) (returnVal []byte) {
	text := any_to_doc(anydata)
	returnVal = []byte(text)
	return
}

// 格式参照
// @	"2006-01-02 15:04:05",
// @	"2006/01/02 15:04:05",
// @	"2006/01/02",
// @	"2006-01-02",
// @	"2006年01月02日",
// @	"2006年01月02日 15时04分05秒",
// @	"2006年01月02日 15:04:05",
// @	"January 02, 2006 15:04:05",
// @	"02 Jan 2006 15:04:05",
// @	"02.01.2006 15:04:05",
// @	"Mon Jan 02 15:04:05 MST 2006",
// @	"2006-01-02T15:04:05",
// @	"2006-01-02 15:04:05Z",
// @	"2006-01-02T15:04:05Z",
// @	"2006-01-02T15:04:05.999999999Z07:00",
// @	"Jan _2 15:04:05",
// @	"Jan _2 15:04:05.000",
// @	"Jan _2 15:04:05.000000",
// @	"Jan _2 15:04:05.000000000",
// @	 "15:04:05",
// @	 "01/02 03:04:05PM '06 -0700",
// @	 "Mon Jan _2 15:04:05 2006",
// @	 "Mon Jan _2 15:04:05 MST 2006",
// @	 "Mon Jan 02 15:04:05 -0700 2006",
// @	 "02 Jan 06 15:04 MST",
// @	 "02 Jan 06 15:04 -0700",
// @	 "Monday, 02-Jan-06 15:04:05 MST",
// @	 "Mon, 02 Jan 2006 15:04:05 MST",
// @	 "Mon, 02 Jan 2006 15:04:05 -0700",
// @	 "2006-01-02T15:04:05Z07:00",
// @	 "3:04PM",

func (LTypeTransfrom) DToTime(anydata any, fromat_option ...string) (returnVal time.Time) {
	text := any_to_doc(anydata)
	if len(fromat_option) > 0 && fromat_option[0] != "" {
		returnVal, _ = time.Parse(fromat_option[0], text)

	} else {
		for _, v := range all_timefromt {
			if restime, err := time.Parse(v, text); err == nil {
				returnVal = restime
				return
			}

		}

	}

	return
}
