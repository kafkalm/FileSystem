package message

import (
	"bytes"
	"encoding/binary"
	"net"
)

func SendMsg(msg []byte,c net.Conn) (error) {
	dataSize := len(msg)
	buf := new(bytes.Buffer)
	err := binary.Write(buf,binary.LittleEndian,int64(dataSize))
	if err != nil {
		return err
	}
	err = binary.Write(buf,binary.LittleEndian,msg)
	if err != nil {
		return err
	}
	_,err = c.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func RecvMsg(c net.Conn) ([]byte,error) {
	// 读取长度
	buf := make([]byte,8)
	_,err := c.Read(buf)
	if err != nil {
		return nil,err
	}

	var dataSize int64
	bufReader := bytes.NewReader(buf)
	err = binary.Read(bufReader,binary.LittleEndian,&dataSize)
	if err != nil {
		return nil,err
	}
	data := make([]byte,dataSize)
	_,err = c.Read(data)
	if err != nil {
		return nil,err
	}
	return data,nil
}