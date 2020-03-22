package hsync

import (
	"fmt"
	"testing"
	"time"
)

func TestClientIPSync(t *testing.T) {
	cis := NewClientIPSync("../hosts-1",
		func(u uint64) string {
			return fmt.Sprintf("jmeter%02d.test.wxstc.org.cn", u)
		})
	cis.StartSync()
	go cis.AddClient("abc", "127.0.0.1")
	go cis.AddClient("abc", "127.0.0.2")
	go cis.AddClient("abcd", "127.0.0.3")
	go cis.AddClient("abcde", "127.0.0.4")
	go cis.DelClient("abc")
	go cis.AddClient("abcdef", "127.0.0.5")
	go cis.DelClient("abcde")
	time.Sleep(1e9)
	cis.StopSync()
}
