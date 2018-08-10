package rtuserver

import (
	"net"
	"github.com/golang/glog"
	"fmt"
	"github.com/orcaman/concurrent-map"
)

type RtuConfig struct {
	Addr string
}

type Function func(rtuServer *RtuServer, c *RtuConn, f *RTUFrame) (data []byte, exception Exception)

type RtuServer struct {
	config *RtuConfig

	listener net.Listener
	conns cmap.ConcurrentMap

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
		conns:cmap.New(),
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

	//并不需要做什么
}

func (rtuServer *RtuServer) OnDisconnected(c *RtuConn, err error) {
	glog.Infof("%s OnDisconnected, %v", c.String(), err)

	if c.IsRegistered() {
		rtuServer.conns.Remove(c.Key())

		//TODO:业务处理
	}
}

func (rtuServer *RtuServer) OnRtuFrame(c *RtuConn, f *RTUFrame) error  {
	var exception Exception
	var data []byte

	response := f.Copy()

	f_code := f.GetFunction()
	if !c.IsRegistered() && f_code != F_REGISTER {
		return fmt.Errorf("receive data while unregitered")
	}

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
	glog.Infof("%s register [%v-%v-%s]", c.String(), area, deviceId, string(deviceName))

	c.Register(area, deviceId, deviceName, remark)
	rtuServer.conns.Set(c.Key(), c)

	//TODO:业务处理
}

func (RtuServer *RtuServer) HeartBeat(c *RtuConn, area uint16, deviceId uint8) {
	glog.Infof("%s HeartBeat [%v-%v]", c.String(), area, deviceId)

	//TODO:业务处理
}