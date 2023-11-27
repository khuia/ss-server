package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/armon/go-socks5"
)

func socks5Start(c *Config) {

	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port
	if err := server.ListenAndServe("tcp", c.Socks5Addr); err != nil {
		panic(err)
	}

}

func socks5Begin() {

	fmt.Println("hello world")
	l, _ := net.Listen("tcp", "127.0.0.1:8080")
	for {
		conn, _ := l.Accept()
		go handle(conn)

	}

}

func handle(conn net.Conn) {
	defer conn.Close()
	auth(conn)
	addr := connect(conn)
	fmt.Println("addr is ", addr)
	targetConn, err := net.Dial("tcp", addr)

	if err != nil {
		fmt.Println("connect target error :", err)
		return
	}
	targetConn.Close()
	_, err = conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})
	if err != nil {
		fmt.Println("write err", err)
		return
	}
	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)

}
func auth(conn net.Conn) {
	buf := make([]byte, 128)
	conn.Read(buf)
	fmt.Println(buf)
	_, err := conn.Write([]byte{5, 0})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}
func connect(conn net.Conn) string {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	buf := make([]byte, 4)
	io.ReadFull(conn, buf[:4])
	fmt.Println(buf[:4])
	addr := ""
	switch buf[3] {
	case 1:
		io.ReadFull(conn, buf)
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
		fmt.Println("ipv4 is ", addr)
	case 3:
		sizeArea := make([]byte, 1)
		io.ReadFull(conn, sizeArea[:])
		fmt.Println("size is ", sizeArea)
		host := make([]byte, sizeArea[0])
		io.ReadFull(conn, host[:])
		fmt.Println("host is ", string(host))
		addr = string(host)

	default:

	}
	io.ReadFull(conn, buf[:2])
	port := binary.BigEndian.Uint16(buf[:2])
	fmt.Printf("addr is %s,port is %d", addr, port)
	return fmt.Sprintf("%s:%d", addr, port)

}
