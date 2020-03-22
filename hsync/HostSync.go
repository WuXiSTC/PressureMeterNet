package hsync

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//此类用于进行hosts文件的同步
type HostSync struct {
	HostTable
	path   string
	origin []byte //原始文件内容
	syncMu *sync.RWMutex
}

func NewHostSync(path string) *HostSync {
	hs := &HostSync{
		HostTable: NewHostTable(),
		path:      path,
		origin:    []byte{},
		syncMu:    new(sync.RWMutex),
	}
	hs.syncMu.Lock() //停止同步时不允许增删改
	return hs
}

func (hs *HostSync) StartSync() {
	file, err := os.Open(hs.path)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	hs.syncMu.Unlock() //启动同步后允许增删改
	origin, _ := ioutil.ReadAll(file)
	hs.origin = origin
	err = hs.HostTable.Read(hs.path)
	if err != nil {
		log.Println(err)
	}
}

func (hs *HostSync) StopSync() {
	hs.syncMu.Lock()                                                                   //停止同步后不允许增删改
	file, err := os.OpenFile(hs.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm) //打开文件流
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	_, err = file.Write(hs.origin)
	if err != nil {
		panic(err)
	}
}

func (hs *HostSync) AddHost(newHost, newIP string) {
	hs.syncMu.RLock()
	defer hs.syncMu.RUnlock()
	hs.HostTable.AddHost(newHost, newIP)
	err := hs.Write(hs.path)
	if err != nil {
		log.Println(err)
	}
}

func (hs *HostSync) DelHost(oldHost string) {
	hs.syncMu.RLock()
	defer hs.syncMu.RUnlock()
	hs.HostTable.DelHost(oldHost)
	err := hs.Write(hs.path)
	if err != nil {
		log.Println(err)
	}
}
