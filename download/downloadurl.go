package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)


func DownloadUrl(url string) error {
	if !IsValidURL(url) {
		return fmt.Errorf("invalid-url")
	}

	// Simulate download same way as wget utility
	Simulate(url)

	// get the filename
	files := strings.Split(url, "/")
	filePath := "./" + files[len(files)-1]


	// send http request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching file: %v", err)
	}

	defer response.Body.Close()

	reqMessage := "sending request, awaiting response..."

	if response.StatusCode == http.StatusOK {
		fmt.Printf("\r"+reqMessage+" status %d OK\n",response.StatusCode)
		err = nil
	} else if response.StatusCode == http.StatusNotFound {
		fmt.Printf(reqMessage+" %d: Not Found\n",response.StatusCode)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	contentLength , err := io.Copy(file, response.Body)
	roundOfSize := RoundOfSize(contentLength)
	fmt.Printf("content size: %d %s\n",contentLength,roundOfSize)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}
	return nil
}

func RoundOfSize(n int64) string {
	sizeFloat := float64(n)
	res := sizeFloat / 1024
	fmt.Println(res)
	return ""
}