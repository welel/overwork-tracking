package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"time"
)

const (
	DataFilePath        string = ".overwork_data.json"
	WorkHourInputFormat string = "HH:MM"
)

type WorkTimeDuration time.Duration

type HistoryRecord struct {
	Date     time.Time        `json:"date"`
	Worked   WorkTimeDuration `json:"worked"`
	NeedWork WorkTimeDuration `json:"need_work"`
	Overwork WorkTimeDuration `json:"overwork"`
}

type Data struct {
	NeedWork WorkTimeDuration `json:"need_work"`
	Overwork WorkTimeDuration `json:"overwork"`
	History  []HistoryRecord  `json:"history"`
}

var data Data

func NewData() *Data {
	return &Data{
		NeedWork: 0,
		Overwork: 0,
		History:  []HistoryRecord{},
	}
}

func NewHistoryRecord(workedDuration WorkTimeDuration) *HistoryRecord {
	return &HistoryRecord{
		Date:     time.Now(),
		Worked:   workedDuration,
		NeedWork: data.NeedWork,
		Overwork: workedDuration - data.NeedWork,
	}
}

func (d WorkTimeDuration) String() string {
	totalMinutes := int64(time.Duration(d).Minutes())
	h := int(totalMinutes / 60)
	sign := ""
	if totalMinutes < 0 {
		sign = "-"
	}
	h = int(math.Abs(float64(h)))
	m := int(math.Abs(float64(totalMinutes % 60)))
	return fmt.Sprintf("%s%02d:%02d", sign, h, m)
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
		return fmt.Errorf("Data content is invalid at '%v': %v", dataFilePath, err)
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
	fmt.Printf("---\nWork Today:\t%s\nOverwork:\t%s\n\n", data.NeedWork, data.Overwork)
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

func BlockUntilEnterPressed(s string) {
	fmt.Printf("\n%s\n", s)
	fmt.Println("-> Press Enter to return to the main screen")
	fmt.Scanln()
}

func ScanDuration(s string, d *WorkTimeDuration) {
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
	*d = WorkTimeDuration(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
}

func RecordWorkingHours() {
	var workedDuration WorkTimeDuration
	ScanDuration("Enter hours worked today (format: '09:15'):", &workedDuration)

	historicalRecord := NewHistoryRecord(workedDuration)
	lastIdx := len(data.History) - 1

	if len(data.History) > 0 && IsSameDate(data.History[lastIdx].Date, historicalRecord.Date) {
		data.Overwork -= data.History[lastIdx].Overwork
		data.History[lastIdx] = *historicalRecord
	} else {
		data.History = append(data.History, *historicalRecord)
	}
	data.Overwork += historicalRecord.Overwork
	go SaveData()
	BlockUntilEnterPressed("Worked hours are recorded.")
}

func ChangeNeedWork() {
	var needWorkedDuration WorkTimeDuration
	ScanDuration("Enter required work hours for today (format: '09:11'):", &needWorkedDuration)
	data.NeedWork = needWorkedDuration
	go SaveData()
	BlockUntilEnterPressed("Today's need work time is changed.")
}

func PrintHistory() {
	fmt.Println("\n________________________________________")
	fmt.Println("| Date  | Worked | Need work | Overwork |")
	fmt.Println("|-------+--------+-----------+----------|")
	for _, record := range data.History {
		fmt.Printf(
			"| %s | %s  | %s     | %6s   |\n",
			record.Date.Format("02.01"),
			record.Worked,
			record.NeedWork,
			record.Overwork,
		)
	}
	fmt.Println("|_______|________|___________|__________|")
	BlockUntilEnterPressed("")
	BlockUntilEnterPressed("") // Why first call doesn't stop the program?
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
