package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 键值列表_类型筛选
func KVListFilter(data any) (retrunError error) {
	switch nowdata := data.(type) {
	case string:
		return
	case int:
		return
	case bool:
		return
	case byte:
		return
	case nil:
		return
	case float64:
		return
	case float32:
		return
	case int8:
		return
	case int16:
		return
	case int32:
		return
	case int64:
		return
	case uint:
		return
	case uint16:
		return
	case uint32:
		return
	case uint64:
		return
	case []byte:
		return
	case time.Time:
		return
	case []map[string]any:
		for _, v := range nowdata {
			retrunError = KVListFilter(v)
			if retrunError != nil {
				return
			}
		}
		return
	case map[string]any:
		for _, v := range nowdata {
			retrunError = KVListFilter(v)
			if retrunError != nil {
				return
			}
		}
		return
	case []any:
		for _, v := range nowdata {
			retrunError = KVListFilter(v)
			if retrunError != nil {
				return
			}
		}
		return
	default:
		retrunError = fmt.Errorf("错误: 类型 %T 不是支持的基本类型", nowdata)
		return
	}

}

// 键值列表_深拷贝
func KVListDeepcopy(data any) (returnValue any, returnError error) {
	switch nowdata := data.(type) {
	case []any:
		returnValue = Deepcopy(nowdata)
		return
	case map[string]any:
		returnValue = Deepcopy(nowdata)
		return
	case []map[string]any:
		tranformData := make([]any, len(nowdata))
		for i, v := range nowdata {
			tranformData[i] = v
		}
		returnValue = Deepcopy(returnValue)
		return
	default:
		returnError = errors.New("只支持 []any 或者 map[string]any  或者 []map[string]any 当前数据类型错误")
		return
	}

}

// 深拷贝
func Deepcopy(data any) (returnData any) {
	if conditionMap, ok := data.(map[string]any); ok {
		newMap := make(map[string]any)
		for k, v := range conditionMap {
			newMap[k] = Deepcopy(v)
		}
		returnData = newMap
		return
	} else if conditionSlice, ok := data.([]any); ok {
		newSlice := make([]any, len(conditionSlice))
		for k, v := range newSlice {
			newSlice[k] = Deepcopy(v)
		}
		returnData = newSlice
		return
	}
	returnData = data
	return
}

// 键值列表_JSON之前处理
func KVList_BeforeHandleJson(data any) (value any) {
	if conditionMap, ok := data.(map[string]any); ok {
		newMap := make(map[string]any)
		for k, v := range conditionMap {
			switch nowData := v.(type) {
			case time.Time:
				newMap[k] = KVList_BeforeHandleJson(nowData.Format("2006-01-02 15:04:05"))
			case []byte:
				newMap[k] = KVList_BeforeHandleJson(string(nowData))

			default:
				newMap[k] = KVList_BeforeHandleJson(v)

			}
		}
		value = newMap
		return
	} else if conditionslices, ok := data.([]any); ok {
		newslices := make([]any, len(conditionslices))
		for k, v := range conditionslices {
			switch nowValue := v.(type) {
			case time.Time:
				newslices[k] = KVList_BeforeHandleJson(nowValue.Format("2006-01-02 15:04:05"))
			case []byte:
				newslices[k] = KVList_BeforeHandleJson(string(nowValue))

			default:
				newslices[k] = KVList_BeforeHandleJson(v)

			}
		}
		value = newslices
		return
	}
	value = data
	return
}

func KVList_SetValue(data any, keyarrpath []string, value any) (return_error error) {
	return_error = KVListFilter(value)
	if return_error != nil {
		return
	}
	newValue, err := KVListDeepcopy(value)
	if err != nil {
		switch nowValue := data.(type) {
		case []any:
			_, return_error = setValue(nowValue, keyarrpath, value)
			return
		case map[string]any:
			_, return_error = setValue(nowValue, keyarrpath, value)
			return
		default:
			return_error = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
			return
		}
	}
	switch nowdata := data.(type) {
	case []any:
		_, return_error = setValue(nowdata, keyarrpath, newValue)
		return
	case map[string]any:
		_, return_error = setValue(nowdata, keyarrpath, newValue)
		return
	default:
		return_error = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
		return
	}
}

// 置值
func setValue(data any, keyArrPath []string, value any) (returnValue any, returnError error) {
	lenKeys := len(keyArrPath)
	if lenKeys == 0 {
		returnError = errors.New("键组路径不能为空")
		return
	}
	if lenKeys == 1 {
		switch nowdata := data.(type) {
		case map[string]any:
			nowdata[keyArrPath[0]] = value
			returnValue = nowdata
			return
		case []any:
			pos, err := strconv.Atoi(keyArrPath[0])

			if pos > len(nowdata)-1 {
				returnError = errors.New("路径错误:切片置值位置超过切片长度")
				return
			}
			if pos < 0 {
				returnError = errors.New("路径错误:切片置值不能<0")
				return
			}
			if err != nil {
				returnError = errors.New("切片路径错误:" + err.Error())
				return
			}
			switch newdata := value.(type) {
			case map[string]any:
				nowdata[pos] = newdata
				returnValue = nowdata
				return

			default:
				nowdata[pos] = newdata
			}

			//当前值[位置] = 值
			returnValue = nowdata
			return
		default:
			addedValue := make(map[string]any)
			addedValue[keyArrPath[0]] = data
			returnValue = addedValue
			return

		}
	}
	switch nowdata := data.(type) {
	case map[string]any:
		nowdata[keyArrPath[0]], returnError = setValue(nowdata[keyArrPath[0]], keyArrPath[1:], data)
		return nowdata, returnError
	case []any:
		pos, err := strconv.Atoi(keyArrPath[0])
		if pos > len(data.([]any)) {
			returnError = errors.New("路径错误:切片置值位置超过切片长度")
			return
		}
		if err != nil {
			returnError = errors.New("切片路径错误:" + err.Error())
			return
		}
		nowdata[pos], returnError = setValue(nowdata[pos], keyArrPath[1:], data)
		return nowdata, returnError

	default:

		addedvalue := make(map[string]any)
		addedvalue[keyArrPath[0]] = keyArrPath[1]
		nowdata = addedvalue
		nowdata.(map[string]any)[keyArrPath[0]], returnError = setValue(addedvalue[keyArrPath[0]], keyArrPath[1:], data)
		return nowdata, returnError
	}

}

