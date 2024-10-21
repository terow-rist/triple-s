package internal

import (
	"encoding/xml"
	"errors"
	"net/http"
)

type Bucket struct {
	Name             string `xml:"Name"`
	DateOfCreation   string `xml:"DateOfCreation"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

type Buckets struct {
	SliceBucket []Bucket `xml:"Bucket"`
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets Buckets  `xml:"Buckets"`
}

type ErrorResponse struct {
	XMLName    xml.Name `xml:"Error"`
	StatusCode string   `xml:"StatusCode"`
	Message    string   `xml:"Message"`
}

type RetrieveAnObject struct {
	XMLName    xml.Name `xml:"RetrieveAnObject"`
	ObjectData string   `xml:"ObjectData"`
}

func listAllMyBucketsResult() ([]byte, error) {
	records, err := readCSV("/buckets.csv")
	if err != nil {
		return nil, err
	}
	buckets := Buckets{}
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			return nil, errors.New("csv record is not correct")
		}
		soloBucket := Bucket{
			Name:             record[0],
			DateOfCreation:   record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets.SliceBucket = append(buckets.SliceBucket, soloBucket)
	}

	result := ListAllMyBucketsResult{
		Buckets: buckets,
	}

	xmlData, err := xml.MarshalIndent(result, "", "   ")
	if err != nil {
		return nil, err
	}
	xmlData = append(xmlData, '\n')
	return xmlData, nil
}

func listObjectData(bucketName, objectKey string) ([]byte, error) {
	records, err := readCSV("/" + bucketName + "/" + "objects.csv")
	if err != nil {
		return nil, err
	}

	var size string
	for _, v := range records {
		if v[0] == objectKey {
			size = v[1]
			break
		}
	}
	objectData := RetrieveAnObject{
		ObjectData: " [" + size + " bytes of object data] ",
	}
	xmlData, err := xml.MarshalIndent(objectData, "", "   ")
	if err != nil {
		return nil, err
	}
	xmlData = append(xmlData, '\n')

	return xmlData, nil
}

func writeXMLError(w http.ResponseWriter, statusCode, message string, code int) {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	}
	xmlData, err := xml.MarshalIndent(errorResponse, "", "   ")
	if err != nil {
		errorResponse.StatusCode = "Internal Server Error"
		errorResponse.Message = "Error: error in xml.MarshalIndent"
	}
	xmlData = append(xmlData, '\n')
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write(xmlData)
}
