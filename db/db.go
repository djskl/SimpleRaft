package db

import (
	"os"
	"bufio"
	"log"
)

func WriteToDisk(cnt string)  {
	outputFile, outputError := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println("An error occurred with file opening")
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	writer.WriteString(cnt+"\n")
	writer.Flush()
}