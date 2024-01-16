package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WRequetsKVlist struct {
	protocolHeader map[string]string
	FResHeader     http.Header //可以直接用 . 路径取值
	FRescook       string
	FResStatus     int
	FResData       []byte
}

func (Class *WRequetsKVlist) init() {
	if Class.protocolHeader == nil {
		Class.protocolHeader = make(map[string]string)
	}
}

func (Class *WRequetsKVlist) ZSetprotocolHeader(protocoltable JKVtable) {
	Class.init()
	table := protocoltable.Dtomap()
	for J, v := range table {
		Class.protocolHeader[J] = fmt.Sprintf("%v", v)
	}
}
func (Class *WRequetsKVlist) QKClear(cookies ...string) {
	Class.init()
	Class.protocolHeader = make(map[string]string)
	Class.FRescook = ""
	Class.FResData = make([]byte, 0)
}

func (Class *WRequetsKVlist) ZSetcookies(cookies ...string) {
	Class.init()
	Class.protocolHeader["cookie"] = strings.Join(cookies, ";")
}

func (Class *WRequetsKVlist) Post(url string, data string) (returnerror error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(data))
	if len(Class.protocolHeader) > 0 {
		for k, v := range Class.protocolHeader {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		returnerror = err
		//resp.Body.Close()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	Class.FResData = body
	if err != nil {
		returnerror = err
		return
	}
	Class.FResHeader = resp.Header
	Class.FResStatus = resp.StatusCode
	cookes := resp.Cookies()
	//网络请求.Fh_cook

	protocolHeaderArr := make([]string, len(cookes))
	for i, v := range cookes {
		protocolHeaderArr[i] = v.Name + "=" + v.Value
	}
	Class.FRescook = strings.Join(protocolHeaderArr, ";")
	return
}
func (Class *WRequetsKVlist) FResData_ToKVTable() (table JKVtable, returnerror error) {
	maptable := make(map[string]any)
	returnerror = json.Unmarshal(Class.FResData, &maptable)
	if returnerror != nil {
		return
	}
	table.ZRLoad(maptable)
	return
}

func (Class *WRequetsKVlist) FResDataToList() (list LList, returnerror error) {
	slices := make([]any, 0)
	returnerror = json.Unmarshal(Class.FResData, &slices)
	if returnerror != nil {
		return
	}
	list.ZRLoad(slices)
	return
}
func (Class *WRequetsKVlist) FResDataToText() (value string) {
	value = string(Class.FResData)
	return
}

func (Class *WRequetsKVlist) Get(url string, data ...string) (returnerror error) {
	client := &http.Client{}
	requestsUrl := url
	if len(data) > 0 && data[0] != "" {
		if strings.HasSuffix(requestsUrl, "?") {
			requestsUrl = requestsUrl + data[0]
		} else {
			requestsUrl = requestsUrl + "?" + data[0]
		}

	}
	req, _ := http.NewRequest("GET", requestsUrl, nil)
	if len(Class.protocolHeader) > 0 {
		for k, v := range Class.protocolHeader {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		returnerror = err
		//resp.Body.Close()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	Class.FResData = body
	if err != nil {
		returnerror = err
		return
	}

	Class.FResHeader = resp.Header
	Class.FResStatus = resp.StatusCode
	Class.FResStatus = resp.StatusCode
	cookes := resp.Cookies()
	protocolHeaderArr := make([]string, len(cookes))
	for i, v := range cookes {
		protocolHeaderArr[i] = v.Name + "=" + v.Value
	}
	Class.FRescook = strings.Join(protocolHeaderArr, ";")
	return
}
