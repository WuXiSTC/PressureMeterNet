package hsync

import "testing"

func TestHostSync(t *testing.T) {
	hs := NewHostSync("../hosts-1")
	hs.StartSync()
	go hs.AddHost("localhost", "127.0.0.1")
	go hs.AddHost("yin", "127.0.0.1")
	go hs.AddHost("wxstc.com", "192.168.0.1")
	go hs.AddHost("yin-v", "192.168.56.102")
	go hs.AddHost("localhost", "127.0.0.1")
	go hs.DelHost("yin")
	go hs.DelHost("wxstc.com")
	go hs.DelHost("yin-v")
	hs.StopSync()
}
