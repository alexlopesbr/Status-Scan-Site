package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const delay = 2
const numberOfCalls = 2

func main() {
	showIntro()
	for {
		showMenu()
		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Printf("Showing Logs...")
			printLogs()
		case 0:
			fmt.Printf("Exiting...")
			os.Exit(0)
		default:
			fmt.Printf("Comand invalid")
			os.Exit(-1)
		}
	}
}

func showIntro() {
	name := "Alexandre"
	fmt.Println("Olá", name)
}

func showMenu() {
	fmt.Println("1- Init monitoring")
	fmt.Println("2- Show Logs")
	fmt.Println("0- Exit program")
}

func readCommand() int {
	var commandRead int
	fmt.Scan(&commandRead)
	fmt.Println("The command is:", commandRead)

	return commandRead
}

func startMonitoring() {
	fmt.Printf("Monitoring...\n")
	websites := folderSitesList()

	for i := 0; i < numberOfCalls; i++ {
		for _, website := range websites {
			testSite(website)
		}
		time.Sleep(delay * time.Second)
	}
}

func testSite(website string) {
	resp, err := http.Get(website)

	if err != nil {
		fmt.Println("An error ocurred:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Website", website, "OK")
		registerLog(website, true)
	} else {
		fmt.Println("Website", website, "with problens. Status Code:", resp.StatusCode)
		registerLog(website, false)
	}
}

func folderSitesList() []string {
	var websites []string
	folder, err := os.Open("websites.txt")

	if err != nil {
		fmt.Println("An error ocurred:", err)
	}

	reader := bufio.NewReader(folder)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		websites = append(websites, line)

		if err == io.EOF {
			break
		}
	}
	folder.Close()

	return websites
}

func registerLog(website string, status bool) {
	folder, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	folder.WriteString(fmt.Sprintf("%s - %s - online %t\n",
		time.Now().Format("02/01/2006 15:04:05"),
		website,
		status,
	))

	folder.Close()
}

func printLogs() {
	folder, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(folder))
}