// 键值列表_子切片添加值
func KVList_SubsliceAddValue(data any, keyarrpath []string, value any) (returnError error) {
	returnError = KVListFilter(value)
	if returnError != nil {
		return
	}
	newData, err := KVListDeepcopy(value)
	if err != nil {
		switch nowData := data.(type) {
		case []any:
			_, returnError = subSliceAddValue(nowData, keyarrpath, value)
			return
		case map[string]any:
			_, returnError = subSliceAddValue(nowData, keyarrpath, value)
			return
		default:
			returnError = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
			return
		}
	}
	switch nowdata := data.(type) {
	case []any:
		_, returnError = subSliceAddValue(nowdata, keyarrpath, newData)
		return
	case map[string]any:
		_, returnError = subSliceAddValue(nowdata, keyarrpath, newData)
		return
	default:
		returnError = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
		return
	}
}

// 子切片添加值
func subSliceAddValue(data any, keyarrpath []string, value any) (returnvalue any, returnerror error) {
	lenKeys := len(keyarrpath)
	if lenKeys == 0 {
		switch nowValue := data.(type) {
		case []any:
			nowValue = append(nowValue, value)
			returnvalue = nowValue
			return
		default:
			returnerror = errors.New("错误:键路径  不存在 或者 不为[]any")
			return
		}
	}
	switch nowValue := data.(type) {
	case map[string]any:
		nowValue[keyarrpath[0]], returnerror = subSliceAddValue(nowValue[keyarrpath[0]], keyarrpath[1:], value)
		return nowValue, returnerror
	case []any:
		pos, err := strconv.Atoi(keyarrpath[0])
		if pos > len(data.([]any)) {
			returnerror = errors.New("路径错误:切片置值位置超过切片长度")
			return
		}
		if err != nil {
			returnerror = errors.New("切片路径错误:" + err.Error())
			return
		}
		nowValue[pos], returnerror = subSliceAddValue(nowValue[pos], keyarrpath[1:], value)
		return nowValue, returnerror
	default:

		returnerror = errors.New("错误:键路径 '" + keyarrpath[0] + "' 不存在 或者 不为[]any | map[string]any")
		return nowValue, returnerror
	}

}

// 键值列表_取值
func KVList_GetValue(data any, keyarrpath []string) (returnvalue any, returnerror error) {
	switch nowdata := data.(type) {
	case []any:
		returnvalue, returnerror = getValue(nowdata, keyarrpath)

	case map[string]any:
		returnvalue, returnerror = getValue(nowdata, keyarrpath)

	default:
		returnerror = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
		return
	}
	value, err := KVListDeepcopy(returnvalue)
	if err != nil {
		return
	}
	returnvalue = value
	return

}

// 取值
func getValue(data any, keyarrpath []string) (returnvalue any, returnerror error) {
	lenKeys := len(keyarrpath)
	if lenKeys == 0 {
		returnerror = errors.New("路径错误:路径长度不能为0")
		return
	}
	if lenKeys == 1 {
		switch nowvalue := data.(type) {
		case map[string]any:
			value, ok := nowvalue[keyarrpath[0]]
			if ok {
				returnvalue = value
			} else {
				returnerror = errors.New("错误:键路径 '" + keyarrpath[0] + "' 不存在")
			}

			return
		case []any:
			pos, err := strconv.Atoi(keyarrpath[0])

			if pos > len(nowvalue)-1 {
				returnerror = errors.New("路径错误:切片置值位置超过切片长度")
				return
			}
			if pos < 0 {
				returnerror = errors.New("路径错误:切片置值不能<0")
				return
			}
			if err != nil {
				returnerror = errors.New("路径错误:" + err.Error())
				return
			}
			returnvalue = nowvalue[pos]
			return
		default:
			returnerror = errors.New("路径错误:非有效路径")
			return

		}
	}
	switch nowvalue := data.(type) {
	case map[string]any:

		value, ok := nowvalue[keyarrpath[0]]
		if !ok {
			returnerror = errors.New("错误:键路径 '" + keyarrpath[0] + "' 不存在")
			return
		}
		returnvalue, returnerror = getValue(value, keyarrpath[1:])
		return
	case []any:
		pos, err := strconv.Atoi(keyarrpath[0])
		if pos > len(nowvalue)-1 {
			returnerror = errors.New("路径错误:切片置值位置超过切片长度")
			return
		}
		if pos < 0 {
			returnerror = errors.New("路径错误:切片置值不能<0")
			return
		}

		if err != nil {
			returnerror = errors.New("路径错误:" + err.Error())
			return
		}
		returnvalue, returnvalue = getValue(nowvalue[pos], keyarrpath[1:])
		return
	default:
		returnerror = errors.New("路径错误:非有效路径")
		return
	}

}

// []any 或者 map[string]any
// 键值列表_删除值
func KVList_DelValue(data any, keyarrpath []string) (returnvalue any, returnerror error) {

	switch nowValue := data.(type) {
	case []any:
		returnvalue, returnerror = delvalue(nowValue, keyarrpath)
		return
	case map[string]any:
		returnvalue, returnerror = delvalue(nowValue, keyarrpath)
		return
	default:
		returnvalue = data
		returnerror = errors.New("只支持 []any 或者 map[string]any 当前数据类型错误")
		return
	}

}

