package main

import (
	"fmt"
	"io"

	"net"
)

func main() {

	l, _ := net.Listen("tcp", "127.0.0.1:8081")

	for {
		conn, _ := l.Accept()

		go hand(conn)
	}

}
func hand(conn net.Conn) {
	srcConn, _ := NewAesConn("hello", conn)
	defer srcConn.Close()

	sockType := make([]byte, 1)
	io.ReadFull(srcConn, sockType)
	addr := ""
	switch sockType[0] {
	case 1:
		data := make([]byte, 11)
		io.ReadFull(srcConn, data)
		addr = string(data)
	case 3:
		sizeArea := make([]byte, 1)
		io.ReadFull(srcConn, sizeArea)
		host := make([]byte, sizeArea[0]+3)
		io.ReadFull(srcConn, host)
		addr = string(host)

	}
	fmt.Println("addr is ", addr)
	destConn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("connect dastConn error", err)
		srcConn.Write([]byte{1})
		return
	}
	srcConn.Write([]byte{0})
	go io.Copy(destConn, srcConn)
	io.Copy(srcConn, destConn)

}
