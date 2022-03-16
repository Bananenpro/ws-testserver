package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Bananenpro/log"
)

var waitingForInput = false
var inputQuestion = ""

func Input(question string) string {
	fmt.Print(question)
	scanner := bufio.NewScanner(os.Stdin)
	waitingForInput = true
	inputQuestion = question
	scanner.Scan()
	waitingForInput = false
	return scanner.Text()
}

func PrintMessage(msg string) {
	fmt.Print("\x1b[2K\r")
	log.Info("Received:", msg)
	if waitingForInput {
		fmt.Print(inputQuestion)
	}
}

func PrintError(msg string) {
	fmt.Println()
	log.Error("ERROR:", msg)
	if waitingForInput {
		fmt.Print(inputQuestion)
	}
}
