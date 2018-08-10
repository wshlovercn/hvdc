package main

import (
	"net"
	"hvdc/rtuserver"
	"encoding/binary"
	"fmt"
)

/**
 Test connection register
 */

const addr = "127.0.0.1:502"

func main() {
	testRegister();

	testUnregister();
}

func testRegister() {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial error: ", err)
		return
	}

	defer c.Close()

	sendRegister(c)

	sendHeartBeat(c)

	sendUnknow(c)

	//wait time out
	packet := make([]byte, 512)
	n, err := c.Read(packet)
	fmt.Printf("testRegister read %d, %v\n", n, err)
}

func testUnregister()  {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial error: ", err)
		return
	}

	defer c.Close()

	sendHeartBeat(c)

	//wait time out
	packet := make([]byte, 512)
	n, err := c.Read(packet)
	fmt.Printf("testUnregister read %d, %v\n", n, err)
}

func sendRegister(c net.Conn) {
	//send
	f := rtuserver.RTUFrame{
		TransactionIdentifier:1,
		ProtocolIdentifier:2,
		Device:1,
		Function:rtuserver.F_REGISTER,
	}

	data := make([]byte, 13)
	binary.BigEndian.PutUint16(data[0:2], 13)
	data[2] = 3
	copy(data[3:9], []byte("HELLO1"))
	copy(data[9:13], []byte("OOOO"))
	f.SetData(data)

	c.Write(f.Bytes())

	//read
	packet := make([]byte, 512)
	if n, err := c.Read(packet); err != nil {
		fmt.Printf("sendHeartBeat read error : %v\n", err)
	} else {
		if r, err := rtuserver.NewRTUFrame(packet[:n]); err != nil {
			fmt.Printf("sendRegister error response: %v\n", err)
		} else {
			fmt.Printf("sendRegister response: %v\n", r)
		}
	}
}

func sendHeartBeat(c net.Conn) {
	//send
	f := rtuserver.RTUFrame{
		TransactionIdentifier:1,
		ProtocolIdentifier:2,
		Device:1,
		Function:rtuserver.F_HEARTBEAT,
	}

	data := make([]byte, 3)
	binary.BigEndian.PutUint16(data[0:2], 13)
	data[2] = 3
	f.SetData(data)

	c.Write(f.Bytes())

	//read
	packet := make([]byte, 512)
	if n, err := c.Read(packet); err != nil {
		fmt.Printf("sendHeartBeat read error : %v\n", err)
	} else {
		if r, err := rtuserver.NewRTUFrame(packet[:n]); err != nil {
			fmt.Printf("sendHeartBeat error response: %v\n", err)
		} else {
			fmt.Printf("sendHeartBeat response: %v\n", r)
		}
	}
}

func sendUnknow(c net.Conn) {
	//send
	f := rtuserver.RTUFrame{
		TransactionIdentifier:1,
		ProtocolIdentifier:2,
		Device:1,
		Function:1,
	}

	data := make([]byte, 3)
	binary.BigEndian.PutUint16(data[0:2], 13)
	data[2] = 3
	f.SetData(data)

	c.Write(f.Bytes())

	//read
	packet := make([]byte, 512)
	if n, err := c.Read(packet); err != nil {
		fmt.Printf("sendUnknow read error : %v\n", err)
	} else {
		if r, err := rtuserver.NewRTUFrame(packet[:n]); err != nil {
			fmt.Printf("sendUnknow: error response: %v\n", err)
		} else {
			fmt.Printf("sendUnknow response: %v\n", r)
		}
	}
}