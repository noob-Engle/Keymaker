package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewLList(isSafe ...bool) (returnList LList) {
	if len(isSafe) > 0 && isSafe[0] {
		var lock sync.RWMutex
		returnList.RWPermissions = &lock
	}
	returnList.data = make([]any, 0)
	return
}
func NewLList_DirectAssign(isSafe bool, value ...any) (returnList LList, returnError error) {
	if isSafe {
		var lock sync.RWMutex
		returnList.RWPermissions = &lock
	}
	returnList.data = make([]any, 0)
	returnError = returnList.TAddValue(value...)
	return
}

type LList struct {
	data          []any
	RWPermissions *sync.RWMutex
}

func (Class *LList) init() {
	if Class.data == nil {
		Class.data = make([]any, 0)
	}
}
func (Class *LList) ZXThreadSafety() {
	if Class.RWPermissions == nil {
		var lock sync.RWMutex
		Class.RWPermissions = &lock
	}
}

// @ 支持 json && []any && LList &&[]基本类型 && 可以json化的切片
func (Class *LList) ZRLoad(loadData any) (returnError error) {
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	Class.data = make([]any, 0)
	switch nowData := loadData.(type) {
	case string:
		returnError = json.Unmarshal([]byte(nowData), &Class.data)
		return
	case []any:
		returnError = KVListFilter(nowData)
		if returnError != nil {
			return
		}
		newValue, _ := KVListDeepcopy(nowData)
		Class.data, _ = newValue.([]any)
		return
	case LList:
		if &Class.data == &nowData.data {
			returnError = errors.New("错误:自己不能载入自己")
			return
		}
		Class.data = nowData.Dtoslice()
		return
	case []int:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []bool:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []float64:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []float32:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []int8:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []int16:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []int32:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []int64:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []uint:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []uint16:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []uint32:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []uint64:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []byte:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []time.Time:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	case []string:
		Class.data = listLoadAuxiliaryConversion(nowData)
		return
	default:
		JSON, err := json.Marshal(loadData)
		if err != nil {
			returnError = err
			return
		}
		returnError = json.Unmarshal([]byte(JSON), &Class.data)
		return
	}

}
func listLoadAuxiliaryConversion[T any](args []T) []any {
	returnArr := make([]any, len(args))
	for i, v := range args {
		returnArr[i] = v
	}
	return returnArr

}

// 列分割符 会将 数据分割成 []any  加了表分割符 会分割成[]map[string]any
func (Class *LList) CSplitText(text string, columnsplitter string, tableSplitter ...string) (returnError error) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	Class.data = make([]any, 0)
	tableSymbol := ""
	if len(tableSplitter) > 0 {
		tableSymbol = tableSplitter[0]
	}

	splitArr := allText.FSliceText(text, columnsplitter)

	oldData := make([]any, len(splitArr))
	for i, v := range splitArr {
		if tableSymbol != "" {
			var table JKVtable
			tableArr := allText.FSliceText(v, tableSymbol)
			var transfromArr = make([]any, len(tableArr))
			for ii, vv := range tableArr {
				transfromArr[ii] = vv
			}
			if len(transfromArr) > 0 {
				returnError = table.CCreate(transfromArr...)
				if returnError != nil {
					return
				}
				oldData[i] = table.Dtomap()

			} else {
				oldData[i] = table.Dtomap()
			}

		} else {
			oldData[i] = v
		}

	}

	Class.data = oldData
	return
}
func (Class *LList) Dtoslice() (returnvalue []any) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	newValue, _ := KVListDeepcopy(Class.data)
	returnvalue, _ = newValue.([]any)
	return
}
func (Class *LList) QClear() bool {
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	Class.data = make([]any, 0)
	return true
}
func (Class *LList) DtoJSON() (returnvalue string) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}

	newvalue, _ := KVListDeepcopy(Class.data)
	transformData := KVList_BeforeHandleJson(newvalue)
	JSON, err := json.Marshal(transformData)
	if err != nil {
		returnvalue = "[]"
		return
	}
	returnvalue = string(JSON)
	return
}

