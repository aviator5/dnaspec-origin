package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks the user a yes/no question and returns true if they answer yes.
func Confirm(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s (y/N): ", question)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
