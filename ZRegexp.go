package pkg

import (
	"regexp"
)

type ZRegexp struct {
	element *regexp.Regexp
}

func (Class *ZRegexp) CCreate(Condition string) (returnerror error) {
	Class.element, returnerror = regexp.Compile(Condition)
	return
}

// 返回 匹配列表_带子项  [{"0":"匹配项1","1":"子_匹配项1","2":"子_匹配项2"},{"0":"匹配项2","1":"匹配项2的子_匹配项1","2":"匹配项2的_子匹配项2"}]
func (Class *ZRegexp) CCreateAndExec(reg string, text string) (list LList, err error) {
	list.QClear()
	Class.element, err = regexp.Compile(reg)
	if err != nil {
		return
	}
	res := Class.element.FindAllStringSubmatch(text, -1)
	for _, ele := range res {
		table := make(map[string]any)
		for i, v := range ele {
			table[allType.DtoText(i)] = v
		}
		list.TAddValue(table)
	}
	return
}

// 返回 匹配列表_带子项  [{"0":"匹配项1","1":"子_匹配项1","2":"子_匹配项2"},{"0":"匹配项2","1":"匹配项2的子_匹配项1","2":"匹配项2的_子匹配项2"}]
// "0"键的值是全匹配值   后面 "1","2","3"...等键 对应子匹配项的值
func (Class *ZRegexp) QGetListWithSubElement(text string) (list LList) {
	list.QClear()
	res := Class.element.FindAllStringSubmatch(text, -1)
	for _, ele := range res {
		table := make(map[string]any)
		for i, v := range ele {
			table[allType.DtoText(i)] = v
		}
		list.TAddValue(table)
	}
	return
}

// [匹配项1,匹配项2,匹配项3]
func (Class *ZRegexp) QGetList(text string) (list LList) {
	list.QClear()
	res := Class.element.FindAllString(text, -1)
	arr := make([]any, len(res))
	for i, v := range res {
		arr[i] = v
	}
	list.ZRLoad(arr)
	return
}

// 返回分割列表
func (Class *ZRegexp) PGetSlice(text string) (list LList) {
	list.QClear()
	res := Class.element.Split(text, -1)
	list.ZRLoad(res)
	return
}

func (Class *ZRegexp) JCheck(text string) bool {
	return Class.element.MatchString(text)
}

func (Class *ZRegexp) PGetReplace(text, replacetext string) string {
	return Class.element.ReplaceAllString(text, replacetext)
}