func (Class *LList) DtoNewList() (returnvalue LList) {
	returnvalue.QClear()
	Class.init()
	if Class.RWPermissions != nil {
		if &Class.RWPermissions == &returnvalue.RWPermissions {
			return
		}
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	err := returnvalue.ZRLoad(Class.data)
	if err != nil {
		return LList{}
	}
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QGetData(pathorindex any, index ...any) (returnvalue any, returnerror error) {
	patharr, err := VCPath(pathorindex, index...)
	if err != nil {
		returnerror = err
		return
	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}

	returnvalue, returnerror = KVList_GetValue(Class.data, patharr)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QWGetText(pathorindex any, index ...any) (returnvalue string) {
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	value = any_to_doc(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QZGetIndex(pathorindex any, index ...any) (returnvalue int) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToInt(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QXGetDecimal(pathorindex any, index ...any) (returnvalue float64) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToDubble(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QLGetBoolean(pathorindex any, index ...any) (returnvalue bool) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToBoolean(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 LB_列表 则走 路径+索引 混合
func (Class *LList) QMGetmap(pathorindex any, index ...any) (returnvalue map[string]any) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue, _ = value.(map[string]any)
	return

}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QJGetKVTable(pathorindex any, index ...any) (returnvalue JKVtable) {
	returnvalue.QClear()
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	newvalue, OK := value.(map[string]any)
	if !OK {
		return
	}
	returnvalue.ZRLoad(newvalue)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QQGetSlice(pathorindex any, index ...any) (returnvalue []any) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue, _ = value.([]any)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) QLGetList(pathorindex any, index ...any) (returnvalue LList) {
	returnvalue.QClear()
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	switch nowdata := value.(type) {
	case []any:
		returnvalue.ZRLoad(nowdata)
		return
	default:
		//返回_错误 = errors.New("错误:被取值的类型不是[]any")
		return
	}
}
func (Class *LList) QZGetString(pathorindex any, index ...any) (returnvalue []byte) {
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	switch nowvalue := value.(type) {
	case []byte:
		returnvalue = nowvalue
		return
	default:
		returnvalue = allType.DToString(nowvalue)
		return
	}
}

func (Class *LList) QSGetDate(pathorindex any, index ...any) (returnvalue time.Time) {
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	returnvalue = allType.DToTime(value)
	return
}
func (Class *LList) QSGetNumber() int {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	return len(Class.data)
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZSetValue(addvalue any, pathorindex any, index ...any) (returnerror error) {
	patharr, err := VCPath(pathorindex, index...)
	if err != nil {
		returnerror = err
		return
	}
	if newvalue, ok := addvalue.(JKVtable); ok {
		addvalue = newvalue.Dtomap()
	} else if newvalue, ok := addvalue.(LList); ok {
		addvalue = newvalue.Dtoslice()
	}

	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	returnerror = KVList_SetValue(Class.data, patharr, addvalue)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZZSetSubSliceAdd(addvalue any, pathorindex any, index ...any) (returnerror error) {
	patharr, err := VCPath(pathorindex, index...)
	if err != nil {
		returnerror = err
		return
	}
	if newvalue, ok := addvalue.(JKVtable); ok {
		addvalue = newvalue.Dtomap()
	} else if newvalue, ok := addvalue.(LList); ok {
		addvalue = newvalue.Dtoslice()
	}

	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	returnerror = KVList_SubsliceAddValue(Class.data, patharr, addvalue)
	return
}

// 单个 或 按顺序 连续添加
func (Class *LList) TAddValue(addvalue ...any) (returnerror error) {
	for i, v := range addvalue {
		if newvalue, ok := v.(JKVtable); ok {
			assign := newvalue.Dtomap()
			addvalue[i] = assign
		} else if newvalue, ok := v.(LList); ok {
			assign := newvalue.Dtoslice()
			addvalue[i] = assign
		} else {
			returnerror = KVListFilter(v)
			if returnerror != nil {
				return
			}
		}

	}

	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	newvalue, _ := KVListDeepcopy(addvalue)
	newslice, ok := newvalue.([]any)
	if !ok {
		returnerror = errors.New("错误,添加失败")
		return
	}
	Class.data = append(Class.data, newslice...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZWSetText(addValue string, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZZSetIndex(addValue int, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZXSetDecimal(addvalue float64, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZLSetBoolean(addvalue bool, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZZSetString(addvalue []byte, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZSSetDate(addValue time.Time, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZMSetmap(addvalue map[string]any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}
func (Class *LList) ZJSetKVTable(addvalue JKVtable, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *LList) ZQSetSlice(addvalue []any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

func (Class *LList) ZLSetList(addvalue LList, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
// 自动会把 []map[string]any 转换成 []any
func (Class *LList) ZMSetmapArr(addvalue []map[string]any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}
func (Class *LList) SDeleteValue(pathorindex any, index ...any) (returnerror error) {
	patharr, err := VCPath(pathorindex, index...)
	if err != nil {
		returnerror = err
		return
	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	tempvalue, err := KVList_DelValue(Class.data, patharr)
	if value, ok := tempvalue.([]any); ok && err == nil {
		Class.data = value
	}
	return
}
func (Class *LList) ZRLoadRepeatFile(filepath string) (returnerror error) {
	var data []byte
	data, returnerror = allfile.DReadFile(filepath)
	if returnerror != nil {
		return
	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	returnerror = json.Unmarshal(data, &Class.data)
	if returnerror != nil {
		return
	}
	return
}

func (Class *LList) BSaveToFile(filepath string) (returnerror error) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	newvalue, _ := KVListDeepcopy(Class.data)
	transformdata := KVList_BeforeHandleJson(newvalue)
	JSON, err := json.Marshal(transformdata)
	if err != nil {
		returnerror = err
		return
	}
	returnerror = allfile.XWriteFile(filepath, []byte(JSON))
	return
}

// 正常支持 列表里套键值表查询  如[]map[]   []any 类型查询  默认 键值="索引值" 如  索引值="无敌" and 索引值 !="哈哈"
// @ 逻辑判断符号 支持  "||" or OR,    "&& and ADN"
// @ 等式判断符号 支持  "~= LIKE  like (包含 支持 单个站位符号 % 如  键 like '%哈哈')",
// "!~= NOTLIKE  notlike (包含 支持 单个站位符号 % )", "<=", ">=", "!=", "=", ">", "<"
// @ 运算符   "=", "+", "-", "*", "/", "%"(求余数)
// @ 反回新列表
func (List *LList) CQuery(condition string) (returnlist LList, returnerror error) {
	returnlist.QClear()
	returnarr := make([]any, 0)
	splitCondition := []string{"||", "&&", "(", ")", "!~=", "~=", "<=", ">=", "!=", "=", "+", "-", "*", "/", "%", ">", "<"}
	conditionSlice := KVList_ConditionSlice(condition, splitCondition)
	for i, v := range conditionSlice {
		if v == "and" || v == "ADN" {
			conditionSlice[i] = "&&"
		} else if v == "or" || v == "OR" {
			conditionSlice[i] = "||"
		} else if v == "NOTLIKE" || v == "notlike" {
			conditionSlice[i] = "!~="
		} else if v == "LIKE" || v == "like" {
			conditionSlice[i] = "~="
		}

	}

	List.init()
	if List.RWPermissions != nil {
		List.RWPermissions.RLock()
		defer List.RWPermissions.RUnlock()
	}
	for _, v := range List.data {
		addArr := make([]string, len(conditionSlice))
		copy(addArr, conditionSlice)
		switch nowdata := v.(type) {
		case map[string]any:

			res, err := mMapQueryWithParam(splitCondition, addArr, nowdata)
			if err != nil {
				returnerror = err
				return
			}
			if res {
				returnarr = append(returnarr, nowdata)
			}
		default:
			newMap := map[string]any{"索引值": nowdata}
			res, err := mMapQueryWithParam(splitCondition, addArr, newMap)
			if err != nil {
				returnerror = err
				return
			}
			if res {
				returnarr = append(returnarr, newMap)
			}
		}
	}
	returnlist.ZRLoad(returnarr)
	return
}

// []any 类型查询  默认 键值="索引值" 如  索引值="无敌" and 索引值 !="哈哈"
// @ 逻辑判断符号 支持  "||" or OR,    "&& and ADN"
// @ 等式判断符号 支持  "~= LIKE  like (包含 支持 单个站位符号 % 如  键 like '%哈哈')",
// "!~= NOTLIKE  notlike (包含 支持 单个站位符号 % )", "<=", ">=", "!=", "=", ">", "<"
// @ 运算符   "=", "+", "-", "*", "/", "%"(求余数)
// @ 位置  如果 没有对应 则返回 -1
func (Class *LList) CQueryCondition(condition string, startindex ...int) (returnindex int, returnerror error) {
	checkStartIndex := 0
	returnindex = -1
	if len(startindex) > 0 {
		checkStartIndex = startindex[0]
	}
	spiltCondition := []string{"||", "&&", "(", ")", "!~=", "~=", "<=", ">=", "!=", "=", "+", "-", "*", "/", "%", ">", "<"}
	conditionSlice := KVList_ConditionSlice(condition, spiltCondition)
	for i, v := range conditionSlice {
		if v == "and" || v == "ADN" {
			conditionSlice[i] = "&&"
		} else if v == "or" || v == "OR" {
			conditionSlice[i] = "||"
		} else if v == "NOTLIKE" || v == "notlike" {
			conditionSlice[i] = "!~="
		} else if v == "LIKE" || v == "like" {
			conditionSlice[i] = "~="
		}

	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	for i, v := range Class.data {
		addarr := make([]string, len(conditionSlice))
		copy(addarr, conditionSlice)
		if i < checkStartIndex {
			continue
		}
		switch nowdata := v.(type) {
		case map[string]any:
			result, err := mMapQueryWithParam(conditionSlice, addarr, nowdata)
			if err != nil {
				returnerror = err

				return
			}
			if result {
				returnindex = i
				return
			}
		default:
			newmap := map[string]any{"索引值": nowdata}
			result, err := mMapQueryWithParam(spiltCondition, addarr, newmap)
			if err != nil {
				returnerror = err
				return
			}
			if result {
				returnindex = i
				return
			}
		}
	}
	return
}

// @ 主键名称=""的时候 直接查 普通切列表
func (Class *LList) CFind(keyname string, content any, startindex ...int) (returnindex int) {
	checkstartindex := 0
	returnindex = -1
	if len(startindex) > 0 {
		checkstartindex = startindex[0]
	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	if keyname == "" {
		for i, v := range Class.data {
			if i < checkstartindex {
				continue
			}
			if v == content || any_to_doc(v) == any_to_doc(content) {
				returnindex = i
				return
			}

		}
	} else {
		for i, v := range Class.data {
			if i < checkstartindex {
				continue
			}
			switch nowdata := v.(type) {
			case map[string]any:
				if newvalue, ok := nowdata[keyname]; ok && (newvalue == content || any_to_doc(newvalue) == any_to_doc(content)) {
					returnindex = i
					return
				}

			}

		}

	}

	return
}

// 类型 和 数据 完全一致的查找
// @ 主键名称=""的时候 直接查 普通切列表
func (Class *LList) CFindAbsoluteMatch(keyname string, content any, startindex ...int) (returnindex int) {
	checkstartindex := 0
	returnindex = -1
	if len(startindex) > 0 {
		checkstartindex = startindex[0]
	}
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}

	if keyname == "" {
		for i, v := range Class.data {
			if i < checkstartindex {
				continue
			}
			if v == content {
				returnindex = i
				return
			}

		}
	} else {

		for i, v := range Class.data {
			if i < checkstartindex {
				continue
			}
			switch nowdata := v.(type) {
			case map[string]any:
				if newdata, ok := nowdata[keyname]; ok && newdata == content {
					returnindex = i
					return
				}

			}

		}
	}

	return
}

// @唯一标识的键 就是用哪个键的值做索引  要值求唯一性 否则 会被最后一个 值的值替换
// {唯一标识的键的值:{原来列表数据}} 并且会在原数据键值表添加一个  "原列表位置" 的值 标明数据在原来列表的位置
// 只支持 []map[string]any
// 原理列位表的值  oldpos
func (Class *LList) QSGetIndex(uniquekey string) (returnvalue JKVtable, returnerror error) {
	returnvalue.QClear()
	returnvalue.ZXThreadSafety()
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	for i, v := range Class.data {
		switch nowdata := v.(type) {
		case map[string]any:
			if value, ok := nowdata[uniquekey]; ok {
				uniquevalue := any_to_doc(value)
				if uniquevalue != "" {
					nowdata["oldpos"] = i
					returnvalue.ZSetValue(nowdata, uniquevalue)
				}
			} else {
				returnerror = errors.New("错误:唯一标识的键 出现不存在 ")
				return
			}
		default:
			returnerror = errors.New("错误:列表值类型错误  必须是 map[string]any")
			return
		}
	}
	return

}

// L连接到文本
func (Class *LList) LLinkText(mergeconnectors string) (returnvalue string) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	slice := make([]string, len(Class.data))
	for i, v := range Class.data {
		slice[i] = fmt.Sprintf("%v", v)
	}
	returnvalue = strings.Join(slice, mergeconnectors)
	return
}

// 把其它列表的值合并过来  相同的会被替换, 没有的会添加进来
func (Class *LList) HMergeLists(argslist LList) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	list := argslist.Dtoslice()
	Class.data = append(Class.data, list...)
}

// 把其它列表的值合并成一个 新的返回列表  相同的会被替换, 没有的会添加进来(不影响原列表 )
func (Class *LList) HMergeListsToNewList(argslist LList) (returnlist LList, returnerror error) {
	if &returnlist.data == &Class.data || &returnlist.data == &argslist.data {
		returnerror = errors.New("新列表 不能是原有列表")
		return
	}
	returnlist.QClear()
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	newvalue, _ := KVListDeepcopy(Class.data)
	returnlist.ZRLoad(newvalue)
	list := argslist.Dtoslice()
	returnlist.data = append(returnlist.data, list...)
	return
}

// @键参数 排序条件的键  如果键参数="" 则按整切片排序
// @ 排序类型 1=文本 2=整数 3=浮点数 其他数 默认 文本
func (Class *LList) PSort(key string, sorttype int, isup bool) (returnerror error) {
	if sorttype != 2 && sorttype != 3 {
		sorttype = 1
	}
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	sort.SliceStable(Class.data, func(i, j int) bool {
		if key == "" {
			predata := Class.data[i]
			postdata := Class.data[j]
			if sorttype == 2 {
				value1 := allType.DToInt(predata)
				value2 := allType.DToInt(postdata)
				if isup {
					return value1 < value2
				} else {
					return value1 > value2
				}

			} else if sorttype == 3 {
				value1 := allType.DToDubble(predata)
				value2 := allType.DToDubble(postdata)
				if isup {
					return value1 < value2
				} else {
					return value1 > value1
				}

			} else {
				value1 := allType.DtoText(predata)
				value2 := allType.DtoText(postdata)
				if isup {
					return value1 < value2
				} else {
					return value1 > value2
				}

			}

		} else {
			predata := Class.data[i]
			prevalue, ok1 := predata.(map[string]any)
			postdata := Class.data[j]
			postvalue, ok2 := postdata.(map[string]any)
			if !ok1 || !ok2 {
				return false
			}
			if sorttype == 2 {
				value1 := allType.DToInt(prevalue[key])
				value2 := allType.DToInt(postvalue[key])
				if isup {
					return value1 < value2
				} else {
					return value1 > value2
				}

			} else if sorttype == 3 {
				value1 := allType.DToDubble(postvalue[key])
				value2 := allType.DToDubble(postvalue[key])
				if isup {
					return value1 < value2
				} else {
					return value1 > value2
				}

			} else {
				值1 := allType.DtoText(prevalue[key])
				值2 := allType.DtoText(postvalue[key])
				if isup {
					return 值1 < 值2
				} else {
					return 值1 > 值2
				}

			}

		}

	})

	return
}
