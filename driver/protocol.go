package main

import (
	"math"
	"time"
)

type MessageType uint8

type NodeAddress uint16

type Header struct {
	From        NodeAddress
	To          NodeAddress
	Id          uint16
	Type        MessageType
	PayloadSize uint8
}

type Message struct {
	Node    NodeAddress
	Type    MessageType
	Payload []byte
}

type FPanelCommand struct {
	Pin      uint8
	Duration uint8 // unit 50ms
}
type FPanelStatus struct {
	PowerLED bool
	DiskLED  bool
}

const (
	TypeSerialMessage MessageType = 2
	TypeFPanelControl             = 3
	TypeStatus                    = 4
	TypeInvalid                   = 5
	TypeHello                     = 6
	TypeEcho                      = 7
	TypeMaxNum                    = 8
)

const MasterNode NodeAddress = 0
const MasterSerialBauds = 38400
const SerialBufferSize = 64 - 8

var SerialWait = (time.Duration(math.Floor(float64(MasterSerialBauds)/10.0/float64(SerialBufferSize))) + 10) * time.Millisecond
