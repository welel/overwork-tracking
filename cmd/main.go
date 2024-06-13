package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

const (
	DataFilePath string = "data.json"
)

type HistoryRecord struct {
	Date     time.Time     `json:"date"`
	Worked   time.Duration `json:"worked"`
	NeedWork time.Duration `json:"need_work"`
	Overwork time.Duration `json:"overwork"`
}

type Data struct {
	NeedWork time.Duration   `json:"need_work"`
	Overwork time.Duration   `json:"overwork"`
	History  []HistoryRecord `json:"history"`
}

var data Data

func NewData() *Data {
	return &Data{
		NeedWork: 0,
		Overwork: 0,
		History:  []HistoryRecord{},
	}
}

func CreateMemoryFile() (err error) {
	var file *os.File
	dataFilePath := path.Join(".", DataFilePath)

	if _, err = os.Stat(dataFilePath); err == nil {
		// File exists
		file, err = os.Open(dataFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

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

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if !json.Valid(fileContent) {
		return fmt.Errorf("Data content is invalid at '%v'", dataFilePath)
	}
	return nil
}

func StartupEnvironment() error {
	err := CreateMemoryFile()
	if err != nil {
		return err
	}
	return nil
}

func SaveData() (err error) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	dataFilePath := path.Join(".", DataFilePath)
	tempDataFilePath := path.Join(".", "_temp_"+DataFilePath)

	file, err := os.Create(tempDataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}
	file.Close()

	if err = os.Remove(dataFilePath); err != nil {
		return err
	}

	if err = os.Rename(tempDataFilePath, dataFilePath); err != nil {
		return err
	}
	return nil
}

func ShowMainScreen() {
	fmt.Printf("---\nWork Today: %s\nOverwork: %s\n\n", data.NeedWork, data.Overwork)
	fmt.Println("1. Record Working Hours")
	fmt.Println("2. Change Work Today")
	fmt.Println("3. Print history")
	fmt.Println("---")
	fmt.Print("Select an option: ")
}

func main() {
	err := StartupEnvironment()
	if err != nil {
		fmt.Printf("Can't startup a program: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Started successfully...")

	var option int
	for {
		ShowMainScreen()
		if _, err = fmt.Scan(&option); err != nil {
			option = 0
		}
		switch option {
		case 1:
			fmt.Println(1)
		case 2:
			fmt.Println(2)
		case 3:
			fmt.Println(3)
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}

// data.History = append(data.History, HistoryRecord{
// 	Date:     time.Now(),
// 	Worked:   0,
// 	NeedWork: 0,
// 	Overwork: 0,
// })
// err = SaveData()
// if err != nil {
// 	fmt.Printf("Failed to save data: %v\n", err)
// 	os.Exit(1)
// }
