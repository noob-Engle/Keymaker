package pkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type WFile struct {
}

// QX取运行目录
func (*WFile) QXGetRunDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return dir
}

//调用格式： 〈逻辑型〉 创建目录 （文本型 欲创建的目录名称） - 系统核心支持库->磁盘操作
//英文名称：MkDir
//创建一个新的目录。成功返回真，失败返回假。本命令为初级命令。
//参数<1>的名称为“欲创建的目录名称”，类型为“文本型（text）”。
//
//操作系统需求： Windows、Linux

func (*WFile) CCreateDir(DirName string) error {
	return os.Mkdir(DirName, os.ModePerm)
}

//调用格式： 〈逻辑型〉 删除目录 （文本型 欲删除的目录名称） - 系统核心支持库->磁盘操作
//英文名称：RmDir
//删除一个存在的目录及其中的所有子目录和下属文件，请务必谨慎使用本命令。成功返回真，失败返回假。本命令为初级命令。
//参数<1>的名称为“欲删除的目录名称”，类型为“文本型（text）”。该目录应实际存在，如果目录中存在文件或子目录，将被一并删除，因此使用本命令请千万慎重。
//
//操作系统需求： Windows、Linux

func (*WFile) SDelDir(DirName string) error {
	return os.RemoveAll(DirName)
}

// I复制文件 调用格式： 〈逻辑型〉 复制文件 （文本型 被复制的文件名，文本型 复制到的文件名） - 系统核心支持库->磁盘操作
// 英文名称：FileCopy
// 成功返回真，失败返回假。本命令为初级命令。
// 参数<1>的名称为“被复制的文件名”，类型为“文本型（text）”。
// 参数<2>的名称为“复制到的文件名”，类型为“文本型（text）”。
//
// 操作系统需求： Windows、Linux
func (*WFile) FCopyFile(filename string, targetFile string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)
	dst, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

// 移动文件  源文件filename  目标位置targetPos
func (*WFile) YMoveFile(filename string, targetPos string) error {
	return os.Rename(filename, targetPos)
}

// 调用格式： 〈逻辑型〉 删除文件 （文本型 欲删除的文件名） - 系统核心支持库->磁盘操作
// 英文名称：kill
// 成功返回真，失败返回假。本命令为初级命令。
// 参数<1>的名称为“欲删除的文件名”，类型为“文本型（text）”。
//
// 操作系统需求： Windows、Linux
func (*WFile) SDelFile(filename string) error {
	return os.Remove(filename)
}

// 移动文件  源文件filename  目标位置targetPos
func (*WFile) GChangename(name string, targetname string) error {
	return os.Rename(name, targetname)
}

func (*WFile) PFileIsExist(filename string) bool {
	if stat, err := os.Stat(filename); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func (*WFile) QWGetFileSize(filename string) int64 {
	f, err := os.Stat(filename)
	if err == nil {
		return f.Size()
	} else {
		return -1
	}
}

func (*WFile) DReadFile(filename string) (data []byte, err error) {
	data, err = os.ReadFile(filename)
	if len(data) >= 3 && data[0] == 239 && data[1] == 187 && data[2] == 191 {
		if allText.PHasSuffix(filename, ".txt") || allText.PHasSuffix(filename, ".json") || allText.PHasSuffix(filename, ".sql") {
			data = data[3:]
		}
	}
	return
}

func (*WFile) XWriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, os.ModePerm)
}

func (*WFile) XAddContent(filename string, data []byte) (result error) {
	dst, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		result = err
		return
	}
	defer dst.Close()
	_, err = dst.Write(data)
	if err1 := dst.Close(); err1 != nil && err == nil {
		result = err1
	}
	result = err
	return
}

func (Class *WFile) YExecFile(filename string) (result error) {
	if !Class.PFileIsExist(filename) {
		result = errors.New("错误:文件不存在")
		return
	}
	cmd := exec.Command("cmd", "/c", "start", filename)
	result = cmd.Run()
	return
}

// 文件名 必须是 绝对的 全路径
//func (Class *WFile) YExecFileHideWindow(filename string) (result error) {
//	if !Class.PFileIsExist(filename) {
//		result = errors.New("错误:文件不存在")
//		return
//	}
//	cmd := exec.Command(filename)
//
//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏黑框
//	result = cmd.Run()
//	return
//}

func (Class *WFile) QGetFileArrInDir(path string) (filearr []string, returnerror error) {
	files, err := os.ReadDir(path)
	if err != nil {
		returnerror = err
		return
	}
	filearr = make([]string, len(files))
	for i, v := range files {
		filearr[i] = v.Name()
	}
	return
}
