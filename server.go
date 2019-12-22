package main

import (
	"FileSystem/config"
	. "FileSystem/file"
	"FileSystem/message"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 开始监听
	l,err := net.Listen("tcp",config.Host)
	if err != nil {
		return
	}
	defer l.Close()

	fmt.Println("Start")
	for {
		conn,err := l.Accept()
		if err != nil {
			continue
		}
		go handleConnServer(conn)
	}
	fmt.Println("End")
}

func handleConnServer(c net.Conn) {
	defer c.Close()
	fmt.Println(c.RemoteAddr().String() + " is connected!")

	var filecontainer *FileContainer
	var filecontainer_pos string

	// 登录注册流程
	for {
		buf, err := message.RecvMsg(c)
		if err != nil {
			return
		}
		op := strings.Split(string(buf), " ")
		switch op[0] {
			case "regist":
				if len(op) > 2 {
					res, fcpos, err := regist(op[1], op[2])
					if res != "Regist Success" {
						err = message.SendMsg([]byte(res),c)
						if err != nil {
							return
						}
					} else {
						err = message.SendMsg([]byte(res),c)
						if err != nil {
							return
						}
						filecontainer,err = ReadFileContainer(fcpos)
						if err != nil {
							return
						}
						filecontainer_pos = fcpos
					}
				} else {
					err = message.SendMsg([]byte("Please Login First"),c)
					if err != nil {
						return
					}
				}
			case "login":
				if len(op) > 2 {
					res, fcpos,flag, err := login(op[1], op[2])
					if res != "Login Success" {
						err = message.SendMsg([]byte(res), c)
						if err != nil {
							return
						}
					} else {
						filecontainer, err = ReadFileContainer(fcpos)
						if err != nil {
							return
						}
						if flag {
							err = message.SendMsg([]byte(res + "\nDo You Want Get Share Files? Y/N"),c)
							if err != nil {
								return
							}
							for {
								data, err := message.RecvMsg(c)
								if err != nil {
									return
								}
								switch string(data) {
									case "Y":
										msg := shareGet(op[1],filecontainer)
										err = message.SendMsg([]byte(msg),c)
									case "y":
										msg := shareGet(op[1],filecontainer)
										err = message.SendMsg([]byte(msg),c)
									case "N":
										err = message.SendMsg([]byte("Share Files Will Remind You Next Time"),c)
									case "n":
										err = message.SendMsg([]byte("Share Files Will Remind You Next Time"),c)
									default:
										err = message.SendMsg([]byte("Do You Want Get Share Files? Y/N"), c)
									}
								if string(data) == "Y" || string(data) == "y" || string(data) == "N" || string(data) == "n" {
									break
									}
								}
						} else {
							err = message.SendMsg([]byte(res), c)
							if err != nil {
								return
							}
						}
						filecontainer_pos = fcpos
					}
				} else {
					err = message.SendMsg([]byte("Please Login First"),c)
					if err != nil {
						return
					}
				}
			case "exit":
				err = message.SendMsg([]byte("Connection Closed!"),c)
				fmt.Println(c.RemoteAddr().String() + " is disconnected!")
				return
			default:
				err = message.SendMsg([]byte("Please Login First"),c)
				if err != nil {
					return
				}
		}
		if filecontainer != nil {
			break
		}
	}

	cur := filecontainer.DirectoryTree	// 当前目录

	// 交互流程
	for {
		buf,err := message.RecvMsg(c)
		if err != nil {
			return
		}
		// 操作列表
		op := strings.Split(string(buf)," ")

		var msg string

		switch op[0] {
			case "cd":
				if len(op) >1 {
					msg = cd(op[1], &cur, filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "ls":
				msg = ls(cur)
			case "push":
				if len(op) > 2 {
					msg = push(op[2], filecontainer, c)
				} else {
					msg = "Operation Error"
				}
			case "pull":
				if len(op) > 2 {
					msg = pull(op[1], cur, filecontainer, c)
				} else {
					msg = "Operation Error"
				}
			case "mv":
				if len(op) > 2 {
					msg = mv(op[1],op[2],filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "cp":
				if len(op) > 2 {
					msg = cp(op[1],op[2],filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "del":
				if len(op) > 1 {
					msg = del(op[1], filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "deldir":
				if len(op) > 1 {
					msg = delDir(op[1],filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "mkdir":
				if len(op) > 1 {
					msg = mkdir(op[1],filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "share":
				if len(op) > 2 {
					msg = share(op[1],op[2],filecontainer)
				} else {
					msg = "Operation Error"
				}
			case "info":
				msg = info(filecontainer)
			case "exit":
				msg = "Connection Closed!"
			default:
				msg = " "
		}
		err = message.SendMsg([]byte(msg),c)
		if err != nil {
			return
		}
		if msg == "Connection Closed!" {
			break
		}
	}
	// 保证异常退出时用户容器信息更新
	WriteFileContainer(filecontainer,filecontainer_pos)
	fmt.Println(c.RemoteAddr().String() + " is disconnected!")
	return
}

// 注册 注册时分配一个空的文件容器
func regist(username,password string) (string,string,error) {
	db,err := sql.Open("mysql",config.DSN)
	if err != nil {
		return "Regist Failed","",err
	}
	defer db.Close()

	// err == nil代表查询到了该用户名
	var id string
	err = db.QueryRow("SELECT username FROM fc_management WHERE username = \"" + username + "\"").Scan(&id)
	if err == nil {
		return "Username already exists", "", err
	}

	stmt,err := db.Prepare("INSERT INTO fc_management SET username = ?,password = ?,fcpos = ?,sharepos = ?,sharestat = 0")
	if err != nil {
		return "Regist Failed","",err
	}

	filecontainer := &FileContainer{TotalSpace:1073741824,UsedSpace:0,Authority:[]string{username}, DirectoryTree:&Tree{Name:username}}
	fcpos := randomPath()
	sharepos := randomPath()
	_ , err = stmt.Exec(username,password,fcpos,sharepos)
	if err != nil {
		return "Regist Failed","",err
	}

	err = WriteFileContainer(filecontainer,fcpos)
	if err != nil {
		_,err = db.Exec("DELETE FROM fc_management WHERE username = \"" + username + "\"")
		return "Regist Failed","",err
	}
	err = FileCreate(sharepos)
	if err != nil {
		_,err = db.Exec("DELETE FROM fc_management WHERE username = \"" + username + "\"")
		return "Regist Failed","",err
	}
	return "Regist Success",fcpos,nil
}

// 登录
func login(username,password string) (string,string,bool,error) {
	var flag bool
	db,err := sql.Open("mysql",config.DSN)
	if err != nil {
		return "Login Failed","",flag,err
	}
	defer db.Close()

	row := db.QueryRow("SELECT password,fcpos,sharepos FROM fc_management WHERE username = \"" + username + "\"")
	var pswd, fcpos,sharepos string
	err = row.Scan(&pswd, &fcpos,&sharepos)
	if err != nil {
		return "Login Failed", "",flag,err
	}
	if pswd != password {
		return "Password incorrect", "", flag, nil
	}
	data,err := FileReadRaw(sharepos)
	if len(data) > 0 {
		flag = true
	}
	return "Login Success",fcpos,flag, nil
}

// 随机生成一个保存的路径
func randomPath() (string) {
	for {
		t := md5.New()
		t.Write([]byte(strconv.FormatInt(rand.Int63(), 10)))
		if !FileExist(config.FileAddress + hex.EncodeToString(t.Sum(nil))) {
			return config.FileAddress + hex.EncodeToString(t.Sum(nil))
		}
	}
}

func cd(dsc string,cur **Tree,filecontainer *FileContainer) string {
	if dsc == ".." {
		if (*cur).ParentNode != nil {
			*cur = (*cur).ParentNode
		}
	} else if dsc == "." {
		return GetPath(*cur)
	} else {
		path := GetPath(*cur)
		*cur = FindFile(filecontainer,path + dsc)
		// 当前目录下没找到 就找绝对路径
		if path + dsc + "/" != GetPath(*cur) {
			*cur = FindFile(filecontainer,dsc)
		}
	}
	return GetPath(*cur)
}

func ls(cur *Tree) (string) {
	var res string
	if cur.ChildNodes != nil && len(cur.ChildNodes) > 0{
		for _, t := range cur.ChildNodes {
			res += (" " + t.Name)
		}
		return res[1:]
	} else {
		return " "
	}
}

func push(dsc string,filecontainer *FileContainer,c net.Conn) (string) {
	// 判断是否是目录
	dscParts := strings.Split(dsc,"/")
	dscPath := strings.Join(dscParts[:len(dscParts)-1],"/")
	flag := IsDir(filecontainer,dscPath)
	if !flag {
		err := message.SendMsg([]byte("Error"),c)
		if err != nil {
			return "Error : Push File Failed"
		}
		return "Error : Push File Failed"
	} else {
		err := message.SendMsg([]byte("OK"),c)
		if err != nil {
			return "Error : Push File Failed"
		}
	}

	// 接收文件大小
	size,err := message.RecvMsg(c)
	size_int64,err := strconv.ParseInt(string(size),10,64)
	if err != nil {
		return "Error : Push File Failed"
	}
	if filecontainer.UsedSpace + size_int64 > filecontainer.TotalSpace {
		err = message.SendMsg([]byte("Fail"),c)
		return "Error : Space is not enough"
	} else {
		err = message.SendMsg([]byte("Success"),c)
		filecontainer.UsedSpace += size_int64
	}

	// 接收文件总数
	totalNum,err := message.RecvMsg(c)
	if err != nil {
		return "Error : Push File Failed"
	}
	num,err := strconv.Atoi(string(totalNum))

	if err != nil {
		return "Error : Push File Failed"
	}

	var preAddress string
	// 接收文件
	for i:=0;i<num;i++ {
		data,err := message.RecvMsg(c)
		if err != nil {
			return "Error : Push File Failed"
		}
		f := &File{PreAddress:preAddress,Content:data}
		preAddress = randomPath()
		err = FileWriteToDisk(f,preAddress)
		if err != nil {
			return "Error : Push File Failed"
		}
	}
	err = AddFile(filecontainer,dsc,preAddress,"")
	if err != nil {
		return "Error : Push File Failed"
	}
	return "Push File Successed"
}

func pull(filename string,cur *Tree,filecontainer *FileContainer,c net.Conn) (string) {
	path := GetPath(cur)
	pos := FindFile(filecontainer,path + filename)

	filepathParts := strings.Split(filename,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return "Error : Filename is wrong " + filename
	}

	file,err := FileReadFromDisk(pos.Address)
	if err != nil {
		return "Error : Pull File Failed"
	}

	for {
		err = message.SendMsg(file.Content,c)
		if err != nil {
			return "Error : Pull File Failed"
		}
		if file.PreAddress == "" {
			break
		} else {
			file,err = FileReadFromDisk(file.PreAddress)
			if err != nil {
				return "Error : Pull File Failed"
			}
		}
	}

	err = message.SendMsg([]byte("END"),c)

	if err != nil {
		return "Error : Pull File Failed"
	}
	return "Pull File Successed"
}

func mv(src,dsc string,filecontainer *FileContainer) (string) {
	err := MoveFile(filecontainer,src,dsc)
	if err != nil {
		return "Error : Move File Failed"
	}
	return "Move File Successed"
}

func cp(src,dsc string,filecontainer *FileContainer) (string) {
	err := CopyFile(filecontainer,src,dsc)
	if err != nil {
		return "Error : Copy File Failed"
	}
	return "Copy File Successed"
}

func del(filename string,filecontainer *FileContainer) (string) {
	pos := FindFile(filecontainer,filename)

	filepathParts := strings.Split(filename,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return "Error : Filename is wrong " + filename
	}

	var curAddress string = pos.Address

	if curAddress == "" {
		return "Error : This is a Dir,Please Use deldir " + filename
	}

	var filesize int64 = 0

	for {
		fileInfo,err := os.Stat(curAddress)
		if err != nil {
			return "Error : Delete File Failed"
		}
		filesize += fileInfo.Size()
		file,err := FileReadFromDisk(curAddress)
		if err != nil {
			return "Error : Delete File Failed"
		}
		err = os.Remove(curAddress)
		if file.PreAddress == "" {
			break
		} else {
			curAddress = file.PreAddress
		}
		if err != nil {
			return "Error : Delete File Failed"
		}
	}
	err := DeleteFile(filecontainer,filename)
	if err != nil {
		return "Error : Delete File Failed"
	}
	filecontainer.UsedSpace -= filesize
	return "Delete File Success"
}

func recdel(pos *Tree,filecontainer *FileContainer) (error) {
	for _,child := range pos.ChildNodes {
		if child.Address == "" {
			err := recdel(child,filecontainer)
			if err != nil {
				return err
			}
		} else {
			name := GetPath(child)
			name = name[:len(name)-1]
			res := del(name, filecontainer)
			if res == "Delete File Success" {
				continue
			} else {
				return errors.New("Delete File Failed")
			}
		}
	}
	dirname := GetPath(pos)
	err := DeleteFile(filecontainer,dirname[:len(dirname)-1])
	if err != nil {
		return err
	}
	return nil
}

func delDir(dirname string,filecontainer *FileContainer) (string) {
	pos := FindFile(filecontainer,dirname)

	filepathParts := strings.Split(dirname,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return "Error : Dirname is wrong " + dirname
	}
	if pos.Address != "" {
		return "Error : This is a File,Please Use del " + dirname
	}
	err := recdel(pos,filecontainer)
	if err != nil {
		return "Error : Del Dir Failed"
	}
	return "Del Dir Success"
}

func mkdir(dirname string,filecontainer *FileContainer) (string) {
	parent := strings.Split(dirname,"/")
	parentPath := strings.Join(parent[:len(parent)-1],"")
	if IsDir(filecontainer,parentPath) {
		err := AddFile(filecontainer,dirname,"","")
		if err != nil {
			return "Error : Make Dir Failed"
		}
		return "Make Dir Successed"
	} else {
		return "Error : Make Dir Failed"
	}
}

func recshare(pos *Tree,filecontainer *FileContainer,lastPath string) (error) {
	err := AddFile(filecontainer,lastPath + "/" + pos.Name,pos.Address,pos.Md5Address)
	if err != nil {
		return err
	}
	lastPath = lastPath + "/" + pos.Name
	for _,child := range pos.ChildNodes {
		if child.Address == "" {
			err := recshare(child,filecontainer,lastPath)
			if err != nil {
				return err
			}
		} else {
			err := AddFile(filecontainer,lastPath + "/" + child.Name,child.Address,child.Md5Address)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func share(filename,username string,filecontainer *FileContainer) (string) {
	pos := FindFile(filecontainer,filename)

	filepathParts := strings.Split(filename,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return "Error : Filename is wrong " + filename
	}

	db,err := sql.Open("mysql",config.DSN)
	if err != nil {
		return "Error : Share File Failed"
	}
	defer db.Close()

	var sharepos string
	var sharestat int
	for {
		row := db.QueryRow("SELECT sharepos,sharestat FROM fc_management WHERE username = \"" + username + "\"")
		err = row.Scan(&sharepos, &sharestat)
		if sharepos == "" {
			return "Error : User doesn't exist"
		}
		if sharestat == 0 {
			_,err = db.Exec("UPDATE fc_management SET sharestat = 1 WHERE username = \"" + username + "\"")
			if err != nil {
				return "Error : Share File Failed"
			}
			break
		}
		time.Sleep(time.Second)
	}
	defer db.Exec("UPDATE fc_management SET sharestat = 0 WHERE username = \"" + username + "\"")
	_tree := &Tree{Name:filecontainer.Authority[0],ChildNodes:[]*Tree{pos}}
	tree := &Tree{Name:"share",ChildNodes:[]*Tree{_tree}}
	err = WriteShareContainer(tree,sharepos)
	if err != nil {
		return "Error : Share File Failed"
	}
	return "Share File Success"
}

func shareGet(username string,filecontainer *FileContainer) (string) {
	db,err := sql.Open("mysql",config.DSN)
	if err != nil {
		return "Error : Share File Failed"
	}
	defer db.Close()

	var sharepos string
	var sharestat int
	for {
		row := db.QueryRow("SELECT sharepos,sharestat FROM fc_management WHERE username = \"" + username + "\"")
		err = row.Scan(&sharepos, &sharestat)
		if sharestat == 0 {
			_,err = db.Exec("UPDATE fc_management SET sharestat = 1 WHERE username = \"" + username + "\"")
			if err != nil {
				return "Error : Get Share File Failed"
			}
			break
		}
		time.Sleep(time.Second)
	}
	defer db.Exec("UPDATE fc_management SET sharestat = 0 WHERE username = \"" + username + "\"")
	trees,err := ReadShareContainer(sharepos)
	if err != nil {
		return "Error : Get Share File Failed"
	}
	for _,tree := range trees {
		err = recshare(tree,filecontainer,filecontainer.Authority[0])
		if err != nil {
			return "Error : Get Share File Failed"
		}
	}
	err = FileWrite(sharepos,[]byte(""))
	return "Get Share File Successs"
}

func info(filecontainer *FileContainer) (string) {
	totalSpace := "TotalSpace : " + strconv.FormatInt(filecontainer.TotalSpace / 1024,10) + " MB"
	usedSpace := "UsedSpace : " + strconv.FormatInt(filecontainer.UsedSpace / 1024,10) + " MB"
	authority := "Authority : " + strings.Join(filecontainer.Authority,",")
	return totalSpace + "\n" + usedSpace + "\n" + authority + "\n"
}