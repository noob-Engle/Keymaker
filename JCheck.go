package pkg

import (
	"crypto/md5"
	"fmt"
)

type JCheck struct {
}

func (*JCheck) Md5_32Check(content any, isuppercase ...bool) (value string) {
	text := allType.DtoText(content)
	value = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	if len(isuppercase) > 0 && isuppercase[0] {
		value = allText.DDToUppercase(value)
	}
	return
}
func (Class *JCheck) JEncryptandDecrypt_lin(text string, key any, isEncrypt bool) (value string) {
	newkey := Class.Md5_32Check(key, true)
	if isEncrypt {
		value = decrypt(text, newkey)
	} else {
		value = encrypt(text, newkey)
	}

	return
}
