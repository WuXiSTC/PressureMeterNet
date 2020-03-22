package hsync

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"
)

//存储站点IP信息
type HostTable struct {
	contents []string                       //此变量存储将要写入文件的内容，IP或注释
	ips      map[string]map[string]struct{} //IP表，一个IP对应多个host
	hosts    map[string]string              //host表，一个host对应一个IP
	mu       *sync.RWMutex
}

func NewHostTable() HostTable {
	return HostTable{
		contents: []string{},
		ips:      map[string]map[string]struct{}{},
		hosts:    map[string]string{},
		mu:       new(sync.RWMutex),
	}
}

//添加一个Host-IP对
func (ht *HostTable) AddHost(newHost, newIP string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.addHost(newHost, newIP)
}
func (ht *HostTable) addHost(newHost, newIP string) {
	if oldIP, ok := ht.hosts[newHost]; ok { //如果此host已存在
		if oldIP == newIP { //并且IP也是一样
			return //就不用改了直接退出
		}
		ipset := ht.ips[oldIP]
		delete(ipset, newHost) //否则就先从IP表中删除这个host
		if len(ipset) < 1 {    //如果表空了
			delete(ht.ips, oldIP) //那就把这个host集合都删了
		}
	}
	ht.hosts[newHost] = newIP
	if oldHosts, ok := ht.ips[newIP]; ok { //如果此IP存在
		oldHosts[newHost] = struct{}{} //就直接添加
	} else { //如果此IP不存在
		ht.ips[newIP] = map[string]struct{}{newHost: {}} //则新建IP表
	}
}

//删除一个Host-IP对
func (ht *HostTable) DelHost(oldHost string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.delHost(oldHost)
}
func (ht *HostTable) delHost(oldHost string) {
	if oldIP, ok := ht.hosts[oldHost]; ok { //如果此host已存在
		if hostset, ok := ht.ips[oldIP]; ok {
			delete(hostset, oldHost) //那就先从IP集合中删除这个host
			if len(hostset) < 1 {    //如果集合空了
				delete(ht.ips, oldIP) //就把这个集合都删了
			}
		}
	}
	delete(ht.hosts, oldHost)
}

var reg, _ = regexp.Compile("^\\s+|\\s+$")
var regx, _ = regexp.Compile("\\s+|#[\\s\\S]*$")

//从正文内容的一行中载入数据，lineNum是行数
func (ht *HostTable) loadAline(line string, lineNum uint64) {
	if length := uint64(len(ht.contents)); length <= lineNum { //如果内容不够长
		ht.contents = append(ht.contents, make([]string, (lineNum+1)-length)...) //就先扩展
	}
	line = reg.ReplaceAllString(line, "") //删除开头和结尾的空字符
	if line == "" || line[0] == '#' {     //如果是空字符串或者第一个非空字符是#（注释）
		ht.contents[lineNum] = line //记录内容
		return                      //直接返回
	}
	line = regx.ReplaceAllString(line, " ") //删除行内多余的空字符和注释
	content := strings.Split(line, " ")     //拆分字符串
	ip, hosts := content[0], content[1:]    //拆出IP和Hosts
	ht.contents[lineNum] = ip               //记录IP
	for _, host := range hosts {            //写入数据
		ht.addHost(host, ip)
	}
}

//通过一个IP和对应Host集合构造文件内容
func getAipContent(ip string, hostset map[string]struct{}) string {
	if len(hostset) < 1 { //没有host
		return "" //返回空
	}
	i, hosts := 0, make([]string, len(hostset))
	for host := range hostset { //构造host列表
		hosts[i] = host
		i++
	}
	return strings.Join(append([]string{ip}, hosts...), " ") //与IP一起拼成文件内容
}

//获取正文内容的一行数据，lineNum是行数
func (ht *HostTable) getAline(lineNum uint64) string {
	if content := ht.contents[lineNum]; content == "" || content[0] == '#' { //如果是注释或空
		return content //就返回注释
	} else if hostset, ok := ht.ips[content]; ok { //否则将此content视为IP。如果此ip存在且有对应hosts
		return getAipContent(content, hostset) //就返回IP列表
	}
	return ""
}

func (ht *HostTable) Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	ht.mu.Lock()
	defer ht.mu.Unlock()
	var i uint64
	for i = 0; scanner.Scan(); i++ {
		ht.loadAline(scanner.Text(), i)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	ht.contents = ht.contents[0:i]
	return nil
}

func (ht *HostTable) Write(path string) error {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm) //打开文件流
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	unwrittenIP := make(map[string]map[string]struct{}, len(ht.ips)) //未写入的IP集合
	for ip, hosts := range ht.ips {                                  //初始化未写入集合
		unwrittenIP[ip] = hosts
	}
	writer := bufio.NewWriter(file)
	for i, content := range ht.contents { //遍历内容表中的每一行
		if content != "" && content[0] != '#' { //如果不是空字符串或注释，说明是IP
			delete(unwrittenIP, content) //已遍历
		}
		_, _ = writer.WriteString(ht.getAline(uint64(i)) + "\n") //写入
	}
	for ip, hosts := range unwrittenIP { //遍历剩余ip
		_, _ = writer.WriteString(getAipContent(ip, hosts) + "\n") //写入
	}
	_ = writer.Flush()
	return nil
}
