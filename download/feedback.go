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

// Feedback will give user additional information that is happening during the download process
// It takes a variable number of arguments depending on the download that is happennin
func Feedback(URL string) {
	start := "start at " + GetCurrentTime()
	fmt.Println(start)
}
