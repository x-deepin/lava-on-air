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
	TypeFPanelControl       = 3
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

const MasterNode uint16 = 0

type DevCreator func(node uint16, Type uint8) Device

func Register(drv *Driver, creator DevCreator, Type uint8) {
}

type Device interface {
	Read()
	Write([]byte)
}

type Driver struct {
	rfNode  uint16
	channel io.ReadWriter
	nextId  uint16
	devices map[uint16]*PTS
}

func NewDriver(port string, bauds int, rfNode uint16) (*Driver, error) {
	c := &serial.Config{Name: port, Baud: bauds}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return &Driver{
		rfNode:  rfNode,
		channel: s,
		nextId:  0,
	}, nil
}

func (drv *Driver) write(h Header, payload []byte) error {
	if h.PayloadSize != uint8(len(payload)) {
		return fmt.Errorf("Invalid header")
	}

	binary.Write(drv.channel, binary.BigEndian, h)
	binary.Write(drv.channel, binary.BigEndian, payload)
	return nil
}

func (drv *Driver) NextId() uint16 {
	drv.nextId++
	return drv.nextId
}

func (drv *Driver) SendFPanelMessage(to uint16) error {
	h := Header{
		From:        drv.rfNode,
		To:          to,
		Id:          drv.nextId,
		Type:        TypeFPanelControl,
		PayloadSize: 1,
	}
	drv.nextId++
	return drv.write(h, []byte{3})
}

func (drv *Driver) Show(h Header, payload []byte) {
	switch h.Type {
	case TypeSerialMessage:
		fmt.Printf("Serial Message[%d] (%d --> %d Len:%d):\n\t %q\n", h.Type, h.From, h.To, h.PayloadSize, payload)
	case TypeFPanelControl:
		fmt.Printf("Control Message[%d](%d --> %d Len:%d):\n\t %v\n", h.Type, h.From, h.To, h.PayloadSize, payload)
	case TypeStatus:
		fmt.Printf("Status Message[%d](%d --> %d Len:%d):\n\t %v\n", h.Type, h.From, h.To, h.PayloadSize, payload)
	case TypeHello:
		fmt.Printf("Hello Message[%d] (%d --> %d Len:%d):\n\t %q\n", h.From, h.To, h.PayloadSize, payload)
	case TypeInvalid:
	default:
		fmt.Printf("Unknown Message(%d --> %d Type:%d Len:?)\n", h.From, h.To, h.Type)
	}
}

func (drv *Driver) Read() {
	h := Header{}
	binary.Read(drv.channel, binary.LittleEndian, &h)

	buf := make([]byte, h.PayloadSize)
	if h.PayloadSize > 0 {
		binary.Read(drv.channel, binary.LittleEndian, buf)
	}

	// Tee to devices
	drv.Show(h, buf)

	pts.Write(buf)
}

func main() {
	d, err := NewDriver("/dev/ttyUSB0", 57600, MasterNode)
	if err != nil {
		log.Println(err)
	}
	print("Starting...\n")
	for {
		d.Read()
	}
}
