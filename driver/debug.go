package main

import (
	"fmt"
)

func DebugMessage(in bool, h Header, payload []byte, ok bool) {
	words := map[MessageType]string{
		TypeSerialMessage: "Serial",
		TypeFPanelControl: "FPanelControl",
		TypeStatus:        "Status",
		TypeInvalid:       "Invalid",
		TypeHello:         "Hello",
		TypeEcho:          "echo",
	}

	if !ok && h.Type != TypeInvalid {
		fmt.Printf("No device can handle the %s message %v %q!!!\n", words[h.Type], h, payload)
	}

	if _, ok := words[h.Type]; !ok {
		fmt.Printf("Unknown Message(%d --> %d Type:%d Len:%d)\n", h.From, h.To, h.Type, h.PayloadSize)
		return
	}

	direction := fmt.Sprintf("to node %d from node %d", h.To, h.From)
	if in {
		direction = fmt.Sprintf("from node %d to node %d", h.From, h.To)
	}
	fmt.Printf("%s Message(%d) with %d bytes %s\n",
		words[h.Type],
		h.Id,
		h.PayloadSize,
		direction,
	)

	switch h.Type {
	case TypeSerialMessage, TypeHello:
		fmt.Printf("\t%q\n", payload)
	case TypeInvalid:
		fmt.Printf("\t%t\n", payload)
	case TypeEcho:
		fmt.Printf("\t%t\n", payload[0:8])
		fmt.Printf("\t%q\n", payload[8:])
	case TypeFPanelControl, TypeStatus:
		fmt.Printf("\t%v\n", payload)
	}
}
