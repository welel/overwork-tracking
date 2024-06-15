package main

import (
	"fmt"
	"os"

	"github.com/welel/overwork-tracking/internal"
)

func main() {
	data, err := internal.StartupEnvironment()
	if err != nil {
		fmt.Printf("Can't startup a program: %v\n", err)
		os.Exit(1)
	}

	var option int
	for {
		internal.ShowMainScreen(data)
		if _, err = fmt.Scan(&option); err != nil {
			option = 0
		}
		switch option {
		case 1:
			internal.RecordWorkingHours(data)
		case 2:
			internal.ChangeNeedWork(data)
		case 3:
			internal.PrintHistory(data)
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
