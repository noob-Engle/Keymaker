package pkg

import "strings"

type WText struct {
}

/*
寻找文本
text : 原始文本
str : 寻找文本
返回值为:整型 pos
*/
func (*WText) XFindText(text, str string) (pos int) {
	pos = strings.Index(text, str)
	return
}

/*
X寻找文本位置组_正向找
text : 原始文本
str : 寻找文本
*/
func (*WText) XFindTextPosArrForward(text, str string) (returnval []int) {
	times := strings.Count(text, str)
	pos := 0
	count := 0
	for i := 0; i < times; i++ {
		pos = strings.Index(text[count:], str)

		if pos < 0 {
			return
		}
		count = count + pos
		returnval = append(returnval, count)
		count = count + len(text)
	}
	return
}
func (*WText) XFindTextPosArrBack(text, str string) (retrunVal []int) {
	times := strings.Count(text, str)
	pos := len(text)
	for i := 0; i < times; i++ {
		pos = strings.LastIndex(text[:pos], str)
		if pos < 0 {
			return
		}
		retrunVal = append(retrunVal, pos)
	}
	return
}

// @返回字符串 text 中字符串  寻找文本 最后一次出现的位置
func (*WText) XFindTextPosBack(text, str string) (pos int) {
	pos = strings.LastIndex(text, str)
	return
}

//	文本_取左边
//
// @text 被取的文本
// @寻找文本 关键词
// @启始寻找位置  <0 为 从右往左边 1个汉字一个3个位置
func (Class *WText) QZGetLeft(text, str string, stratIndex ...int) (result string) {

	getStratIndex := 0
	if len(stratIndex) > 0 {
		getStratIndex = stratIndex[0]
	}

	if getStratIndex >= len(text) || -getStratIndex >= len(text) {
		return
	}
	if getStratIndex >= 0 {
		FormerText := text[:getStratIndex]
		Posttext := text[getStratIndex:]
		pos := strings.Index(Posttext, text)
		if pos >= 0 {
			result = FormerText + Posttext[:pos]
		}
		return

	} else if getStratIndex < 0 {
		getStratIndex = len(text) + getStratIndex + 1
		FormerText := text[:getStratIndex]
		pos := strings.LastIndex(FormerText, str)
		if pos >= 0 {
			result = FormerText[:pos]
		}

	}

	return

}

//	文本_取右边
//
// @text 被取的文本
// @寻找文本 关键词
// @启始寻找位置  -1 为 从右往左边 1个汉字一个3个位置
func (Class *WText) QYGetRight(text, str string, startIndex ...int) (result string) {

	getStratIndex := 0
	if len(startIndex) > 0 {
		getStratIndex = startIndex[0]
	}
	if getStratIndex >= len(text) || -getStratIndex >= len(text) {
		return
	}

	if getStratIndex >= 0 {
		backText := text[getStratIndex:]
		pos := strings.Index(backText, str)
		if pos >= 0 {
			result = backText[pos:]
		}

	} else if getStratIndex < 0 {
		getStratIndex = len(text) + getStratIndex + 1
		forwordText := text[:getStratIndex]
		backText := text[getStratIndex:]
		位置 := strings.LastIndex(forwordText, str)
		if 位置 >= 0 {
			result = forwordText[位置:] + backText
		}

	}
	if len(result) >= len(str) {
		result = result[len(str):]
		// 结果组 := strings.Split(结果, "")
		// 结果组=结果组[1:]
		// 结果=strings.Join(结果组, "")
	}

	return

}

func (Class *WText) QZGetMiddle(text, startText, endText string) (result string) {
	resText := Class.QYGetRight(text, startText)
	result = Class.QZGetLeft(resText, endText)
	return
}

func (*WText) QCGetAvilabeCount(text, str string) (count int) {
	count = strings.Count(text, str)
	return
}

func (*WText) QHGetWords(text string) string {
	var Regexp ZRegexp
	Regexp.CCreate("[\u4e00-\u9fa5]+")
	column := Regexp.QGetList(text)
	return column.LLinkText("")

}

// 通过 {1},{2},{3} 进行站为替换对应文本位置
func (*WText) CCreateText(text string, reactText ...any) (returnVal string) {
	returnVal = createText(text, 0, reactText...)
	return
}

// 分割文本
func (*WText) FSliceText(text, str string) (returnSlice []string) {
	if text == "" {
		returnSlice = make([]string, 0)
		return
	}
	returnSlice = strings.Split(text, str)
	return
}
func (*WText) FSplitText_Whitespace(text string) (returnSlice []string) {
	returnSlice = strings.Fields(text)
	return
}

func (*WText) LConcatText_slice(text []string, merge string) (result string) {
	result = strings.Join(text, merge)
	return
}