// 删除值
func delvalue(data any, keyarrpath []string) (returnvalue any, returnerror error) {
	lenKeys := len(keyarrpath)
	if lenKeys == 0 {
		returnerror = errors.New("路径错误:路径长度不能为0")
		return
	}
	if lenKeys == 1 {
		switch nowValue := data.(type) {
		case map[string]any:
			_, ok := nowValue[keyarrpath[0]]
			if ok {
				delete(nowValue, keyarrpath[0])
				returnvalue = nowValue
			} else {
				returnerror = errors.New("错误:键路径 '" + keyarrpath[0] + "' 不存在")
			}
			return
		case []any:
			pos, err := strconv.Atoi(keyarrpath[0])
			if pos > len(nowValue)-1 {
				returnerror = errors.New("路径错误:切片置值位置超过切片长度")
				return
			}
			if pos < 0 {
				returnerror = errors.New("路径错误:切片置值不能<0")
				return
			}

			if err != nil {
				returnerror = errors.New("路径错误:" + err.Error())
				return
			}
			//fmt.Println(len(当前值[:位置]))
			//fmt.Println(len(当前值[位置+1:]))
			if len(nowValue) == 1 {
				returnvalue = nowValue[:pos]
				return
			} else if len(nowValue[:pos]) == 0 && len(nowValue[pos+1:]) != 0 {
				returnvalue = nowValue[pos+1:]
				return
			} else if len(nowValue[:pos]) != 0 && len(nowValue[pos+1:]) == 0 {
				returnvalue = nowValue[:pos]
				return
			} else {
				nowValue = append(nowValue[:pos], nowValue[pos+1:]...)
				returnvalue = nowValue
			}
			return
		default:
			returnerror = errors.New("路径错误:非有效路径")
			return

		}
	}
	switch nowvalue := data.(type) {
	case map[string]any:
		value, ok := nowvalue[keyarrpath[0]]
		if !ok {
			returnerror = errors.New("错误:键路径 '" + keyarrpath[0] + "' 不存在")
			return
		}
		returnvalue, returnerror = delvalue(value, keyarrpath[1:])
		return
	case []any:
		pos, err := strconv.Atoi(keyarrpath[0])
		if pos > len(nowvalue)-1 {
			returnerror = errors.New("路径错误:切片置值位置超过切片长度")
			return
		}
		if pos < 0 {
			returnerror = errors.New("路径错误:切片置值不能<0")
			return
		}
		if err != nil {
			returnerror = errors.New("路径错误:" + err.Error())
			return
		}
		returnvalue, returnerror = delvalue(nowvalue[pos], keyarrpath[1:])
		return
	default:
		returnerror = errors.New("路径错误:非有效路径")
		return
	}

}

// func 键值列表_路径分割(路径 string) (返回_值 []string) {
// 	返回_值 = 路径分割(路径, ".", false)
// 	for i, v := range 返回_值 {
// 		//全_文本.PDQ_判断前缀(v, "[")
// 		if len(v) >= 3 && 全_文本.PDQ_判断前缀(v, "[") && 全_文本.PDH_判断后缀(v, "]") {
// 			返回_值[i] = v[1 : len(v)-1]
// 		}
// 	}
// 	return
// }

// func 路径分割(路径 string, sep string, 保持符号 bool) (返回_值 []string) {
// 	//sep := "."
// 	separator := sep
// 	keepQuotes := 保持符号
// 	singleQuoteOpen := false
// 	doubleQuoteOpen := false
// 	var tokenBuffer []string
// 	var ret []string

// 	arr := 全_文本.FGW_分割文本(路径, "")
// 	for _, element := range arr {
// 		matches := false
// 		if separator == element {
// 			matches = true
// 		}

// 		if element == "'" && !doubleQuoteOpen {
// 			if keepQuotes {
// 				tokenBuffer = append(tokenBuffer, element)
// 			}
// 			singleQuoteOpen = !singleQuoteOpen
// 			continue
// 		} else if element == `"` && !singleQuoteOpen {
// 			if keepQuotes {
// 				tokenBuffer = append(tokenBuffer, element)
// 			}
// 			doubleQuoteOpen = !doubleQuoteOpen
// 			continue
// 		}

// 		if !singleQuoteOpen && !doubleQuoteOpen && matches {
// 			if len(tokenBuffer) > 0 {
// 				ret = append(ret, 全_文本.HBW_合并文本_切片(tokenBuffer, ""))
// 				tokenBuffer = make([]string, 0)
// 			} else if sep != "" {
// 				ret = append(ret, element)
// 			}
// 		} else {
// 			tokenBuffer = append(tokenBuffer, element)
// 		}
// 	}
// 	if len(tokenBuffer) > 0 {
// 		ret = append(ret, 全_文本.HBW_合并文本_切片(tokenBuffer, ""))
// 	} else if sep != "" {
// 		ret = append(ret, "")
// 	}
// 	return ret
// }

