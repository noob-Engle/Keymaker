package pkg

import "C"
import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"
)

type JKVtable struct {
	data          map[string]any
	RWPermissions *sync.RWMutex
}

func NewJKVTable(isSafe ...bool) (returnKVTable JKVtable) {
	if len(isSafe) > 0 && isSafe[0] {
		var lock sync.RWMutex
		returnKVTable.RWPermissions = &lock
	}
	returnKVTable.data = make(map[string]any)
	return
}

func NewjkvtableDirectassign(isSafe bool, key ...any) (returnKVTable JKVtable, returnError error) {
	if isSafe {
		var lock sync.RWMutex
		returnKVTable.RWPermissions = &lock
	}
	returnKVTable.data = make(map[string]any)
	returnError = returnKVTable.LContinuousAssign(key...)
	return
}

func (Class *JKVtable) init() {
	if Class.data == nil {
		Class.data = make(map[string]any)
	}
}

func (Class *JKVtable) ZXThreadSafety() {
	if Class.RWPermissions == nil {
		var lock sync.RWMutex
		Class.RWPermissions = &lock
	}
}

// ZRLoad @ 支持 json(字节集的json) && map[string]any && J键值表 && 可以json化的map
func (Class *JKVtable) ZRLoad(loadData any) (returnError error) {
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	Class.data = make(map[string]any)
	switch nowData := loadData.(type) {
	case string:
		returnError = json.Unmarshal([]byte(nowData), &Class.data)
		return
	case []byte:
		returnError = json.Unmarshal(nowData, &Class.data)
		return
	case map[string]any:
		returnError = KVListFilter(nowData)
		if returnError != nil {
			return
		}
		newData, err := KVListDeepcopy(nowData)
		Class.data, _ = newData.(map[string]any)
		returnError = err
		return
	case JKVtable:
		if &Class.data == &nowData.data {
			returnError = errors.New("错误:自己不能载入自己")
			return
		}
		Class.data = nowData.Dtomap()
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

func (Class *JKVtable) Dtomap() (returnValue map[string]any) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	newData, _ := KVListDeepcopy(Class.data)
	returnValue, _ = newData.(map[string]any)
	return
}

func (Class *JKVtable) QClear() bool {
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	Class.data = make(map[string]any)
	return true
}
func (Class *JKVtable) DtoJSON() (returnValue string) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}

	newData, _ := KVListDeepcopy(Class.data)
	transfromData := KVList_BeforeHandleJson(newData)
	JSON, err := json.Marshal(transfromData)
	if err != nil {
		returnValue = "{}"
		return
	}
	returnValue = string(JSON)
	return
}
func (Class *JKVtable) DToNewKVTable() (returnValue JKVtable) {
	returnValue.QClear()
	Class.init()
	if Class.RWPermissions != nil {
		if &Class.RWPermissions == &returnValue.RWPermissions {
			return
		}
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	returnValue.ZRLoad(Class.data)
	return
}

// QGetData 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QGetData(PathOrIndex any, Index ...any) (returnData any, returnError error) {
	PathArr, err := VCPath(PathOrIndex, Index...)
	if err != nil {
		returnError = err
		return
	}

	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}

	returnData, returnError = KVList_GetValue(Class.data, PathArr)
	return
}

// QWGetText 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QWGetText(PathorIndex any, Index ...any) (returnvalue string) {
	value, returnerror := Class.QGetData(PathorIndex, Index...)
	if returnerror != nil {
		return
	}
	returnvalue = any_to_doc(value)
	return
}

// QZGetIndex 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QZGetIndex(pathorindex any, index ...any) (returnvalue int) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToInt(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QXGetDecimal(pathorindex any, index ...any) (returnvalue float64) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToDubble(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QLGetBoolean(pathorindex any, index ...any) (returnvalue bool) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue = allType.DToBoolean(value)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 J键值表 则走 路径+索引 混合
func (Class *JKVtable) QMGetMap(pathorindex any, index ...any) (returnvalue map[string]any) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue, _ = value.(map[string]any)
	return

}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QJGetKVTable(pathorindex any, index ...any) (returnvalue JKVtable) {
	returnvalue.QClear()
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	newData, OK := value.(map[string]any)
	if !OK {
		return
	}
	returnvalue.ZRLoad(newData)
	return
}