// @W_文本_子文本替换
//
//	@替换次数 <0 或 大于存在数 为替换全部
//
// @从左往右边  true   重右到左 则填false
func (*WText) TReplaceText_withDirection(text, oldstr, newstr string, count int, lefttoright bool) (returnVal string) {
	existCount := strings.Count(text, oldstr)
	if existCount == 0 {
		return text
	}
	if count < 0 || count >= existCount {
		returnVal = strings.Join(strings.Split(text, oldstr), newstr)
		return
	}
	if count > 0 {
		if !lefttoright {
			returnVal = text
			for i := 0; i < count; i++ {
				pos := strings.LastIndex(returnVal, oldstr)
				returnVal = returnVal[:pos] + newstr + returnVal[pos+len(oldstr):]
			}
			return
		} else {
			returnVal = text
			for i := 0; i < count; i++ {
				pos := strings.Index(returnVal, oldstr)
				returnVal = returnVal[:pos] + newstr + returnVal[pos+len(oldstr):]
			}
			return

		}
	}
	return text
}

// @将字符串 text 前n个不重叠  被替换字符 子串都替换为  替换字符 的新字符串
// @如果n<0会替换所有old子串。
func (*WText) TReplaceText(text, oldstr, newstr string, count int) (returnVal string) {

	returnVal = strings.Replace(text, oldstr, newstr, count)
	return
}
func (*WText) PHasPreix(text, prefixText string) (returnVal bool) {
	returnVal = strings.HasPrefix(text, prefixText)
	return
}
func (*WText) PHasSuffix(text, suffixText string) (returnVal bool) {
	returnVal = strings.HasSuffix(text, suffixText)
	return
}
func (*WText) PIsExist(text, findtext string) (returnVal bool) {
	returnVal = strings.ContainsAny(text, findtext)
	return
}

// @判断字符串s是否包含unicode的码值r
func (*WText) PisExist_unicode(text string, findtext rune) (returnVal bool) {
	returnVal = strings.ContainsRune(text, findtext)
	return
}

// @判断字符串 text 是否包含  寻找文本  字符串中的任意一个字符
func (*WText) PIsContains(text string, findtext string) (returnVal bool) {
	returnVal = strings.ContainsAny(text, findtext)
	return
}

// @判断s和t两个UTF-8字符串是否相等，忽略大小写
func (*WText) PIsSame_utf8(text, compareText string) (returnVal bool) {
	returnVal = strings.EqualFold(text, compareText)
	return
}

/*
@ 按字典顺序比较a和b字符串的大小
@ 如果 a > b，返回一个大于 0 的数
@ 如果 a > b，如果 a == b，返回 0
@ 如果 a < b，返回一个小于 0 的数
*/
func (*WText) PCompareText(text, comparetext string) (returnVal int) {
	returnVal = strings.Compare(text, comparetext)
	return
}

func (*WText) DXToLowercase(text string) (returnVal string) {
	returnVal = strings.ToLower(text)
	return
}
func (*WText) DDToUppercase(text string) (returnVal string) {
	returnVal = strings.ToTitle(text)
	return
}

// @将字符串 text 首尾包含在  条件字符 中的任一字符去掉
func (*WText) STrimStr(text, condtion string) (returnVal string) {
	returnVal = strings.Trim(text, condtion)
	return
}

// 删除首部文本包含
func (*WText) STrimLeft(text, condtion string) (returnVal string) {
	returnVal = strings.TrimLeft(text, condtion)
	return
}

// 删除首部文本
func (*WText) STrimPrefix(text, condtion string) (returnVal string) {
	returnVal = strings.TrimPrefix(text, condtion)
	return
}

// 删除尾部文本包含
func (*WText) STrimRight(text, condtion string) (returnVal string) {
	returnVal = strings.TrimRight(text, condtion)
	return
}

// 删除尾部文本包含
func (*WText) STrimSuffix(text, condtion string) (returnVal string) {
	returnVal = strings.TrimSuffix(text, condtion)
	return
}

// 删除尾部空格
func (*WText) STrimSpace(text string) (returnVal string) {
	returnVal = strings.TrimSpace(text)
	return
}

// 去重复文本
func (*WText) QCDeduplicateText(text string) string {
	slice := strings.Split(text, "")
	returnText := ""
	for _, v := range slice {
		pos := strings.Index(text, v)
		if pos >= 0 {
			continue
		}
		returnText = returnText + v
	}
	return returnText
}

func (*WText) SFIsNumer(text string) bool {
	var Regexp ZRegexp
	Regexp.CCreate(`^[0-9]*$`)
	return Regexp.JCheck(text)
}
