package internal

import (
	"encoding/csv"
	"os"
	"time"
)

func writeCSV(bucketName string) error {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, _ := os.Stat(file.Name()); fileInfo.Size() == 0 {
		writer.Write([]string{"NameOfBucket", "DateOfCreation"})
	}

	err = writer.Write([]string{bucketName, time.Now().Format("2006/01/02 15:04:05")})
	if err != nil {
		return err
	}
	return nil
}

func readCSV() ([][]string, error) {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func deleteRecord(target string) error {
	file, err := os.OpenFile("buckets/buckets.csv", os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 && record[0] != "NameOfBucket" && record[0] == target {
			continue
		}
		updatedRecords = append(updatedRecords, record)
	}

	if file.Truncate(0) != nil {
		return err
	}
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return err
	}
	return nil
}

func elementExists(element string) (bool, error) {
	records, err := readCSV()
	if err != nil {
		return false, err
	}
	for _, row := range records {
		for _, bucket := range row {
			if bucket == element {
				return true, nil
			}
		}
	}
	return false, nil
}
