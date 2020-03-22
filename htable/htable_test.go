package htable

import "testing"

func TestHostTable(t *testing.T) {
	ht := NewHostTable()
	if err := ht.Read("../hosts"); err != nil {
		panic(err)
	}
	ht.AddHost("localhost", "127.0.0.1")
	ht.AddHost("yin", "127.0.0.1")
	ht.AddHost("wxstc.com", "192.168.0.1")
	ht.AddHost("yin-v", "192.168.56.102")
	if err := ht.Write("../hosts-3"); err != nil {
		panic(err)
	}
	ht.DelHost("yin")
	ht.AddHost("localhost", "127.0.0.1")
	ht.AddHost("localhost", "127.0.0.1")
	if err := ht.Write("../hosts-1"); err != nil {
		panic(err)
	}
	if err := ht.Write("../hosts-2"); err != nil {
		panic(err)
	}
	ht.DelHost("localhost")
	ht.DelHost("localhast")
	if err := ht.Write("../hosts-1"); err != nil {
		panic(err)
	}

	ht1 := NewHostTable()
	if err := ht1.Read("../hosts-3"); err != nil {
		panic(err)
	}
	if err := ht1.Write("../hosts-3-1"); err != nil {
		panic(err)
	}
	ht1.DelHost("localhost")
	ht1.DelHost("yin")
	ht1.DelHost("wxstc.com")
	ht1.DelHost("yin-v")
	if err := ht1.Write("../hosts-3-2"); err != nil {
		panic(err)
	}
}
