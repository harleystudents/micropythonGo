package micropythongo

import (
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

type MicroPyDevice struct {
	Port       *serial.Port // Pointer to a serial port
	MPYVersion string       // MicroPython Version
}

func connectMPYDevice(devicePort serial.Port) (*MicroPyDevice, error) {
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
