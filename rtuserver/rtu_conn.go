package rtuserver

import (
	"github.com/golang/glog"
	"sync"
	"sync/atomic"
	"hvdc/baselib/utils"
	"net"
	"fmt"
	"time"
)

type RtuConnCallback interface {
	OnConnected(c *RtuConn)
	OnDisconnected(c *RtuConn, err error)
	OnRtuFrame(c *RtuConn, f *RTUFrame) error
}

type RtuConn struct {
	conn     net.Conn
	callback RtuConnCallback

	once     sync.Once
	closeCh  chan struct{}
	writeCh  chan *RTUFrame

	err      atomic.Value
	running  utils.AtomicBoolean
	timeout  time.Duration

	registerd  utils.AtomicBoolean
	area       uint16
	deviceId   uint8
	deviceName []byte
	remark     []byte
}

func NewRtuConn(conn net.Conn, callback RtuConnCallback) *RtuConn {
	rtuConn := &RtuConn{
		conn:conn,
		callback:callback,

		closeCh:make(chan struct{}, 0),
		writeCh:make(chan *RTUFrame, 1),
		timeout:5,
	}
	rtuConn.SetRunning(false)
	rtuConn.SetErr(nil)

	return rtuConn
}

func (c *RtuConn) IsRunning() bool {
	return c.running.Get()
}

func (c *RtuConn) SetRunning(r bool)  {
	c.running.Set(r)
}

func (c *RtuConn) SetErr(err error)  {
	if err != nil {
		c.err.Store(err)
	}
}

func (c *RtuConn) GetErr() error {
	err := c.err.Load()
	return err.(error)
}

func (c *RtuConn) IsRegistered() bool {
	return c.registerd.Get()
}

func (c *RtuConn) SetRegistered() {
	c.registerd.Set(true)
}

func (c *RtuConn) Register(area uint16, deviceId uint8, deviceName []byte, remark []byte) {
	if c.IsRegistered() {
		return
	}

	c.SetRegistered()
	c.area = area
	c.deviceId = deviceId
	c.deviceName = deviceName
	c.remark = remark
}

func (c *RtuConn) Key() string {
	if c.IsRegistered() {
		return fmt.Sprintf("%v-%v", c.area, c.deviceId)
	} else {
		return ""
	}
}

func (c *RtuConn) Start()  {
	go c.readLoop()
	go c.writeLoop()

	c.callback.OnConnected(c)
}

func (c *RtuConn) Read() (*RTUFrame, error) {
	packet := make([]byte, 512)

	c.conn.SetReadDeadline(time.Now().Add(4 * c.timeout * time.Second))
	bytesRead, err := c.conn.Read(packet)
	if err != nil {
		glog.Errorf("read error %v", err)
		return nil, err
	}
	packet = packet[:bytesRead]

	frame, err := NewRTUFrame(packet)
	if err != nil {
		glog.Errorf("bad packet error %v", err)
		return nil, err
	}
	return frame, nil
}

func (c *RtuConn) Write(f *RTUFrame)  {
	c.writeCh <- f
}

func (c *RtuConn) Close() {
	c.once.Do(func() {
		close(c.closeCh)

		c.conn.Close()

		if c.callback != nil {
			c.callback.OnDisconnected(c, c.GetErr())
		}
	})
}


func (c *RtuConn) readLoop()  {
	defer c.Close()

	for {
		select {
		case <- c.closeCh:
			return

		default:
			if f, err := c.Read(); err != nil {
				glog.Errorf("%s read error %v", c.String(), err)
				c.SetErr(err)
				return
			} else {
				if err := c.callback.OnRtuFrame(c, f); err != nil {
					glog.Errorf("%s OnRtuFrame error %v", c.String(), err)
					c.SetErr(err)
					return
				}
			}
		}
	}
}

func (c *RtuConn) writeLoop() {
	defer c.Close()

	for {
		select {
		case <- c.closeCh:
			return

		case f := <- c.writeCh:
			if _, err := c.conn.Write(f.Bytes()); err != nil {
				glog.Errorf("%s write error %v", err)
				c.SetErr(err)
				return
			}
		}
	}
}

func (c *RtuConn) String() string {
	s := fmt.Sprintf("RtuConn{key:%s, addr:%s}", c.Key(), c.conn.RemoteAddr().String())
	return s
}