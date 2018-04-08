package gokits

import "testing"

func TestKey(t *testing.T) {
	factor := "ywlSRb80TaCQ4b7b"
	t.Logf("%s", KermitStr(factor))
}

func TestCrypto(t *testing.T) {
	scrypto, _ := NewCrypto("ywlSRb80TaCQ4b7b", "7E66")
	usr, _ := scrypto.EncryptStr("goman")
	t.Logf("%s", usr)
	if d, _ := scrypto.DecryptStr(usr); d != "goman" {
		t.Fail()
	}
	pwd, _ := scrypto.EncryptStr("123456")
	t.Logf("%s", pwd)
	if d, _ := scrypto.DecryptStr(pwd); d != "123456" {
		t.Fail()
	}
}
