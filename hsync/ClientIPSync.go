package hsync

import (
	"sync"
)

//此类进一步封装了HostSync，将Client和Host视为一对一关系
//对每一个到达的Client，都读取其广播的IP，并按照格式自动生成一个Host，在hosts文件中将这个IP和Host进行关联
type ClientIPSync struct {
	*HostSync
	end          int64               //此项记录order表中最后一项的下一位在哪
	orderClients []string            //此项记录顺序-ClientID
	clientsOrder map[string]uint64   //此项记录ClientID-顺序
	clientsIP    map[string]string   //此项记录ClientID-IP地址
	formatter    func(uint64) string //此函数用于通过客户端的序号构造出Host
	mu           *sync.Mutex
}

func NewClientIPSync(path string, formatter func(uint64) string) *ClientIPSync {
	return &ClientIPSync{
		HostSync:     NewHostSync(path),
		end:          0,
		orderClients: []string{},
		clientsOrder: map[string]uint64{},
		clientsIP:    map[string]string{},
		formatter:    formatter,
		mu:           new(sync.Mutex),
	}
}

func (cis *ClientIPSync) AddClient(ClientID, IP string) {
	cis.mu.Lock()
	defer cis.mu.Unlock()
	cis.addClient(ClientID, IP)
}
func (cis *ClientIPSync) addClient(ClientID, IP string) {
	order, ok := cis.clientsOrder[ClientID] //先获取顺序
	if !ok {                                //顺序不存在，就要进行添加
		if cis.end >= int64(len(cis.orderClients)) { //如果长度不够
			cis.orderClients = append(cis.orderClients, ClientID) //就扩展
		} else { //否则直接放
			cis.orderClients[cis.end] = ClientID
		}
		cis.clientsOrder[ClientID] = uint64(cis.end) //然后记录顺序
		order = uint64(cis.end)
		cis.end++
	}
	cis.clientsIP[ClientID] = IP                   //记录IP
	cis.HostSync.AddHost(cis.formatter(order), IP) //就然后直接改IP
}

func (cis *ClientIPSync) DelClient(ClientID string) {
	cis.mu.Lock()
	defer cis.mu.Unlock()
	cis.delClient(ClientID)
}
func (cis *ClientIPSync) delClient(ClientID string) {
	if cis.end <= 0 { //如果是空的
		return //那就打扰了
	}
	order, ok := cis.clientsOrder[ClientID] //先获取顺序
	if !ok {                                //顺序不存在
		return //直接退出
	}
	delete(cis.clientsOrder, ClientID)                   //先删除顺序
	delete(cis.clientsIP, ClientID)                      //再删除IP
	cis.end--                                            //末尾指针前移一位
	cis.HostSync.DelHost(cis.formatter(uint64(cis.end))) //直接删除末位
	if uint64(cis.end) == order {                        //如果自己就是最末位
		return //就直接退出
	}
	//不是最末位就用最末位的那个值顶替被删值
	cis.orderClients[order], cis.clientsOrder[cis.orderClients[cis.end]] = cis.orderClients[cis.end], order
	cis.HostSync.AddHost(cis.formatter(order), cis.clientsIP[cis.orderClients[order]]) //就然后直接改IP
}
