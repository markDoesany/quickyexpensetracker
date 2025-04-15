package services

import "fmt"

func ProcessCommand(command, psid, mid, token string) {
	fmt.Printf("Processing Command: %s, PSID: %s, MID: %s\n", command, psid, mid)
	fmt.Printf("token: %s\n", token)
}