// 按照条件 和 非条件 分割  被' '包裹的 不会被分割 如果 ”内有'  可以用 " \\' "  或者是 ` \' `
// 键值列表_条件分割
func KVList_ConditionSlice(text string, conditionarr []string) (returnvalue []string) {
	segmented_text := strings.Split(text, "")
	text_arr := ""
	cutoff := false
	toexecuted := false
	for i, v := range segmented_text {
		text_arr = text_arr + v

		if v == "'" {
			cutoff = !cutoff
			toexecuted = true
			if i > 0 && segmented_text[i-1] == `\` {
				cutoff = !cutoff
				toexecuted = false
			}

		}
		if cutoff && toexecuted {
			toexecuted = false
			segmented_text = segmented_text[:len(segmented_text)-1]
			temptext := conditionSplitWithSpace(text_arr, conditionarr, conditionarr)
			nowarr := strings.Fields(temptext)
			returnvalue = append(returnvalue, nowarr...)
			text_arr = "'"
		} else if !cutoff && toexecuted {
			toexecuted = false
			returnvalue = append(returnvalue, text_arr)
			text_arr = ""

		}

	}
	if len(text_arr) > 0 {
		temptext := conditionSplitWithSpace(text_arr, conditionarr, conditionarr)
		nowarr := strings.Fields(temptext)
		returnvalue = append(returnvalue, nowarr...)

	}
	for i, v := range returnvalue {
		returnvalue[i] = strings.Replace(v, `\'`, "'", -1)

	}
	return
}

// 条件分割_加空格
func conditionSplitWithSpace(text string, conditionarr, referencearr []string) (returnvalue string) {
	count := len(conditionarr)
	if count == 0 {
		return
	}
	conditionvalue := conditionarr[0]
	if count == 1 {
		slicearr := strings.Split(text, conditionvalue)
		arr_num := len(slicearr)
		nowarr := make([]string, 0)

		for i, v := range slicearr {
			if i == arr_num-1 {
				nowarr = append(nowarr, v)

			} else {
				nowarr = append(nowarr, v, conditionvalue)
			}

		}
		returnvalue = strings.Join(nowarr, " ")
		return
	}

	slicearr := strings.Split(text, conditionvalue)
	arr_num := len(slicearr)
	nowarr := make([]string, 0)

	for i, v := range slicearr {
		if i == arr_num-1 {
			nowarr = append(nowarr, v)

		} else {
			nowarr = append(nowarr, v, conditionvalue)
		}

	}
	for i, v := range nowarr {
		repeat := false
		for _, X := range referencearr {
			if X == v {
				repeat = true
				break
			}
		}
		if !repeat {
			nowarr[i] = conditionSplitWithSpace(v, conditionarr[1:], referencearr)
		}
	}
	returnvalue = strings.Join(nowarr, " ")

	return
}

// 判断 数组里 ( 和 ) 的成对 位置越前的越有括号优先级 并返回 括号顺序的 位置数组
// 键值列表_查询__括号位置算法
func KVList_GetBracketPos(conditionarr []string) (returnposarr [][]int, returnerror error) {
	returnposarr, returnerror = BracketPos(conditionarr, -1)
	if returnerror != nil {
		return
	}
	exist := false
	for i, v := range conditionarr {
		if v == ")" {
			exist = false
			for _, pos := range returnposarr {
				if i == pos[1] {
					exist = true
				}
			}
			if !exist {
				returnerror = errors.New("错误:括号不对称,多了 ) 号")
				return
			}
		}

	}
	return
}

// 括号位置算法
func BracketPos(conditionarr []string, startLeft int) (returnposarr [][]int, returnerror error) {
	startindex := -1
	startnum := 0
	//结束括号位置 := -1
	endnum := 0
	for i, v := range conditionarr {

		if startindex == -1 && v == "(" {
			if i > startLeft {
				startindex = i
				startnum++
			}
		} else if startindex != -1 && v == "(" {
			if i > startLeft {
				startnum++
			}
		} else if startindex != -1 && v == ")" {
			if i > startLeft {
				endnum++
			}
			if startnum == endnum {
				posarr := []int{startindex, i}
				returnposarr, returnerror = BracketPos(conditionarr, startindex)
				returnposarr = append(returnposarr, posarr)
				return

			}

		}
		// else if 开始括号位置 != -1 && v == ")" && i> 左开始位置{
		// 	fmt.Println(99999,i,左开始位置,v)
		// 	//右括号存在=true
		// }

	}

	if startnum != endnum {

		returnerror = errors.New("错误:括号不对称,少了 ) 号")
		return

	}

	//fmt.Println("开始括号数,结束括号数", 开始括号数, 结束括号数)
	return
}

// 可以将 带有 ()+-*/ 数字 文本  和其他符号的 文本  分割后成数组后 拿进来用  会自动 计算 里面的 +-*/ 并 按优先级处理掉括号
func KVListQueryOperationHandleByOrder(conditionarr []string) (resultarr []string, returnerror error) {
	//fmt.Println(1, 条件组)
	arrlenght := len(conditionarr)
	for i, v := range conditionarr {
		if v == "*" || v == "/" || v == "%" {
			if i == 0 || i == arrlenght {
				returnerror = errors.New("错误: * 位置错误")
				return
			}
			left, err1 := strconv.ParseFloat(conditionarr[i-1], 64)
			right, err2 := strconv.ParseFloat(conditionarr[i+1], 64)
			if err1 != nil && err2 != nil {
				if conditionarr[i-1] == ")" && conditionarr[i+1] == "(" {

				} else {
					returnerror = errors.New("错误: 符号" + v + " 计算类型错误")
					return

				}

			} else if err1 == nil && err2 == nil {
				value := 0.0
				if v == "*" {
					value = left * right
				} else if v == "/" {
					value = left / right
				} else {
					value = float64(int(left) % int(right))
				}

				if i <= 1 && i >= arrlenght-2 {
					newarr := []string{fmt.Sprintf("%v", value)}
					resultarr = newarr
				} else if i <= 1 {
					newarr := make([]string, 0)

					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, conditionarr[i+2:]...)
					resultarr = newarr
				} else if i >= arrlenght-2 {
					newarr := make([]string, 0)
					newarr = append(newarr, conditionarr[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					resultarr = newarr
				} else {
					newarr := make([]string, 0)
					newarr = append(newarr, conditionarr[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, conditionarr[i+2:]...)
					resultarr = newarr
				}

				if len(resultarr) >= 3 && i >= 2 && i+1 <= len(resultarr) {

					if resultarr[i-2] == "(" && resultarr[i] == ")" {
						newarr := make([]string, 0)
						newarr = append(newarr, resultarr[:i-2]...)
						newarr = append(newarr, resultarr[i-1])
						if i+1 <= len(resultarr) {
							newarr = append(newarr, resultarr[i+1:]...)
						}

						resultarr = newarr
					}
				}
				resultarr, returnerror = KVListQueryOperationHandleByOrder(resultarr)
				return

			}

		}

	}

	for i, v := range conditionarr {
		if v == "+" || v == "-" {
			if i == 0 {
				if v == "-" && arrlenght >= 2 {
					newarr := make([]string, 0)
					value := v + conditionarr[1]
					newarr = append(newarr, value)
					if arrlenght >= 3 {
						newarr = append(newarr, conditionarr[i+2:]...)
					}
					resultarr = newarr
					if len(resultarr) >= 3 && i >= 2 && i+1 <= len(resultarr) {
						if resultarr[i-2] == "(" && resultarr[i] == ")" {

							newarr := make([]string, 0)
							newarr = append(newarr, resultarr[:i-2]...)
							newarr = append(newarr, resultarr[i-1])
							if i+1 <= len(resultarr) {
								newarr = append(newarr, resultarr[i+1:]...)
							}

							resultarr = newarr
						}
					}

					resultarr, returnerror = KVListQueryOperationHandleByOrder(resultarr)
					return

				}
				returnerror = errors.New("错误: * 位置错误1")
				return
			}
			if i == arrlenght {
				returnerror = errors.New("错误: * 位置错误2")
				return
			}

			left, err1 := strconv.ParseFloat(conditionarr[i-1], 64)

			right, err2 := strconv.ParseFloat(conditionarr[i+1], 64)
			value := ""
			if err1 != nil && err2 != nil {
				if v == "+" && strings.HasPrefix(conditionarr[i-1], "'") && strings.HasSuffix(conditionarr[i-1], "'") && strings.HasPrefix(conditionarr[i+1], "'") && strings.HasSuffix(conditionarr[i+1], "'") {
					value = conditionarr[i-1][:len(conditionarr[i-1])-1] + conditionarr[i+1][1:len(conditionarr[i+1])]
					// 值 = 条件组[i-1] + 条件组[i+1]
					// 值 = strings.Replace(值, "'", "", -1)
					// 值 = "'" + 值 + "'"
				} else if conditionarr[i-1] == ")" && conditionarr[i+1] == "(" {
					continue

				} else if v == "-" && conditionarr[i+1] == "(" {
					continue
				} else {
					returnerror = errors.New("错误: 符号" + v + " 计算类型错误")
					return
				}

			} else if err1 == nil && err2 == nil {

				if v == "+" {
					value = fmt.Sprintf("%v", left+right)
				} else {
					value = fmt.Sprintf("%v", left-right)
				}
			} else if err1 != nil && err2 == nil && v == "-" {
				newarr := make([]string, 0)
				invalue := v + conditionarr[i+1]
				newarr = append(newarr, conditionarr[:i]...)
				newarr = append(newarr, invalue)
				if arrlenght >= i+3 {
					newarr = append(newarr, conditionarr[i+2:]...)
				}
				resultarr = newarr
				if len(resultarr) >= 3 && i >= 1 && i+2 <= len(resultarr) {

					if resultarr[i-1] == "(" && resultarr[i+1] == ")" {

						newarr := make([]string, 0)
						newarr = append(newarr, resultarr[:i-1]...)
						newarr = append(newarr, resultarr[i])
						if i+2 <= len(resultarr) {
							newarr = append(newarr, resultarr[i+2:]...)
						}

						resultarr = newarr

					}
				}

				resultarr, returnerror = KVListQueryOperationHandleByOrder(resultarr)
				return

			} else if (err1 == nil || err2 == nil) && v == "+" {
				if (strings.HasPrefix(conditionarr[i-1], "'") && strings.HasSuffix(conditionarr[i-1], "'")) || (strings.HasPrefix(conditionarr[i+1], "'") && strings.HasSuffix(conditionarr[i+1], "'")) {
					value1 := ""
					value2 := ""
					if strings.HasPrefix(conditionarr[i-1], "'") && strings.HasSuffix(conditionarr[i-1], "'") {
						value1 = conditionarr[i-1][:len(conditionarr[i-1])-1]
					} else {
						value1 = "'" + conditionarr[i-1]
					}
					if strings.HasPrefix(conditionarr[i+1], "'") && strings.HasSuffix(conditionarr[i+1], "'") {
						value2 = conditionarr[i+1][1:len(conditionarr[i+1])]
					} else {
						value2 = conditionarr[i+1] + "'"
					}
					value = value1 + value2
				} else {
					continue
				}
			} else {

				continue
			}

			if i <= 1 && i >= arrlenght-2 {
				newarr := []string{value}
				resultarr = newarr
			} else if i <= 1 {
				newarr := make([]string, 0)

				newarr = append(newarr, value)
				newarr = append(newarr, conditionarr[i+2:]...)
				resultarr = newarr
			} else if i >= arrlenght-2 {
				newarr := make([]string, 0)
				newarr = append(newarr, conditionarr[:i-1]...)
				newarr = append(newarr, value)
				resultarr = newarr
			} else {
				newarr := make([]string, 0)
				newarr = append(newarr, conditionarr[:i-1]...)
				newarr = append(newarr, value)
				newarr = append(newarr, conditionarr[i+2:]...)
				resultarr = newarr
			}
			if len(resultarr) >= 3 && i >= 2 && i+1 <= len(resultarr) {

				if resultarr[i-2] == "(" && resultarr[i] == ")" {
					newarr := make([]string, 0)
					newarr = append(newarr, resultarr[:i-2]...)
					newarr = append(newarr, resultarr[i-1])
					if i+1 <= len(resultarr) {
						newarr = append(newarr, resultarr[i+1:]...)
					}

					resultarr = newarr
				}
			}

			resultarr, returnerror = KVListQueryOperationHandleByOrder(resultarr)
			return
		}
	}

	resultarr = conditionarr
	if len(resultarr) == 3 {

		if resultarr[0] == "(" && resultarr[2] == ")" {
			newarr := make([]string, 0)
			newarr = append(newarr, resultarr[1])
			resultarr = newarr
			return
		}
	}
	return
}

// 条件数组内 有 独立的  = > < >= <= !=  ~= !~=  符合两边有 文本或者数值  会判断 并用 文本型的 ture 或者 false 替换掉
func KVListQueryEquationLogicHandle(conditionarr []string) (resultarr []string, returnerror error) {
	//fmt.Println(条件组)
	arrLength := len(conditionarr)
	for i, v := range conditionarr {
		if v == "=" || v == ">" || v == ">=" || v == "<" || v == "<=" || v == "!=" || v == "~=" || v == "!~=" {
			if i == 0 || i == arrLength {
				returnerror = errors.New("错误:不能以 " + v + " 开头结尾 位置错误")
				return
			}
			_, err1 := strconv.ParseFloat(conditionarr[i-1], 64)
			if err1 != nil && !strings.HasPrefix(conditionarr[i-1], "'") && !strings.HasSuffix(conditionarr[i-1], "'") {
				returnerror = errors.New("错误:判断符 " + conditionarr[i-1] + v + " 错误")
				return
			}
			_, err2 := strconv.ParseFloat(conditionarr[i+1], 64)
			if err2 != nil && !strings.HasPrefix(conditionarr[i+1], "'") && !strings.HasSuffix(conditionarr[i+1], "'") {
				returnerror = errors.New("错误:判断符 " + conditionarr[i+1] + v + " 错误")
				return
			}

			value := equationLogicHandle(conditionarr[i-1], v, conditionarr[i+1])

			if i <= 1 && i >= arrLength-2 {
				newarr := []string{fmt.Sprintf("%v", value)}
				resultarr = newarr
			} else if i <= 1 {
				newarr := make([]string, 0)

				newarr = append(newarr, fmt.Sprintf("%v", value))
				newarr = append(newarr, conditionarr[i+2:]...)
				resultarr = newarr
			} else if i >= arrLength-2 {
				newarr := make([]string, 0)
				newarr = append(newarr, conditionarr[:i-1]...)
				newarr = append(newarr, fmt.Sprintf("%v", value))
				resultarr = newarr
			} else {
				newarr := make([]string, 0)
				newarr = append(newarr, conditionarr[:i-1]...)
				newarr = append(newarr, fmt.Sprintf("%v", value))
				newarr = append(newarr, conditionarr[i+2:]...)
				resultarr = newarr
			}
			resultarr, returnerror = KVListQueryOperationHandleByOrder(resultarr)
			return

		}

	}
	resultarr = conditionarr
	return
}

// 以 原始文本 的数据类型为基准  进行类型转换->判断文本或数值比较    若类型转换失败 则以文本类型判断
// 如 "原始文本" = "判断文本"  返回判断结果
// 支持  = > < >= <= !=  ~= !~=
func equationLogicHandle(text, judges, judgestext string) (result bool) {

	if judges == "=" {
		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				result = text == judgestext
				return
			} else {
				result = text == ("'" + judgestext + "'")
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				newvalue := judgestext[1 : len(judgestext)-1]
				result = text == newvalue
				// 数值2, err := strconv.ParseFloat(新值, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 ==新值
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 == 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text == value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 == num2

				return
			}
		}

	} else if judges == "!=" {

		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				result = text != judgestext
				return
			} else {
				result = text != ("'" + judgestext + "'")
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				newvalue := judgestext[1 : len(judgestext)-1]
				result = text != newvalue
				// 数值2, err := strconv.ParseFloat(新值, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 !=新值
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 != 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text == value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 == num2

				//返回_值 = 原始文本 != 判断文本
				return
			}
		}
	} else if judges == ">" {

		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 := text[1 : len(text)-1]
				value2 := judgestext[1 : len(judgestext)-1]
				result = value1 > value2
				return
			} else {
				vaule1 := text[1 : len(text)-1]
				result = vaule1 > judgestext
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value2 := judgestext[1 : len(judgestext)-1]
				result = text > value2
				// 数值2, err := strconv.ParseFloat(新值2, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 > 新值2
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 > 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text > value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 > num2
				return
			}
		}
	} else if judges == "<" {

		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 := text[1 : len(text)-1]
				value2 := judgestext[1 : len(judgestext)-1]
				result = value1 < value2
				return
			} else {
				value1 := text[1 : len(text)-1]
				result = value1 < judgestext
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value2 := judgestext[1 : len(judgestext)-1]
				result = text < value2
				// 数值2, err := strconv.ParseFloat(新值2, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 < 新值2
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 < 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text < value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 < num2
				return
			}
		}
	} else if judges == "<=" {

		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 := text[1 : len(text)-1]
				value2 := judgestext[1 : len(judgestext)-1]
				result = value1 <= value2
				return
			} else {
				value1 := text[1 : len(text)-1]
				result = value1 <= judgestext
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value2 := judgestext[1 : len(judgestext)-1]
				result = text <= value2
				// 数值2, err := strconv.ParseFloat(新值2, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 <= 新值2
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 <= 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text <= value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 <= num2
				return
			}
		}
	} else if judges == ">=" {

		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 := text[1 : len(text)-1]
				value2 := judgestext[1 : len(judgestext)-1]
				result = value1 >= value2
				return
			} else {
				value1 := text[1 : len(text)-1]
				result = value1 >= judgestext
				return
			}

		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value2 := judgestext[1 : len(judgestext)-1]
				result = text >= value2
				// 数值2, err := strconv.ParseFloat(新值2, 64)
				// if err != nil {
				// 	返回_值 = 原始文本 >= 新值2
				// 	return
				// }
				// 数值1, _ := strconv.ParseFloat(原始文本, 64)
				// 返回_值 = 数值1 >= 数值2
				return
			} else {
				value2 := judgestext
				num2, err := strconv.ParseFloat(value2, 64)
				if err != nil {
					result = text >= value2
					return
				}
				num1, _ := strconv.ParseFloat(text, 64)
				result = num1 >= num2
				return
			}
		}
	} else if judges == "~=" {
		value1 := ""
		value2 := ""
		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 = text[1 : len(text)-1]
				value2 = judgestext[1 : len(judgestext)-1]
			} else {
				value1 = text[1 : len(text)-1]
				value2 = judgestext
			}
		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 = text
				value2 = judgestext[1 : len(judgestext)-1]
			} else {
				value1 = text
				value2 = judgestext
			}
		}

		if value1 == "" && value2 == "" {
			result = true
			return
		} else if value2 == "" {
			result = false
			return
		}
		if strings.HasPrefix(value2, "%") {
			slices := strings.Split(value2, "%")
			if len(slices) == 2 {
				result = strings.HasSuffix(value1, slices[1])
				return
			} else if len(slices) == 0 {
				result = true
				return
			} else {
				result = false
				return
			}

		} else if strings.HasSuffix(value2, "%") {
			slices := strings.Split(value2, "%")
			if len(slices) == 2 {
				result = strings.HasPrefix(value1, slices[0])
				return
			} else if len(slices) == 0 {
				result = true
				return
			} else {
				result = false
				return
			}

		} else if strings.ContainsAny(value2, "%") {
			slices := strings.Split(value2, "%")

			if len(slices) == 2 {
				result = strings.HasPrefix(value1, slices[0]) && strings.HasSuffix(value1, slices[1])
				return
			} else if len(slices) == 0 {
				result = true
				return
			} else {
				result = false
				return
			}

		} else {
			result = strings.Contains(value1, value2)
			//返回_值 = strings.ContainsAny(新值1, 新值2)
			return
		}

	} else if judges == "!~=" {
		value1 := ""
		value2 := ""
		if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 = text[1 : len(text)-1]
				value2 = judgestext[1 : len(judgestext)-1]
			} else {
				value1 = text[1 : len(text)-1]
				value2 = judgestext
			}
		} else {
			if strings.HasPrefix(judgestext, "'") && strings.HasSuffix(judgestext, "'") {
				value1 = text
				value2 = judgestext[1 : len(judgestext)-1]
			} else {
				value1 = text
				value2 = judgestext
			}
		}

		if value2 == "" && value1 == "" {
			result = false
			return
		} else if value2 == "" {
			result = false
			return
		}
		if strings.HasPrefix(value2, "%") {
			slices := strings.Split(value2, "%")
			if len(slices) == 2 {
				result = !strings.HasSuffix(value1, slices[1])
				return
			} else if len(slices) == 0 {
				result = false
				return
			} else {
				result = false
				return
			}

		} else if strings.HasSuffix(value2, "%") {
			slices := strings.Split(value2, "%")
			if len(slices) == 2 {
				result = !strings.HasPrefix(value1, slices[0])
				return
			} else if len(slices) == 0 {
				result = false
				return
			} else {
				result = false
				return
			}

		} else if strings.ContainsAny(value2, "%") {
			slices := strings.Split(value2, "%")
			if len(slices) == 2 {
				result = !(strings.HasPrefix(value1, slices[0]) && strings.HasSuffix(value1, slices[1]))
				return
			} else if len(slices) == 0 {
				result = false
				return
			} else {
				result = false
				return
			}

		} else {
			result = strings.Index(value1, value2) <= -1
			//返回_值 = !strings.ContainsAny(新值1, 新值2)
			return
		}

	}
	return
}

