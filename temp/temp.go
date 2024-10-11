package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// DeleteRecord deletes a record from a CSV file based on the specified value in the first column.
func DeleteRecord(filename, valueToDelete string) error {
	// Step 1: Read all records from the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Step 2: Filter out the record to delete
	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 && record[0] == valueToDelete {
			// Skip the record that matches the value to delete
			continue
		}
		updatedRecords = append(updatedRecords, record)
	}

	// Step 3: Truncate the file before writing the updated records
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Step 4: Write the updated records back to the CSV file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := DeleteRecord("buckets/buckets.csvsdw", "b") // Replace with the actual bucket name you want to delete
	if err != nil {
		fmt.Println("Error deleting record:", err)
		fmt.Println("Error deleting record:" + err.Error())

	} else {
		fmt.Println("Record deleted successfully!")
	}
}