// 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QQGetSlice(pathorindex any, index ...any) (returnvalue []any) {
	value, _ := Class.QGetData(pathorindex, index...)
	returnvalue, _ = value.([]any)
	return
}

// QLGetList 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) QLGetList(pathorindex any, index ...any) (returnvalue LList) {
	returnvalue.QClear()
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	switch nowData := value.(type) {
	case []any:
		returnvalue.ZRLoad(nowData)
		return
	default:
		//返回_错误 = errors.New("错误:被取值的类型不是[]any")
		return
	}
}
func (Class *JKVtable) QZGetString(pathorindex any, index ...any) (returnvalue []byte) {
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	switch nowdata := value.(type) {
	case []byte:
		returnvalue = nowdata
		return
	default:
		returnvalue = allType.DToString(nowdata)
		return
	}
}

func (Class *JKVtable) QSGetDate(pathorindex any, index ...any) (returnvalue time.Time) {
	value, err := Class.QGetData(pathorindex, index...)
	if err != nil {
		return
	}
	returnvalue = allType.DToTime(value)
	return
}
func (Class *JKVtable) QSGetNum() int {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	return len(Class.data)
}

// QJGetArr @ 1为升序 2为降序 空 或者 其它为不排序
func (Class *JKVtable) QJGetArr(sorttype ...int) []string {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	jc := 0
	keyarr := make([]string, len(Class.data))
	for k := range Class.data {
		keyarr[jc] = k
		jc++
	}
	if len(sorttype) == 0 {
		return keyarr
	} else if sorttype[0] == 1 {
		sort.Strings(keyarr)
	} else if sorttype[0] == 2 {
		sort.Sort(sort.Reverse(sort.StringSlice(keyarr)))
	}
	return keyarr
}

// ZSetValue 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZSetValue(addValue any, pathorindex any, index ...any) (return_error error) {
	path_arr, err := VCPath(pathorindex, index...)
	if err != nil {
		return_error = err
		return
	}
	if newValue, ok := addValue.(JKVtable); ok {
		addValue = newValue.Dtomap()
	} else if newValue, ok := addValue.(LList); ok {
		addValue = newValue.Dtoslice()
	}

	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}

	return_error = KVList_SetValue(Class.data, path_arr, addValue)
	return
}

// ZZSetSUbSliceAdd 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZZSetSUbSliceAdd(addvalue any, pathorindex any, index ...any) (returnerror error) {
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
func (Class *JKVtable) LContinuousAssign(value ...any) (return_error error) {
	Class.init()
	if len(value)%2 != 0 {
		return_error = errors.New("错误:键值必须为一键 一值 的 双数")
	}
	var key string
	for i, v := range value {
		if (i+1)%2 != 0 {
			switch nowValue := v.(type) {
			case string:
				key = nowValue
			default:
				return_error = errors.New("错误:键值必须为 string")
				return
			}
		} else {
			return_error = Class.ZSetValue(v, key)
			if return_error != nil {
				return
			}
		}
	}
	return
}
func (Class *JKVtable) CCreate(value ...any) (returnerror error) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		Class.data = make(map[string]any)
		Class.RWPermissions.Unlock()
	} else {
		Class.data = make(map[string]any)
	}
	if len(value)%2 != 0 {
		returnerror = errors.New("错误:键值必须为一键 一值 的 双数")
		return
	}
	var key string
	for i, v := range value {
		if (i+1)%2 != 0 {
			switch nowvalue := v.(type) {
			case string:
				key = nowvalue
			default:
				returnerror = errors.New("错误:键值必须为 string")
				return
			}
		} else {
			returnerror = Class.ZSetValue(v, key)
			if returnerror != nil {
				return
			}
		}
	}

	return
}

