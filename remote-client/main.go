package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// Server URL
	serverURL       string
	programFilename string
)

var (
	currentProgramID string
	reloadProgram    chan bool
	results          chan string
)

func init() {
	flag.StringVar(&serverURL, "server", "", "Server URL")
	flag.StringVar(&programFilename, "program", "/tmp/remote-code", "Program filename")
}

// Get the current program ID
func getProgramID() (string, error) {
	res, err := http.Get(serverURL + "/program/id")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getProgram() ([]byte, error) {
	res, err := http.Get(serverURL + "/program")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func publishResult(result string) error {
	_, err := http.Post(serverURL+"/results", "application/json", bytes.NewBuffer([]byte(result)))
	if err != nil {
		return err
	}
	return nil
}

func fetchProgram(filename string) (bool, error) {
	// Get the current programID and only update
	// if it has changed
	programID, err := getProgramID()
	if err != nil {
		return false, err
	}
	if programID == currentProgramID {
		return false, nil
	}
	currentProgramID = programID

	// Get the program and safe to file
	text, err := getProgram()
	if err != nil {
		return false, err
	}

	// Write the program to file. DO NOT USE THIS IN PRODUCTION.
	if os.WriteFile(filename, text, 0644) != nil {
		return false, err
	}
	return true, nil
}

// Read the stdout of the programm and forward
// results to the channel.
func readProgramOutput(out io.ReadCloser) {
	reader := bufio.NewReader(out)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		result, isResult := strings.CutPrefix(line, "[RESULT] ")
		if isResult {
			results <- result
		}
	}

}

// Start a program and keep it running in the backgound.
// If the program crashes, it will be restarted.
func spawnProgram(filename string, stop chan bool) {
	for {
		log.Println("starting program")
		cmd := exec.Command("../code-runner/code-runner", "-script", filename, "-out", "/tmp/code-runner.out", "-stdout")
		cmd.Stderr = os.Stderr

		// Pipe stdout and handle output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		go readProgramOutput(stdout)

		if err := cmd.Start(); err != nil {
			log.Println("error starting program:", err)
		}

		go cmd.Wait()
		// go cmd.Run()

		// Check if the process is still running or
		// stop if if stop was triggered
		for {
			select {
			case <-stop:
				log.Println("stopping process")
				if cmd.Process != nil {
					log.Println("killing process")
					cmd.Process.Kill()
				}
				return // We are done here.
			default:
				if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
					break
				}
				if cmd.Process == nil {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// Run the program
func runProgram(filename string, reload chan bool) {
	stop := make(chan bool)
	running := false
	for {
		<-reload
		if running {
			stop <- true
		}
		go spawnProgram(filename, stop)
		running = true
	}
}

func handleResults() {
	for result := range results {
		publishResult(result)
	}
}

func main() {
	flag.Parse()
	if serverURL == "" {
		flag.Usage()
		return
	}

	// The channel is used to trigger a reload of the
	// program. This means stopping the current process and
	// starting the runner again.
	reloadProgram = make(chan bool)
	results = make(chan string)

	go runProgram(programFilename, reloadProgram)

	go handleResults()

	// Main loop
	for {
		changed, err := fetchProgram(programFilename)
		if err != nil {
			log.Println("error fetching program:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if changed {
			log.Println("program changed, reloading")
			reloadProgram <- true
		}

		time.Sleep(1 * time.Second)
	}
}
