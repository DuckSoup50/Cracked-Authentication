package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	auth()
}

func continuation() {
	fmt.Print("Hello World!")
}

func auth() {

	usergroups := []string{
		"3",   // Moderator
		"4",   // Administrator
		"6",   // Trial Moderator
		"11",  // Premium
		"12",  // Supreme
		"91",  // Retired Staff
		"92",  // Dreams
		"93",  // Infinity
		"94",  // Coder
		"95",  // SE-God
		"97",  // Reverser
		"98",  // Disinfector
		"99",  // Boss of bosses
		"100", // Irefunder
		"101", // Section Moderator
	}

	// Hwid = Hashed Processor ID + System UUID

	uuidCmd := exec.Command("wmic", "csproduct", "get", "uuid")
	uuidOut, err := uuidCmd.Output()

	if err != nil {
		panic(err)
	}

	systemId := strings.Split(string(uuidOut), "\n")[1]

	pidCmd := exec.Command("wmic", "cpu", "get", "ProcessorId")
	pidOut, err := pidCmd.Output()

	if err != nil {
		panic(err)
	}

	procId := strings.Split(string(pidOut), "\n")[1]

	h := sha256.New()
	h.Write([]byte(procId + systemId))
	hwid := fmt.Sprintf("%x", h.Sum((nil)))

	authFile := "Key.cio"
	var Key string

	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		fmt.Print("Enter your auth key: ")
		fmt.Scan(&Key)

		err := ioutil.WriteFile(authFile, []byte(Key), 0644)

		if err != nil {
			panic(err)
		}
	} else {
		inputKey, err := ioutil.ReadFile(authFile)

		if err != nil {
			panic(err)
		}

		Key = strings.TrimSpace(string(inputKey))
	}

	payload := fmt.Sprintf("a=auth&k=%s&hwid=%s", Key, hwid)

	resp, err := http.Post("https://cracked.io/auth.php", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(payload)))

	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println()

	if resp.StatusCode == http.StatusOK {
		var res map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if errMsg, ok := res["error"].(string); ok {
			fmt.Println(errMsg)
		}

		if auth, ok := res["auth"].(bool); ok && auth == true {
			// Check if user has permissions
			username := res["username"].(string)
			group := res["group"].(string)

			hasPerms := false

			for _, ug := range usergroups {
				if ug == group {
					fmt.Println("Successfully authenticated! Welcome back: " + username + "!")
					hasPerms = true
					// Code continues here -> call whatever you what
					break
				}
			}

			if !hasPerms {
				fmt.Println("You have to be at least Premium+ to be able to use this tool!")
				os.Exit(1)
			}
		} else {
			// fmt.Printf("Authentication failed!")
		}

	}

}