// 如 [( false !! true ) && true] 这样的文本数组  最后会自动 计算 并返回 [true]  这样的文件结果数组 ,,,返回值 有  [true]  [false]
func KVListQueryTextBoolHandle(condition []string) (resultarr []string, returnerror error) {
	//fmt.Println(条件组)
	arrlength := len(condition)
	for i, v := range condition {

		if v == "&&" {
			if i == 0 || i == arrlength {
				returnerror = errors.New("错误:不能以 " + v + " 开头结尾 位置错误")
				return
			}
			if (condition[i-1] == "true" || condition[i-1] == "false") && (condition[i+1] == "true" || condition[i+1] == "false") {

				value1, err1 := strconv.ParseBool(condition[i-1])
				if err1 != nil {
					err1 = returnerror
					return
				}
				value2, err2 := strconv.ParseBool(condition[i+1])
				if err2 != nil {
					err1 = returnerror
					return
				}
				value := value1 && value2
				if i <= 1 && i >= arrlength-2 {
					newarr := []string{fmt.Sprintf("%v", value)}
					resultarr = newarr
				} else if i <= 1 {
					newarr := make([]string, 0)

					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, condition[i+2:]...)
					resultarr = newarr
				} else if i >= arrlength-2 {
					newarr := make([]string, 0)
					newarr = append(newarr, condition[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					resultarr = newarr
				} else {
					newarr := make([]string, 0)
					newarr = append(newarr, condition[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, condition[i+2:]...)
					resultarr = newarr
				}
				if len(resultarr) >= 3 && i >= 2 && i+1 <= len(resultarr) {

					if resultarr[i-2] == "(" && resultarr[i] == ")" {
						newarr := make([]string, 0)
						newarr = append(newarr, resultarr[:i-2]...)
						newarr = append(newarr, resultarr[i-1])
						if i+1 <= len(resultarr) {
							newarr = append(newarr, resultarr[i+1:]...)
						}

						resultarr = newarr
					}
				}

				resultarr, returnerror = KVListQueryTextBoolHandle(resultarr)
				return

			}

		}

	}

	for i, v := range condition {

		if v == "||" {
			if i == 0 || i == arrlength {
				returnerror = errors.New("错误:不能以 " + v + " 开头结尾 位置错误")
				return
			}
			if (condition[i-1] == "true" || condition[i-1] == "false") && (condition[i+1] == "true" || condition[i+1] == "false") {

				value1, err1 := strconv.ParseBool(condition[i-1])
				if err1 != nil {
					err1 = returnerror
					return
				}
				value2, err2 := strconv.ParseBool(condition[i+1])
				if err2 != nil {
					err1 = returnerror
					return
				}
				value := value1 || value2
				if i <= 1 && i >= arrlength-2 {
					newarr := []string{fmt.Sprintf("%v", value)}
					resultarr = newarr
				} else if i <= 1 {
					newarr := make([]string, 0)

					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, condition[i+2:]...)
					resultarr = newarr
				} else if i >= arrlength-2 {
					newarr := make([]string, 0)
					newarr = append(newarr, condition[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					resultarr = newarr
				} else {
					newarr := make([]string, 0)
					newarr = append(newarr, condition[:i-1]...)
					newarr = append(newarr, fmt.Sprintf("%v", value))
					newarr = append(newarr, condition[i+2:]...)
					resultarr = newarr
				}
				if len(resultarr) >= 3 && i >= 2 && i+1 <= len(resultarr) {

					if resultarr[i-2] == "(" && resultarr[i] == ")" {
						newarr := make([]string, 0)
						newarr = append(newarr, resultarr[:i-2]...)
						newarr = append(newarr, resultarr[i-1])
						if i+1 <= len(resultarr) {
							newarr = append(newarr, resultarr[i+1:]...)
						}

						resultarr = newarr
					}
				}

				resultarr, returnerror = KVListQueryTextBoolHandle(resultarr)
				return

			}

		}

	}
	resultarr = condition
	//fmt.Println("最后结果",结果组)
	if len(resultarr) != 1 {
		returnerror = errors.New("错误:查询语法错误")
		return
	}
	if resultarr[0] != "true" && resultarr[0] != "false" {
		returnerror = errors.New("错误:查询语法错误")
		return
	}

	return
}

// 列表 内部用
func mMapQueryWithParam(condition, slices []string, maps map[string]any) (value bool, returnerror error) {
	value = false
	isexist := false

	for i, v := range slices {
		if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {

			continue
		}
		_, err := strconv.ParseFloat(v, 64)
		if err == nil {

			continue
		}

		isexist = false
		for _, vv := range condition {
			if v == vv {
				isexist = true
				break
			}
		}
		if isexist {

			continue
		}

		res, ok := maps[v]

		if !ok {
			//返回_错误 = errors.New("错误: " + v + " 不存在")
			return
		}
		text := any_to_doc(res)
		_, err1 := strconv.ParseFloat(text, 64)
		if err1 == nil {
			slices[i] = text
		} else {
			slices[i] = "'" + text + "'"
		}
	}

	//fmt.Println(条件切片)
	resSlices, err := KVListQueryOperationHandleByOrder(slices)
	if err != nil {
		returnerror = err
		return
	}
	//fmt.Println(结果切片)
	_, err2 := KVList_GetBracketPos(resSlices)
	if err2 != nil {
		returnerror = err2
		return
	}
	//fmt.Println(返回位置组)

	for _, v := range resSlices {
		if v == "+" || v == "-" || v == "*" || v == "/" || v == "%" {
			returnerror = errors.New("错误:运算符 " + v + " 残余错误")
			return
		}

	}
	boolSlices, err3 := KVListQueryEquationLogicHandle(resSlices)
	if err3 != nil {
		returnerror = err3
		return
	}
	//fmt.Println(逻辑切片)
	resarr, err4 := KVListQueryTextBoolHandle(boolSlices)
	if err4 != nil {
		returnerror = err4
		return
	}
	// //fmt.Println("最后结果",结果组)
	// if len(结果组) != 1 {
	// 	返回_错误 = errors.New("错误:查询语法错误")
	// 	return
	// }
	// if 结果组[0] != "true" && 结果组[0] != "false" {
	// 	返回_错误 = errors.New("错误:查询语法错误")
	// 	return
	// }

	value, returnerror = strconv.ParseBool(resarr[0])
	return
}

func VCPath(pathorindex any, index ...any) (patharr []string, return_error error) {
	indexArr := append([]any{pathorindex}, index...)
	var path_arr = make([]string, len(indexArr))
	for i, v := range indexArr {
		path_arr[i] = fmt.Sprintf("%v", v)
	}
	path := strings.Join(path_arr, ".")
	patharr = strings.Split(path, ".")
	for i, val := range patharr {
		if len(val) >= 3 && strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			patharr[i] = val[1 : len(val)-1]
		}
	}

	if len(patharr) == 0 {
		return_error = errors.New("错误:路径解析错误,路径不能为空")
		return
	}
	return

}
