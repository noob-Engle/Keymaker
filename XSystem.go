package pkg

import (
	"net"
)

type XSystem struct {
}

func (XSystem) QGetMacList() (list LList, returnerr error) {
	list.QClear()
	Interfaces, err := net.Interfaces()
	if err != nil {
		returnerr = err
		return
	}
	for _, v := range Interfaces {
		address := v.HardwareAddr.String()
		list.TAddValue(address)
	}

	return
}

//func (XSystem) QGetCPUID() (value string, returnerror error) {
//	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
//	if runtime.GOOS == "windows" {
//		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏黑框
//	}
//
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		returnerror = err
//		return
//	}
//	str := string(out)
//	reg := regexp.MustCompile(`\s+`)
//	str = reg.ReplaceAllString(str, "")
//	value = str[11:]
//	return
//}
