package download

import (
	"fmt"
	"time"
)

// GetCurrentTime gets the current time and date
func GetCurrentTime() string {
	// Get the current time
	currentTime := time.Now()

	// Format time to print up to seconds
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return formattedTime
}

// Simulate will try to simulate the download of the resource the same way wget does it
func Simulate(URL string) {
	start := "start at " + GetCurrentTime()
	fmt.Println(start)
}
