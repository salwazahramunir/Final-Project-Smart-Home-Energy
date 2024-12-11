package service

import (
	"errors"
	repository "smart-home-energy/internal/repository/fileRepository"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if fileContent == "" {
		return nil, errors.New("CSV file is empty, no data found")
	}

	splitData := strings.Split(fileContent, "\n")

	result := map[string][]string{}

	headers := strings.Split(splitData[0], ",")
	body := splitData[1:]

	for i, header := range headers {
		if _, ok := result[header]; !ok {
			result[header] = []string{}
		}

		for _, row := range body {
			splitRow := strings.Split(row, ",")

			if len(splitRow) != len(headers) {
				return nil, errors.New("CSV required fields are missing")
			}

			result[header] = append(result[header], splitRow[i])
		}

	}

	return result, nil // TODO: replace this
}
