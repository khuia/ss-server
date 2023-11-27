package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
)

func server_aes_advPortForward(c *Config) {
	// 监听本地地址
	localAddr := c.LocalAddr
	remoteAddr := c.RemoteAddr
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatal("监听本地地址失败:", err)
	}
	defer listener.Close()
	log.Printf("开始监听本地地址：%s\n", localAddr)

	for {
		// 等待客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Fatal("接受客户端连接失败:", err)
		}
		log.Printf("接受客户端连接：%s\n", clientConn.RemoteAddr())

		// 连接目标服务器
		serverConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			log.Fatal("连接目标服务器失败:", err)
		}
		log.Printf("连接目标服务器：%s\n", remoteAddr)

		// 启动Go协程将客户端数据转发到目标服务器
		go func() {
			buf := make([]byte, 4070)
			n, err := clientConn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatal("error", err)
				}

			}

			ciphertext, err := decrypt(c.Key, buf[:n])
			if err != nil {
				fmt.Println("encrypt err:", err)
			}
			serverConn.Write(ciphertext[:])
			//_, err := io.Copy(serverConn, clientConn)
			/*if err != nil {
				log.Printf("从客户端到目标服务器转发数据失败：%s\n", err)
			}*/
		}()

		// 启动Go协程将目标服务器数据转发到客户端
		go func() {

			buf := make([]byte, 4070)
			n, err := serverConn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatal("serverRead err:", err)
				}
			}
			detxt, err := encrypt(c.Key, buf[:n])
			if err != nil {
				if err != io.EOF {
					log.Fatal("decrypt err:", err)
				}
			}
			clientConn.Write(detxt[:])
			/*_, err := io.Copy(clientConn, serverConn)
			if err != nil {
				log.Printf("从目标服务器到客户端转发数据失败：%s\n", err)
			}*/
		}()
	}
}

func server_advPortForward(c *Config) {
	// 监听本地地址
	localAddr := c.LocalAddr
	remoteAddr := c.RemoteAddr
	/*listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatal("监听本地地址失败:", err)
	}
	defer listener.Close()
	log.Printf("开始监听本地地址：%s\n", localAddr)
	*/
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", localAddr, config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server is listening on :", localAddr)
	for {
		// 等待客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Fatal("接受客户端连接失败:", err)
		}
		log.Printf("接受客户端连接：%s\n", clientConn.RemoteAddr())

		// 连接目标服务器
		serverConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			log.Fatal("连接目标服务器失败:", err)
		}
		log.Printf("连接目标服务器：%s\n", remoteAddr)

		// 启动Go协程将客户端数据转发到目标服务器
		go func() {
			_, err := io.Copy(serverConn, clientConn)
			if err != nil {
				log.Printf("从客户端到目标服务器转发数据失败：%s\n", err)
			}
		}()

		// 启动Go协程将目标服务器数据转发到客户端
		go func() {
			_, err := io.Copy(clientConn, serverConn)
			if err != nil {
				log.Printf("从目标服务器到客户端转发数据失败：%s\n", err)
			}
		}()
	}
}

func PortForward(c *Config) {
	// 监听本地地址
	localAddr := c.LocalAddr
	remoteAddr := c.RemoteAddr
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatal("监听本地地址失败:", err)
	}
	defer listener.Close()
	log.Printf("开始监听本地地址：%s\n", localAddr)

	for {
		// 等待客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Fatal("接受客户端连接失败:", err)
		}
		log.Printf("接受客户端连接：%s\n", clientConn.RemoteAddr())

		// 连接目标服务器
		serverConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			log.Fatal("连接目标服务器失败:", err)
		}
		log.Printf("连接目标服务器：%s\n", remoteAddr)

		// 启动Go协程将客户端数据转发到目标服务器
		go func() {
			_, err := io.Copy(serverConn, clientConn)
			if err != nil {
				log.Printf("从客户端到目标服务器转发数据失败：%s\n", err)
			}
		}()

		// 启动Go协程将目标服务器数据转发到客户端
		go func() {
			_, err := io.Copy(clientConn, serverConn)
			if err != nil {
				log.Printf("从目标服务器到客户端转发数据失败：%s\n", err)
			}
		}()
	}
}

// 普通转发
func Forward(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	fmt.Printf("transferring data from %s to %s\n", src.(net.Conn).RemoteAddr(), dst.(net.Conn).RemoteAddr())

	srcData, err := io.ReadAll(src)
	if err != nil {
		fmt.Println("failed to read data:", err)
		return
	}

	if _, err := dst.Write(srcData); err != nil {
		fmt.Println("failed to write data:", err)
		return
	}
}

// 进阶转发
func advForward(des net.Conn, src net.Conn) {
	defer des.Close()
	defer src.Close()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		fmt.Println("4096读取的数据为", buf[:n])
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		if _, err := des.Write(buf[:n]); err != nil {
			log.Println(err)
			break
		}
	}
}
