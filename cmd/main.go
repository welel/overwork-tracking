package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

const (
	DataFilePath        string = "data.json"
	WorkHourInputFormat string = "HH:MM"
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

// Creates JSON file with data for the program by DataFilePath.
func CreateDataFile() (err error) {
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

// Loads data from JSON file.
func LoadData() (err error) {
	dataFilePath := path.Join(".", DataFilePath)

	file, err := os.Open(dataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return err
	}
	return nil
}

// Creates required files and loads data in memory.
func StartupEnvironment() (err error) {
	err = CreateDataFile()
	if err != nil {
		return
	}
	err = LoadData()
	if err != nil {
		return
	}
	return nil
}

func SaveData() {
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

func ShowMainScreen() {
	fmt.Printf("---\nWork Today: %s\nOverwork: %s\n\n", data.NeedWork, data.Overwork)
	fmt.Println("1. Record Working Hours")
	fmt.Println("2. Change Need Work")
	fmt.Println("3. Print History")
	fmt.Println("---")
	fmt.Print("Select an option: ")
}

func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ReturnToMainScreenOption(s string) {
	fmt.Printf("\n%s\n", s)
	fmt.Println("-> Press Enter to return to the main screen")
	rd := bufio.NewReader(os.Stdin)
	rd.ReadString('\n')
}

func ScanDurationFromStdin(s string, d *time.Duration) {
	var h, m int
	fmt.Println(s)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan() // Skip empty line
	for {
		sc.Scan()
		if _, err := fmt.Sscanf(sc.Text(), "%d:%d", &h, &m); err != nil {
			fmt.Println("Wrong format! Input in this format: HH:MM")
		} else if h > 24 || h < 0 || m > 59 || m < 0 {
			fmt.Println("Wrong format! HH must be from 00 to 24 and MM from 00 to 59.")
		} else {
			break
		}
	}
	*d = time.Duration(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
}

func RecordWorkingHours() {
	var workedDuration time.Duration
	ScanDurationFromStdin("Enter hours worked today (format: '09:15'):", &workedDuration)
	historicalRecord := HistoryRecord{
		Date:     time.Now(),
		Worked:   workedDuration,
		NeedWork: data.NeedWork,
		Overwork: workedDuration - data.NeedWork,
	}
	if len(data.History) > 0 && IsSameDate(data.History[len(data.History)-1].Date, time.Now()) {
		data.Overwork -= data.History[len(data.History)-1].Overwork
		data.History[len(data.History)-1] = historicalRecord
	} else {
		data.History = append(data.History, historicalRecord)
	}
	data.Overwork += historicalRecord.Overwork
	go SaveData()
	ReturnToMainScreenOption("Worked hours are recorded.")
}

func ChangeNeedWork() {
	var needWorkedDuration time.Duration
	ScanDurationFromStdin("Enter required work hours for today (format: '09:11'):", &needWorkedDuration)
	data.NeedWork = needWorkedDuration
	go SaveData()
	ReturnToMainScreenOption("Today's need work time is changed.")
}

func PrintHistory() {
	panic("Not implemented")

}

func main() {
	err := StartupEnvironment()
	if err != nil {
		fmt.Printf("Can't startup a program: %v\n", err)
		os.Exit(1)
	}

	var option int
	for {
		ShowMainScreen()
		if _, err = fmt.Scan(&option); err != nil {
			option = 0
		}
		switch option {
		case 1:
			RecordWorkingHours()
		case 2:
			ChangeNeedWork()
		case 3:
			PrintHistory()
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
