package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

type Driver interface {
	SetMessageChan(chan Message)

	HandleMessage(Message) error

	Types() []MessageType

	Start()

	RemoveDevice(NodeAddress) error

	InstallDevice(NodeAddress) error
}

const DeviceDirectory = "/dev/shm/lava-air"

var __drivers__ = make(map[MessageType]Driver)

type Route struct {
	rfNode  NodeAddress
	network io.ReadWriter
	nextId  uint16
	drivers map[MessageType]Driver
	msgChan chan Message
}

func NewRoute(port string, bauds int) (*Route, error) {
	c := &serial.Config{Name: port, Baud: bauds}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	msgChan := make(chan Message)
	for _, drv := range __drivers__ {
		drv.SetMessageChan(msgChan)
		go drv.Start()
	}

	return &Route{
		rfNode:  MasterNode,
		network: s,
		msgChan: msgChan,
		drivers: __drivers__,
	}, nil
}

// HandleMessage Read message from ttyUSB0
func (route *Route) HandleMessage() {
	header := Header{}
	binary.Read(route.network, binary.LittleEndian, &header)

	payload := make([]byte, header.PayloadSize)
	if header.PayloadSize > 0 {
		binary.Read(route.network, binary.LittleEndian, payload)
	}

	msg := Message{
		Type:    header.Type,
		Node:    header.From,
		Payload: payload,
	}

	for t, drv := range route.drivers {
		if t != msg.Type {
			continue
		}
		err := drv.HandleMessage(msg)
		DebugMessage(true, header, payload, true)
		if err == nil {
			break
		}
		fmt.Println("E:", err)
	}
	DebugMessage(true, header, payload, false)
}

// Start write message to ttyUSB0
func (route *Route) Start() {
	go func() {
		for {
			route.HandleMessage()
		}
	}()

	buf := bytes.NewBuffer(nil)
	for msg := range route.msgChan {
		route.nextId++
		h := Header{
			From:        route.rfNode,
			To:          msg.Node,
			Id:          route.nextId,
			Type:        msg.Type,
			PayloadSize: uint8(len(msg.Payload)),
		}

		err := binary.Write(buf, binary.LittleEndian, h)
		if err != nil {
			fmt.Println("Panic at write header", err)
			continue
		}
		err = binary.Write(buf, binary.LittleEndian, msg.Payload)
		if err != nil {
			fmt.Println("Panic at write payload", err)
			continue
		}

		time.Sleep(SerialWait)

		route.network.Write(buf.Bytes())

		buf.Reset()

		DebugMessage(false, h, msg.Payload, true)
	}
}

func (route *Route) AttachNode(node NodeAddress) error {
	for _, drv := range route.drivers {
		drv.InstallDevice(node)
	}
	return nil
}

func main() {
	route, err := NewRoute("/dev/ttyUSB0", 38400)
	if err != nil {
		log.Println(err)
		return
	}

	nodes := []NodeAddress{2}

	for _, node := range nodes {
		route.AttachNode(node)
	}

	route.Start()
}
