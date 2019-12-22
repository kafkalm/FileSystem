package main

import (
	"FileSystem/config"
	. "FileSystem/file"
	"FileSystem/message"
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)
func main() {
	conn,err := net.Dial("tcp",config.Host)
	if err != nil {
		fmt.Println("Client Connect Failed")
		return
	}
	fmt.Println("Start")
	handleConnClient(conn)
	fmt.Println("End")
}

func handleConnClient(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(os.Stdin)
	help := `>cd path 移动到path路径下>ls 显示当前目录信息
	>push filename dsc 上传文件到云端  本地绝对路径 -> 云端绝对路径
	>pull filename dsc 从云端获取文件  云端相对路径 -> 本地绝对路径
	>mv src dsc 移动文件
	>cp src dsc 复制文件
	>del filename 删除文件
	>deldir dirname 删除目录及目录下所有文件
	>mkdir dirname 创建目录  绝对路径目录
	>share filename username
	>info 查看自己文件容器信息`
	for {
		input,err := reader.ReadString('\n')
		input = strings.Trim(input,"\n")
		if err != nil {
			continue
		}
		op := strings.Split(input," ")

		switch op[0] {
			case "push" :
				if len(op) > 2 {
					if !FileExist(op[1]) {
						fmt.Println("Error : File doesn't exist")
						continue
					}
					flag := JudgeOs()
					filepath := strings.Split(op[1], flag)
					filename := filepath[len(filepath)-1]
					input = op[0] + " " + op[1] + " " + op[2] + "/" + filename
					err = message.SendMsg([]byte(input), c)
					err = pushClient(op[1], c)
				} else {
					err = message.SendMsg([]byte(input),c)
				}
			case "pull" :
				if len(op) > 2 {
					err = message.SendMsg([]byte(input),c)
					err = pullClient(op[1],op[2],c)
				} else {
					err = message.SendMsg([]byte(input),c)
				}
			case "help":
				fmt.Println(help)
				continue
			default:
				err = message.SendMsg([]byte(input),c)
		}
		msg,err := message.RecvMsg(c)
		if err != nil {
			continue
		}
		fmt.Println(string(msg))
		if input == "exit" {
			return
		}
	}
}

func pushClient(filename string,c net.Conn) (error) {
	// 清空临时文件
	defer DirClear(config.Tmp)

	msg,err := message.RecvMsg(c)
	if err != nil {
		return err
	}
	if string(msg) == "Error" {
		return errors.New("Dsc isn't a dir")
	}

	fileInfo,err := os.Stat(filename)
	if err != nil {
		return err
	}
	err = message.SendMsg([]byte(strconv.FormatInt(fileInfo.Size(),10)),c)
	flag,err := message.RecvMsg(c)
	if err != nil {
		return err
	}
	if string(flag) == "Fail" {
		return errors.New("Push File Failed")
	}

	err = FileDivide(filename,config.Tmp,config.FileDivideSize)
	if err != nil {
		return err
	}
	files,err := ioutil.ReadDir(config.Tmp)
	if err != nil {
		return err
	}
	err = message.SendMsg([]byte(strconv.Itoa(len(files))),c)
	if err != nil {
		return err
	}
	for _,file := range files {
		msg,err := FileReadRaw(config.Tmp + file.Name())
		if err != nil {
			return err
		}
		err = message.SendMsg(msg,c)
		if err != nil {
			return err
		}
	}
	return nil
}

func pullClient(filename,dsc string,c net.Conn) (error) {
	// 清空临时文件
	defer DirClear(config.Tmp)

	i := 0
	flag := JudgeOs()
	for {
		data,err := message.RecvMsg(c)
		if err != nil {
			return err
		}
		if string(data) == "END" {
			break
		}
		err = FileWrite(config.Tmp + flag + "tmp" + strconv.Itoa(i),data)
		i += 1
		if err != nil {
			return err
		}
	}
	err := FileMerge(config.Tmp,dsc)
	if err != nil {
		return err
	}
	return nil
}