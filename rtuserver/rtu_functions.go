package rtuserver

import "encoding/binary"

const (
	F_REGISTER  = 0x46
	F_HEARTBEAT = 0x47
)

/**
 处理DTU 的注册报文
 Function: 0x46
 */
func Register(rtuServer *RtuServer, c *RtuConn, f *RTUFrame) (data []byte, exception Exception) {
	if f.Length != 2 + 13 {
		return nil, IllegalDataValue
	}

	area := binary.BigEndian.Uint16(f.Data[0:2])
	deviceId := f.Data[2]
	deviceName := f.Data[3:9]
	remark := f.Data[9:13]

	rtuServer.Register(c, area, deviceId, deviceName, remark)
	return f.Data, Success
}

/**
 处理DTU 的心跳报文
 Function: 0x47
 */
func HeartBeat(rtuServer *RtuServer, c *RtuConn, f *RTUFrame) (data []byte, exception Exception) {
	if f.Length != 2 + 3 {
		return nil, IllegalDataValue
	}

	area := binary.BigEndian.Uint16(f.Data[0:2])
	deviceId := f.Data[2]

	rtuServer.HeartBeat(c, area, deviceId)
	return f.Data, Success
}