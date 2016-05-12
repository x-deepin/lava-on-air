package main

import (
	"encoding/binary"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
)

const (
	TypeSerialMessage uint8 = 2
	TypeControl             = 3
	TypeStatus              = 4
	TypeInvalid             = 5
	TypeHello               = 6
	TypeMaxNum              = 7
)

type Header struct {
	From        uint16
	To          uint16
	Id          uint16
	Type        uint8
	PayloadSize uint8
}

func read(r io.Reader) {
	h := Header{}
	binary.Read(r, binary.LittleEndian, &h)

	buf := make([]byte, h.PayloadSize)
	if h.PayloadSize > 0 {
		binary.Read(r, binary.LittleEndian, buf)
	}

	switch h.Type {
	case TypeSerialMessage:
		fmt.Printf("Serial Message[%d] (%d --> %d Len:%d):\n\t %q\n", h.Type, h.From, h.To, h.PayloadSize, buf)
	case TypeControl:
		fmt.Printf("Control Message[%d](%d --> %d Len:%d):\n\t %v\n", h.Type, h.From, h.To, h.PayloadSize, buf)
	case TypeStatus:
		panic("Not implement firmware, so shouldn't receive this type message")
	case TypeHello:
		fmt.Printf("Hello Message[%d] (%d --> %d Len:%d):\n\t %q\n", h.From, h.To, h.PayloadSize, buf)
	case TypeInvalid:
	default:
		fmt.Printf("Unknown Message(%d --> %d Type:%d Len:?)\n", h.From, h.To, h.Type)
	}
}

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	print("Starting...\n")
	for {
		read(s)
	}
}
