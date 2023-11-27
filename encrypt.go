package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// EncryptedConn 是自定义的加密连接接口，类似于 net.Conn
type EncryptedConn interface {
	io.ReadWriteCloser
}
type aes_gcm_Conn struct {
	key     string
	srcConn net.Conn
}

func NewAesConn(key string, conn net.Conn) (EncryptedConn, error) {
	a := aes_gcm_Conn{
		key:     key,
		srcConn: conn,
	}
	return &a, nil
}
func (agc *aes_gcm_Conn) Read(b []byte) (n int, err error) {

	var bodySize = make([]byte, 4)
	io.ReadFull(agc.srcConn, bodySize)
	size := binary.BigEndian.Uint32(bodySize)
	fmt.Println("size is ", size)
	buf := make([]byte, size)
	io.ReadFull(agc.srcConn, buf)
	fmt.Println(buf[:])
	plaintext, err := decrypt(agc.key, buf[:])
	fmt.Println("plaintext is ", plaintext)
	if err != nil {
		return 0, err
	}
	n = copy(b, plaintext)
	return n, nil
}

func (agc *aes_gcm_Conn) Write(b []byte) (n int, err error) {

	ciphertext, err := encrypt(agc.key, b)
	fmt.Println("ciphertext is ", ciphertext)
	if err != nil {
		fmt.Println("encrypt error")
		return 0, err
	}
	bodySize := uint32(len(ciphertext))
	data := make([]byte, 4+bodySize)
	binary.BigEndian.PutUint32(data[:4], bodySize)
	copy(data[4:], ciphertext)
	n, err = agc.srcConn.Write(data)
	return n, err

}
func (agc *aes_gcm_Conn) Close() error {
	agc.srcConn.Close()
	return nil
}

func makeKey(key string) []byte {
	k := md5.Sum([]byte(key))
	return k[:16]
}

func encrypt(Key string, plaintext []byte) ([]byte, error) {
	key := makeKey(Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(nonce, ciphertext...) // 将nonce与密文拼接在一起

	return ciphertext, nil
}

func decrypt(Key string, ciphertext []byte) ([]byte, error) {
	key := makeKey(Key)
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("密文长度错误")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
