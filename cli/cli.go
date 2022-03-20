package cli

import (
	"bufio"
	"encoding/json"
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

	ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(scanner.Text()), "."))
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
	tempFile, err := os.CreateTemp("", fmt.Sprintf("input.*.%s", fileExtension))
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create temp file: %s", err))
	}
	tempFile.Close()

	editor, err := executeDefaultEditor(tempFile.Name())
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run default editor '%s': %s", editor, err))
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		PrintError(fmt.Sprintf("Failed to read temp file: %s", err))
	}

	if fileExtension == "json" {
		var jsonData json.RawMessage
		json.Unmarshal(data, &jsonData)
		encodedData, err := json.Marshal(jsonData)
		if err == nil {
			data = encodedData
		}
	}

	content := string(data)
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "")

	os.Remove(tempFile.Name())

	return content
}

// returns the used editor
func executeDefaultEditor(path string) (string, error) {
	var cmd *exec.Cmd
	var editor string

	if runtime.GOOS == "windows" {
		editor = "notepad"
		if _, err := exec.LookPath("code"); err == nil {
			editor = "code"
		}
		cmd = exec.Command("start", "/wait", editor, path)
		if editor == "code" {
			cmd = exec.Command("code", "--wait", path)
		}
	} else {
		editor = "vi"
		if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		}

		if e := os.Getenv("VISUAL"); e != "" {
			editor = e
		} else if e := os.Getenv("EDITOR"); e != "" {
			editor = e
		}

		cmd = exec.Command(editor, path)
		if editor == "code" {
			cmd = exec.Command("code", "--wait", path)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return editor, cmd.Run()
}
