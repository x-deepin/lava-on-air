package main

import (
	"bufio"
	"fmt"
	"strings"
)

type FPanelDriver struct {
	*bufio.Reader
	node   uint16
	msgCh  chan Message
	stauts map[uint16]FPanelStatus
}

func init() {
	p := fmt.Sprintf("%s/fpanel", DeviceDirectory)
	md, err := CreatePTS(p)
	if err != nil {
		panic("Failed init fpanel driver")
	}
	drv := &FPanelDriver{
		Reader: bufio.NewReader(md),
	}
	for _, t := range drv.Types() {
		__drivers__[t] = drv
	}
}

func (drv FPanelDriver) Types() []MessageType {
	return []MessageType{TypeFPanelControl, TypeStatus}
}

func (drv *FPanelDriver) HandleMessage(msg Message) error {
	return nil
}
func (drv *FPanelDriver) InstallDevice(node NodeAddress) error {
	return nil
}
func (drv *FPanelDriver) RemoveDevice(node NodeAddress) error {
	return nil
}
func (drv *FPanelDriver) SetMessageChan(ch chan Message) {
	drv.msgCh = ch
}

func (drv *FPanelDriver) Start() {
	// Read from pts/fpanel
	for {
		line, err := drv.ReadString('\n')
		if err != nil {
			continue
		}
		fields := strings.Fields(line)
		if err != nil {
			continue
		}
		if len(fields) != 4 || strings.ToLower(fields[0]) != "set" {
			continue
		}
		//		node, pin, duration := fields[1], fields[2], fields[3]
		node, pin, duration := NodeAddress(0), 0, 0
		drv.msgCh <- Message{
			Type:    TypeFPanelControl,
			Node:    node,
			Payload: []byte{byte(pin), byte(duration)},
		}
	}
}
