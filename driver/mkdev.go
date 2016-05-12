package main

type PTS struct {
	drv  *Driver
	node uint16
	buf  []byte
}

func (pts *PTS) Write(data []byte) error {
	h := Header{
		From:        pts.drv.rfNode,
		To:          pts.node,
		Id:          pts.drv.NextId(),
		Type:        TypeSerialMessage,
		PayloadSize: uint8(len(data)),
	}
	return pts.drv.write(h, data)
}

func (pts *PTS) Read() error {
	return nil
}

func (pts *PTS) Path() string {
	return ""
}

func NewDevice(node uint16, basePath string) *PTS {
	return nil
}
