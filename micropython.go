package micropythongo

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.bug.st/serial"
)

// A MicroPython device connected over serial.
type MicroPyDevice struct {
	Port       *serial.Port // Pointer to a serial port
	MPYVersion string       // MicroPython Version
}

// Parse a list or tuple from Python
func parsePythonList(input string) []string {
	// Handle both list and tuple syntax
	input = strings.TrimSpace(input)

	// Remove the outer brackets/parentheses
	if (strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]")) ||
		(strings.HasPrefix(input, "(") && strings.HasSuffix(input, ")")) {
		input = input[1 : len(input)-1]
	}

	// Handle empty list/tuple
	if len(input) == 0 {
		return []string{}
	}

	// Extract quoted strings
	re := regexp.MustCompile(`'[^']*'|"[^"]*"`)
	matches := re.FindAllString(input, -1)

	result := make([]string, 0, len(matches))
	for _, match := range matches {
		// Remove the quotes
		cleaned := match[1 : len(match)-1]
		result = append(result, cleaned)
	}

	return result
}

// Connect to a MicroPython device using serial. Returns with a MicroPyDevice.
func ConnectMPYDevice(devicePort serial.Port) (*MicroPyDevice, error) {
	devicePort.Write([]byte("import os, machine\r\n"))
	time.Sleep(100 * time.Millisecond)
	devicePort.Write([]byte("os.uname()\r\n"))

	time.Sleep(100 * time.Millisecond)

	buf := make([]byte, 256)
	n, err := devicePort.Read(buf)
	if err != nil {
		return nil, err
	}
	response := string(buf[:n])

	if strings.Contains(response, "version") {
		newDevice := MicroPyDevice{
			Port:       &devicePort,
			MPYVersion: "placeholder",
		}

		return &newDevice, err
	} else {
		return nil, fmt.Errorf("device is not a valid MicroPython device")
	}
}

// Instruct the MicroPython device to enter bootloader mode, if supported.
func (device MicroPyDevice) EnterBootloader() error {
	_, err := (*device.Port).Write([]byte("machine.bootloader()\r\n"))

	if err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	return nil
}

// Save a file to the MicroPython device's filesystem.
func (device MicroPyDevice) SaveFile(fileName string, contents string) error {
	command := "f = open('" + fileName + "', 'w')\r\n"

	_, err := (*device.Port).Write([]byte(command))

	if err != nil {
		return err
	}

	return nil
}

// List files in the MicroPython device's filesystem.
func (device MicroPyDevice) ListFiles(directory string) ([]string, error) {
	command := "os.listdir('" + directory + "')\r\n"
	_, err := (*device.Port).Write([]byte(command))

	if err != nil {
		return nil, err
	}

	time.Sleep(50 * time.Millisecond)

	// Create a buffer and read into it
	buf := make([]byte, 256)
	n, err := (*device.Port).Read(buf)

	// Chceck for errors
	if err != nil {
		return nil, err
	}

	// Convert the buffer
	response := string(buf[:n])

	files := parsePythonList(response)

	return files, nil
}
