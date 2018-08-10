package rtuserver

import (
	"net"
	"github.com/golang/glog"
	"fmt"
)

type RtuConfig struct {
	Addr string
} 

type RtuServer struct {
	config *RtuConfig
	listener net.Listener
}

func NewRtuServer(config *RtuConfig) *RtuServer {
	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		glog.Fatalf("%v", err)
		return nil
	}

	return &RtuServer{
		config:config,
		listener:l,
	}
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

}

func (rtuServer *RtuServer) OnDisconnected(c *RtuConn, err error) {

}

func (rtuServer *RtuServer) OnRtuFrame(c *RtuConn, f *RTUFrame) error  {
	return nil
}