// ZWSetText 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZWSetText(addValue string, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// ZZSetInt 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZZSetInt(addValue int, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// ZXSet 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZXSet(addValue float64, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// ZLSetBoolean 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZLSetBoolean(addValue bool, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addValue, pathorindex, index...)
	return
}

// ZZSetString 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZZSetString(addvalue []byte, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// ZSSetDate 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZSSetDate(addvalue time.Time, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// ZMSetmap 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZMSetmap(addvalue map[string]any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}
func (Class *JKVtable) ZJSetKVTable(addvalue JKVtable, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// ZQSetSlice 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
func (Class *JKVtable) ZQSetSlice(addvalue []any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

func (Class *JKVtable) ZLSetList(addvalue LList, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

// ZMSetMaparr 路径  用 . 分割  自动去除 前后包裹的 []  如 路径1.路径2.[0].路径4 | 路径1.路径2.0.路径4|路径1.[路径2].0.路径4"
// 索引 如果 后面索引不为空 则走 路径+索引 混合
// 自动会把 []map[string]any 转换成 []any
func (Class *JKVtable) ZMSetMaparr(addvalue []map[string]any, pathorindex any, index ...any) (returnerror error) {
	returnerror = Class.ZSetValue(addvalue, pathorindex, index...)
	return
}

func (Class *JKVtable) SDeleteValue(pathorindex any, index ...any) (returnerror error) {
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

	value, err := KVList_DelValue(Class.data, patharr)
	if v, ok := value.(map[string]any); ok && returnerror == nil {
		Class.data = v
	}
	return
}
func (Class *JKVtable) ZRLoadRepeatFile(path string) (returnerror error) {
	var data []byte
	data, returnerror = allfile.DReadFile(path)
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

func (Class *JKVtable) BSaveToFile(path string) (returnerror error) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	value, _ := KVListDeepcopy(Class.data)
	transformData := KVList_BeforeHandleJson(value)
	JSON, err := json.Marshal(transformData)
	if err != nil {
		returnerror = err
		return
	}
	returnerror = allfile.XWriteFile(path, []byte(JSON))
	return
}

// PKeyIsExist 只支 判断持首层键
func (Class *JKVtable) PKeyIsExist(keyname string) (returnvalue bool) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.RLock()
		defer Class.RWPermissions.RUnlock()
	}
	_, returnvalue = Class.data[keyname]
	return
}

// DToForm @ 1为升序 2为降序 空 或者 其它为不排序
func (Class *JKVtable) DToForm(sorttype ...int) (returnvalue string) {
	num := 0
	if len(sorttype) > 0 {
		num = sorttype[0]
	}
	keyarr := Class.QJGetArr(num)
	var list LList
	for _, v := range keyarr {
		submitvalue := Class.QWGetText(v)
		list.TAddValue(v + "=" + submitvalue)
	}
	returnvalue = list.LLinkText("&")
	return
}

// HMergeTable 把其它键值表的值合并过来  相同的会被替换, 没有的会添加进来
func (Class *JKVtable) HMergeTable(argsKVTable JKVtable) {
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	tabledata := argsKVTable.Dtomap()
	for k, v := range tabledata {
		Class.data[k] = v
	}
}

// HMergeTableToNewTable 把其它键值表的值合并成一个 新的返回表  相同的会被替换, 没有的会添加进来(不影响原表 )
func (Class *JKVtable) HMergeTableToNewTable(argsKVTable JKVtable) (table JKVtable, returnerror error) {
	if &table.data == &Class.data || &table.data == &argsKVTable.data {
		returnerror = errors.New("新表 不能是原有表")
		return

	}

	table.QClear()
	Class.init()
	if Class.RWPermissions != nil {
		Class.RWPermissions.Lock()
		defer Class.RWPermissions.Unlock()
	}
	newvalue, _ := KVListDeepcopy(Class.data)
	table.ZRLoad(newvalue)
	tabledata := argsKVTable.Dtomap()
	for k, v := range tabledata {
		table.data[k] = v
	}
	return
}
