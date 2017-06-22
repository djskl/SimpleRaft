package db

import (
	"os"
	"bufio"
	"log"
	"SimpleRaft/settings"
	"net"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Init() {
	exists, err := PathExists(settings.DB_DIR)
	if err != nil {
		panic(nil)
	}

	if !exists {
		os.MkdirAll(settings.DB_DIR, 0777)
	}else{
		os.Remove(settings.DBPATH)
	}

	f, err := os.Create(settings.DBPATH)
	if err != nil{
		panic("数据库文件创建失败！！！")
	}
	defer f.Close()

}

func WriteToDisk(cnt string) {
	outputFile, outputError := os.OpenFile(settings.DBPATH, os.O_APPEND|os.O_WRONLY, 0666)
	if outputError != nil {
		log.Panic("无法打开数据文件!!!")
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	writer.WriteString(cnt + "\n")
	writer.Flush()
}

func LoadLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panic("无法读取本地IP地址!!!")
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if ip != "192.168.16.1" && strings.Contains(ip, "192.168.16"){
					return ip
				}
			}
		}
	}
	return ""
}

func LoadServerIPS(fileName string) []string {
	allips := []string{
		"192.168.16.2",
		"192.168.16.3",
		"192.168.16.4",
		"192.168.16.5",
		"192.168.16.6",
	}
	return allips
	//f, err := os.Open(fileName)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//var ips []string
	//var idx int = -1
	//buf := bufio.NewReader(f)
	//for {
	//	line, err := buf.ReadString('\n')
	//	line = strings.TrimSpace(line)
	//	if strings.Contains(line, "raft_host") {
	//		if strings.Contains(line, "raft_host0") {
	//			idx = len(ips)
	//		}
	//		ip := line[:strings.Index(line, " ")]
	//		ip = strings.TrimSpace(ip)
	//		ips = append(ips, ip)
	//	}
	//
	//	if err != nil {
	//		if err != io.EOF {
	//			panic(err)
	//		} else {
	//			ips[idx], ips[0] = ips[0], ips[idx]
	//			return ips
	//		}
	//	}
	//}
	//return nil
}
