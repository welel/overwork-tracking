package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

const DataFilePath string = ".overwork_data.json"

// Creates JSON file with data for the program by DataFilePath.
func CreateDataFile() (err error) {
	var file *os.File
	dataFilePath := path.Join(".", DataFilePath)

	if _, err = os.Stat(dataFilePath); err == nil {
		// File exists
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist
		file, err = os.Create(dataFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		newContent, err := json.MarshalIndent(NewData(), "", "\t")
		if err != nil {
			return err
		}

		_, err = file.Write(newContent)
		if err != nil {
			return err
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

	} else {
		return err
	}
	return nil
}

// Loads data from JSON file.
func LoadData() (data *Data, err error) {
	dataFilePath := path.Join(".", DataFilePath)

	file, err := os.Open(dataFilePath)
	if err != nil {
		return data, err
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return data, fmt.Errorf("Data content is invalid at '%v': %v", dataFilePath, err)
	}
	return data, nil
}

// Creates required files and loads data in memory.
func StartupEnvironment() (data *Data, err error) {
	err = CreateDataFile()
	if err != nil {
		return data, err
	}
	data, err = LoadData()
	if err != nil {
		return data, err
	}
	return data, nil
}

func SaveData(data *Data) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	dataFilePath := path.Join(".", DataFilePath)
	tempDataFilePath := path.Join(".", "_temp_"+DataFilePath)

	file, err := os.Create(tempDataFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		panic(err)
	}
	file.Close()

	if err = os.Remove(dataFilePath); err != nil {
		panic(err)
	}

	if err = os.Rename(tempDataFilePath, dataFilePath); err != nil {
		panic(err)
	}
}
