package crypto

import (
	"log"
	"testing"
)

func TestMd5(t *testing.T) {
	t.Log(Md5("123456"))
	t.Log("success")
}

/*
=== RUN   TestHmac256

	crypto_test.go:14: 20b0b3143889f19d4ef957bbdba80edb
	crypto_test.go:16: key:  ac87d721a6171ad1
	crypto_test.go:17: 1b1b91594b40567380e7442eaa60c140

--- PASS: TestHmac256 (0.00s)
PASS
*/
func TestHmac256(t *testing.T) {
	t.Log(Hmac256("123456", ""))
	key := GetIteratorStr(16)
	t.Log("key: ", key)
	t.Log(Hmac256("123456", key))
}

func TestSha256(t *testing.T) {
	t.Log("test Sha256")
	t.Log(Sha256("123456"))
}

// === RUN   TestGetIteratorStr
// 2021/03/11 22:55:01 current str:  64d8f7
// 2021/03/11 22:55:01 current str:  24f978
// 2021/03/11 22:55:01 current str:  249097
// --- PASS: TestGetIteratorStr (0.00s)
// PASS
func TestGetIteratorStr(t *testing.T) {
	for i := 0; i < 1000; i++ {
		log.Println("current str: ", GetIteratorStr(6))
	}
}

var k = GetIteratorStr(16)
var iv = GetIteratorStr(16)

func TestCbc256(t *testing.T) {
	t.Log(AesEncrypt("123456", k, iv))
}

func TestDecodeCbc256(t *testing.T) {
	b, _ := AesEncrypt("123456", k, iv)
	bytes, _ := AesDecrypt(b, k, iv)
	t.Log(string(bytes))
}

/*
*
=== RUN   TestAesEbc
crypto_test.go:54: ebc加密后: 3e75cb8bcd9d5e08
crypto_test.go:57: ebc解密: 123456
--- PASS: TestAesEbc (0.00s)
PASS
*/
func TestAesEbc(t *testing.T) {
	k := GetIteratorStr(8)
	b, _ := EncryptEcb("123456", k)
	t.Log("ebc加密后:", b)

	s, _ := DecryptEcb(b, k)
	t.Log("ebc解密:", s)
}

/*
*
测试aes-256-cbc加密
$ go test -v -test.run=TestAesCbc
=== RUN   TestAesCbc
2021/05/05 18:02:01 /fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ
2021/05/05 18:02:01 HRHtimkjsJktwu6AzH2ji9MP9OLpRBRf35Xcm7zFNmr5Lj8X1rxxJiCIQJqnLC8r
2021/05/05 18:02:01 Sj1ENtUBam7C6PglPZgLZGy9lC8bppce7NS8RExuVa+xWow04Trnlc+kJh+Wz9LL
2021/05/05 18:02:01 中文数字123字母ABC符号!@#$%^&*() <nil>
--- PASS: TestAesCbc (0.00s)
PASS
*/
func TestAesCbc(t *testing.T) {
	str := `中文数字123字母ABC符号!@#$%^&*()`
	k2 := `abcdefghijklmnop`
	iv2 := `1234567890123456`
	b, _ := AesEncrypt(str, k2, iv2)
	log.Println(string(b))

	// log.Println(AesDecrypt(`/fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ`, k2, iv2))

	k2 = `abcdefghijklmnop1234567890123456`
	iv2 = `1234567890123456`
	b, _ = AesEncrypt(str, k2, iv2)
	log.Println(b)

	k2 = `abcdefghijklmnop12345678`
	iv2 = `1234567890123456`
	b, _ = AesEncrypt(str, k2, iv2)
	log.Println(b)

	log.Println(AesDecrypt(b, k2, iv2))
}
