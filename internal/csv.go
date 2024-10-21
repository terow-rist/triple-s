package internal

import (
	"encoding/csv"
	"os"
	"time"

	"triple-s/config"
)

func writeBucketCSV(bucketName string) error {
	file, err := os.OpenFile(config.Directory+"/buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, err := os.Stat(file.Name()); err != nil {
		return err
	} else if fileInfo.Size() == 0 {
		writer.Write([]string{"Name", "DateOfCreation", "LastModifiedTime", "Status"})
	}

	err = writer.Write([]string{bucketName, time.Now().Format("2006/01/02 15:04:05"), time.Now().Format("2006/01/02 15:04:05"), "MarkedForDeletion"})
	if err != nil {
		return err
	}
	return nil
}

func updateBucketCSV(bucketName string) error {
	file, err := os.OpenFile(config.Directory+"/buckets.csv", os.O_RDWR, 0o644)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 && record[0] != "Name" && record[0] == bucketName {
			newRecord := []string{record[0], record[1]}
			newRecord = append(newRecord, time.Now().Format("2006/01/02 15:04:05"))
			newRecord = append(newRecord, "Active")
			updatedRecords = append(updatedRecords, newRecord)
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

func writeObjectCSV(bucketName string, o ObjectMD) error {
	path := "/" + bucketName + "/objects.csv"
	file, err := os.OpenFile(config.Directory+path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, err := os.Stat(file.Name()); err != nil {
		return err
	} else if fileInfo.Size() == 0 {
		writer.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})
	}
	is, err := elementExists(path, o.ObjectKey)
	if err != nil {
		return err
	}
	if is {
		deleteRecord(path, o.ObjectKey)
	}
	err = writer.Write([]string{o.ObjectKey, o.Size, o.ContentType, o.LastModified})
	if err != nil {
		return err
	}
	return nil
}

func readCSV(path string) ([][]string, error) {
	file, err := os.OpenFile(config.Directory+path, os.O_CREATE, 0o644)
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

func deleteRecord(path, target string) error {
	file, err := os.OpenFile(config.Directory+path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 && record[0] == target {
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

func elementExists(path, element string) (bool, error) {
	records, err := readCSV(path)
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
