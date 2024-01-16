package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// any_到文本
func any_to_doc(anydata any) (returndata string) {
	if anydata == nil {
		return
	}
	switch nowVal := anydata.(type) {
	case string:
		returndata = nowVal
		return
	case time.Time:
		returndata = nowVal.Format("2006-01-02 15:04:05")
		return
	case []any, []map[string]any:
		JSON, err := json.Marshal(nowVal)
		if err != nil {
			returndata = "[]"
			return
		}
		returndata = string(JSON)
	case map[string]any:
		JSON, err := json.Marshal(returndata)
		if err != nil {
			returndata = "{}"
			return
		}
		returndata = string(JSON)
	case []byte:
		returndata = string(nowVal)
		return
	default:
		returndata = fmt.Sprintf("%v", returndata)
		return
	}
	return
}

// 创建文本
func createText(originalText string, assignTime int, correspondtext ...any) (returndata string) {
	assignTime++
	if assignTime > len(correspondtext) {
		returndata = originalText
		return
	}
	delimiter := "{" + fmt.Sprintf("%v", assignTime) + "}"
	delimitergroup := strings.Split(originalText, delimiter)
	for i, v := range delimitergroup {
		delimitergroup[i] = createText(v, assignTime, correspondtext...)
	}
	returndata = strings.Join(delimitergroup, any_to_doc(correspondtext[assignTime-1]))
	return
}

// 密钥 长度必须 16/24/32长度
// 加密
func encrypt(encryptdoc string, encryptkey string) string {
	if len(encryptdoc) < 1 || (len(encryptkey) != 16 && len(encryptkey) != 24 && len(encryptkey) != 32) {
		return ""
	}

	// 转成字节数组
	origData := []byte(encryptdoc)
	k := []byte(encryptkey)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("密钥 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	return base64.RawURLEncoding.EncodeToString(cryted)

}

// 密钥 长度必须 16/24/32长度
// 解密
func decrypt(decryptdoc string, keys string) string {
	if len(decryptdoc) < 1 || (len(keys) != 16 && len(keys) != 24 && len(keys) != 32) {
		return ""
	}
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	crytedByte, _ := base64.RawURLEncoding.DecodeString(decryptdoc)
	k := []byte(keys)
	if len(crytedByte) == 0 {
		return ""
	} else if len(crytedByte)%16 != 0 && len(crytedByte)%24 != 0 && len(crytedByte)%32 != 0 {
		return ""
	}

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("密钥 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = pkCS7UnPadding(orig)
	return string(orig)
}

// 补码
func pkCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func pkCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
