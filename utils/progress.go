package utils

import (
	"fmt"
	"time"
)

// ProgressBar zeigt eine einfache ASCII-Fortschrittsanzeige
func ProgressBar(task string, durationSec int) {
	fmt.Printf("%s: [", task)
	for i := 0; i < 50; i++ {
		fmt.Print(" ")
	}
	fmt.Print("]\r")

	for i := 0; i <= 50; i++ {
		fmt.Printf("%s: [", task)
		for j := 0; j < i; j++ {
			fmt.Print("=")
		}
		for j := i; j < 50; j++ {
			fmt.Print(" ")
		}
		fmt.Print("]\r")
		time.Sleep(time.Duration(durationSec) * time.Second / 50)
	}
	fmt.Println()
}
