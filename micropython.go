package micropythongo

import (
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

// A MicroPython device connected over serial.
type MicroPyDevice struct {
	Port       *serial.Port // Pointer to a serial port
	MPYVersion string       // MicroPython Version
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

	// Remove the square brackets
	trimmed := strings.Trim(response, "[]")

	// Split by comma
	parts := strings.Split(trimmed, ",")

	// Clean up each element
	var files []string
	for _, part := range parts {
		// Remove quotes and whitespace
		cleaned := strings.Trim(part, " '\"")
		if cleaned != "" {
			files = append(files, cleaned)
		}
	}

	return files, nil
}
