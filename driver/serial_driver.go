package main

import (
	"fmt"
	"io"
	"time"
)

type SerialDriver struct {
	baseDir string
	devices map[NodeAddress]io.ReadWriteCloser
	msgCh   chan Message
}

func init() {
	drv := &SerialDriver{
		baseDir: DeviceDirectory,
		devices: make(map[NodeAddress]io.ReadWriteCloser),
	}
	for _, t := range drv.Types() {
		__drivers__[t] = drv
	}
}

func (drv *SerialDriver) InstallDevice(node NodeAddress) error {
	if _, ok := drv.devices[node]; ok {
		return nil
	}

	masterPTS, err := CreatePTS(fmt.Sprintf("%s/%d", drv.baseDir, node))
	if err != nil {
		return err
	}

	drv.devices[node] = masterPTS
	return nil
}

func (drv *SerialDriver) RemoveDevice(node NodeAddress) error {
	if dev, ok := drv.devices[node]; ok {
		dev.Close()
	}
	return nil
}

func (drv *SerialDriver) Start() {
	if drv.msgCh == nil {
		panic("Setup message chan before start")
	}

	var buf [SerialBufferSize]byte
	for {
		for node, dev := range drv.devices {
			<-time.After(time.Millisecond * 20)
			n, err := dev.Read(buf[:])
			if err != nil {
				continue
			}
			drv.msgCh <- Message{
				Type:    TypeSerialMessage,
				Node:    node,
				Payload: buf[:n],
			}
		}
	}
}

func (drv *SerialDriver) SetMessageChan(ch chan Message) {
	drv.msgCh = ch
}

func (drv *SerialDriver) HandleMessage(msg Message) error {
	if msg.Type != TypeSerialMessage {
		panic("It should happen.")
	}

	for node, dev := range drv.devices {
		if node != msg.Node {
			continue
		}
		_, err := dev.Write(msg.Payload)
		return err
	}

	return nil
}

func (drv *SerialDriver) Types() []MessageType { return []MessageType{TypeSerialMessage} }
