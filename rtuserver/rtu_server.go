package rtuserver

import (
	"net"
	"github.com/golang/glog"
	"fmt"
)

type RtuConfig struct {
	Addr string
}

type Function func(rtuServer *RtuServer, c *RtuConn, f *RTUFrame) (data []byte, exception Exception)

type RtuServer struct {
	config *RtuConfig
	
	listener net.Listener
	functions [256]Function
}

func NewRtuServer(config *RtuConfig) *RtuServer {
	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		glog.Fatalf("%v", err)
		return nil
	}

	rtuServer := &RtuServer{
		config:config,
		listener:l,
	}

	rtuServer.functions[F_REGISTER] = Register
	rtuServer.functions[F_HEARTBEAT] = HeartBeat

	return rtuServer
}

func (rtuServer *RtuServer) Serve() error {
	if rtuServer.listener == nil {
		return fmt.Errorf("nil listener")
	}

	for {
		tcpConn, err := rtuServer.listener.Accept()
		if err != nil {
			return err
		}

		rtuConn := NewRtuConn(tcpConn, rtuServer)
		rtuConn.Start()
	}
	return nil
}

func (rtuServer *RtuServer) Stop() {
	if rtuServer.listener != nil {
		rtuServer.listener.Close()
	}
}

/**
  RtuConnCallback
 */

func (rtuServer *RtuServer) OnConnected(c *RtuConn)  {
	glog.Infof("%s OnConnected", c.String())
}

func (rtuServer *RtuServer) OnDisconnected(c *RtuConn, err error) {
	glog.Infof("%s OnDisconnected, %v", c.String(), err)
}

func (rtuServer *RtuServer) OnRtuFrame(c *RtuConn, f *RTUFrame) error  {
	var exception Exception
	var data []byte

	response := f.Copy()

	f_code := f.GetFunction()
	if rtuServer.functions[f_code] != nil {
		data, exception = rtuServer.functions[f_code](rtuServer, c, f)
		response.SetData(data)
	} else {
		exception = IllegalFunction
	}

	if exception != Success {
		response.SetException(exception)
	}

	c.Write(response)

	return nil
}

func (rtuServer *RtuServer) Register(c *RtuConn, area uint16, deviceId uint8, deviceName []byte, remark []byte)  {
	
}

func (RtuServer *RtuServer) HeartBeat(c *RtuConn, area uint16, deviceId uint8) {

}