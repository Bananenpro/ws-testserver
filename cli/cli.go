package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Bananenpro/log"
)

var waitingForInput = false
var currentInputPrompt = ""

func AskForMessage() string {
	waitingForInput = true
	currentInputPrompt = "Press enter to send a message..."
	printCurrentInputPrompt()
	fmt.Scanln()

	fmt.Print("\x1b[1A\x1b[2K\r")
	currentInputPrompt = "Enter the file extension: "
	printCurrentInputPrompt()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	waitingForInput = false

	fmt.Print("\x1b[1A\x1b[2K\r")

	ext := strings.TrimPrefix(strings.TrimSpace(scanner.Text()), ".")
	if ext == "" {
		ext = "txt"
	}

	return inputMessage(ext)
}

func PrintMessage(msg string) {
	fmt.Print("\x1b[2K\r")
	log.Info("Received:", msg)
	if waitingForInput {
		printCurrentInputPrompt()
	}
}

func PrintError(msg string) {
	fmt.Println()
	log.Error("ERROR:", msg)
	if waitingForInput {
		printCurrentInputPrompt()
	}
}

func printCurrentInputPrompt() {
	fmt.Print(currentInputPrompt)
}

func inputMessage(fileExtension string) string {
	editor := getDefaultEditorName()
	tempFile, err := os.CreateTemp("", fmt.Sprintf("input.*.%s", fileExtension))
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create temp file: %s", err))
	}
	tempFile.Close()

	cmd := exec.Command(editor, tempFile.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run default editor '%s': %s", editor, err))
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		PrintError(fmt.Sprintf("Failed to read temp file: %s", err))
	}

	content := string(data)
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, " ", "")

	os.Remove(tempFile.Name())

	return content
}

func getDefaultEditorName() string {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("code"); err == nil {
			return "code"
		}

		return "notepad"
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	return "vim"
}
