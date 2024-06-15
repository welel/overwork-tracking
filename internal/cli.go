package internal

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func ShowMainScreen(data *Data) {
	fmt.Printf("---\nWork Today:\t%s\nOverwork:\t%s\n\n", data.NeedWork, data.Overwork)
	fmt.Println("1. Record Working Hours")
	fmt.Println("2. Change Need Work")
	fmt.Println("3. Print History")
	fmt.Println("---")
	fmt.Print("Select an option: ")
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

func RecordWorkingHours(data *Data) {
	var workedDuration WorkTimeDuration
	ScanDuration("Enter hours worked today (format: '09:15'):", &workedDuration)

	historicalRecord := NewHistoryRecord(data.NeedWork, workedDuration)
	lastIdx := len(data.History) - 1

	if len(data.History) > 0 && IsSameDate(data.History[lastIdx].Date, historicalRecord.Date) {
		data.Overwork -= data.History[lastIdx].Overwork
		data.History[lastIdx] = *historicalRecord
	} else {
		data.History = append(data.History, *historicalRecord)
	}
	data.Overwork += historicalRecord.Overwork
	SaveData(data)
	BlockUntilEnterPressed("Worked hours are recorded.")
}

func ChangeNeedWork(data *Data) {
	var needWorkedDuration WorkTimeDuration
	ScanDuration("Enter required work hours for today (format: '09:11'):", &needWorkedDuration)
	data.NeedWork = needWorkedDuration
	SaveData(data)
	BlockUntilEnterPressed("Today's need work time is changed.")
}

func PrintHistory(data *Data) {
	var prevHist HistoryRecord
	fmt.Println("\n________________________________________")
	fmt.Println("| Date  | Worked | Need work | Overwork |")
	fmt.Println("|-------+--------+-----------+----------|")
	for _, record := range data.History {
		if prevHist.Date.Year() != 1 {
			daysBetween := int(record.Date.Sub(prevHist.Date).Hours()/24) - 1
			for i := 0; i < daysBetween; i++ {
				fmt.Println("|       |        |           |          |")
			}
		}
		fmt.Printf(
			"| %s | %s  | %s     | %6s   |\n",
			record.Date.Format("02.01"),
			record.Worked,
			record.NeedWork,
			record.Overwork,
		)
		prevHist = record
	}
	fmt.Println("|_______|________|___________|__________|")
	BlockUntilEnterPressed("")
	BlockUntilEnterPressed("") // Why first call doesn't stop the program?
}
