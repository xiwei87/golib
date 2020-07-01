package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func UriEncode(uri string, encodeSlash bool) string {
	var byte_buf bytes.Buffer
	for _, b := range []byte(uri) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' || (b == '/' && !encodeSlash) {
			byte_buf.WriteByte(b)
		} else {
			byte_buf.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return byte_buf.String()
}

func Md5Str(data string) string {
	h := md5.New()
	h.Write([]byte("hello world"))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func HmacMD5(message string, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewUUID() string {
	var buf [16]byte
	for {
		if _, err := rand.Read(buf[:]); err == nil {
			break
		}
	}
	buf[6] = (buf[6] & 0x0f) | (4 << 4)
	buf[8] = (buf[8] & 0xbf) | 0x80

	res := make([]byte, 36)
	hex.Encode(res[0:8], buf[0:4])
	res[8] = '-'
	hex.Encode(res[9:13], buf[4:6])
	res[13] = '-'
	hex.Encode(res[14:18], buf[6:8])
	res[18] = '-'
	hex.Encode(res[19:23], buf[8:10])
	res[23] = '-'
	hex.Encode(res[24:], buf[10:])
	return string(res)
}

func NewRequestId() string {
	return NewUUID()
}
