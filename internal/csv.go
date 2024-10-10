package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func writeCSV(bucketName string) {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, _ := os.Stat(file.Name()); fileInfo.Size() == 0 {
		writer.Write([]string{"NameOfBucket", "DateOfCreation"})
	}

	err = writer.Write([]string{bucketName, time.Now().Format("2006/01/02 15:04:05")})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func readCSV() [][]string {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return records
}

func elementExists(element string) bool {
	for _, row := range readCSV() {
		for _, bucket := range row {
			if bucket == element {
				return true
			}
		}
	}
	return false
}
