package db

import (
	"os"
	"bufio"
	"log"
	"SimpleRaft/settings"
	"strings"
	"io"
)

func WriteToDisk(cnt string)  {
	outputFile, outputError := os.OpenFile(settings.DBFILE, os.O_APPEND|os.O_WRONLY, 0666)
	if outputError != nil {
		log.Println("An error occurred with file opening")
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	writer.WriteString(cnt+"\n")
	writer.Flush()
}


func LoadServerIPS(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	var ips []string
	var idx int = -1
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Contains(line,"raft_host") {
			if strings.Contains(line,"raft_host0") {
				idx = len(ips)
			}
			ip := line[:strings.Index(line, " ")]
			ip = strings.TrimSpace(ip)
			ips = append(ips, ip)
		}

		if err != nil {
			if err != io.EOF {
				panic(err)
			}else{
				ips[idx], ips[0] = ips[0], ips[idx]
				return ips
			}
		}
	}
	return nil
}