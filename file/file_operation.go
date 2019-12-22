package file

import (
	"FileSystem/config"
	"FileSystem/encrypt"
	"bufio"
	"container/list"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// *** 基本文件读写操作封装 ***
// 创建文件
func FileCreate(filename string) (error) {
	file,err := os.Create(filename)
	defer file.Close()
	return err
}

// 读文件 （按行读 去除换行符）
func FileRead(filename string) ([]string,error) {
	file,err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	defer file.Close()

	var data []string
	// 按行读取 存入data 1.总空间 2.已用空间 3.权限信息 4.文件目录结构树的序列化形式
	buf := bufio.NewReader(file)
	for {
		line,err := buf.ReadString('\n')
		data = append(data,strings.Trim(line,"\n"))
		if err == io.EOF {
			break
		} else if err != nil{
			return nil,err
		}
	}
	return data,nil
}

// 读文件
func FileReadRaw(filename string) ([]byte,error) {
	if data,err := ioutil.ReadFile(filename);err==nil {
		return data,nil
	} else {
		return nil,err
	}
}

// 写文件
func FileWrite(filename string,data []byte) (error) {
	if err := ioutil.WriteFile(filename,data,0666);err == nil{
		return nil
	} else {
		return err
	}
}

// 追加写文件
func FileAppend(filename string,data []byte) (error) {
	file,err := os.OpenFile(filename,os.O_APPEND,0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_ , err = file.Write(data)
	return err
}

// 删除文件元信息
func FileDeleteMeta(filename string,rows int) (error) {
	file,err := os.OpenFile(filename,os.O_RDWR,0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var deleteMeta []byte
	buf := bufio.NewReader(file)
	for i:=0;i<rows;i++ {
		meta,err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		deleteMeta = append(deleteMeta,meta...)
	}
	data,err := FileReadRaw(filename)
	return FileWrite(filename,data[len(deleteMeta):])
}

// 判断文件是否存在
func FileExist(filename string) (bool) {
	_,err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 判断操作系统 确定连接符
func JudgeOs() (string) {
	var flag string
	if runtime.GOOS == "windows" {
		flag = "\\"
	} else {
		flag = "/"
	}
	return flag
}

// 清空文件夹
func DirClear(dirname string) (error) {
	files,err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	flag := JudgeOs()
	for _,file := range files {
		if runtime.GOOS == "windows" {
			os.Remove(dirname + flag + file.Name())
		} else {
			os.Remove(dirname + flag + file.Name())
		}
	}
	return nil
}

// *** FileContainer 结构体操作 ***
// 序列化目录树
func serialize(root *Tree) (string){
	// 没有节点用 N 来表示
	if root == nil{
		return "N,"
	}
	// 序列化形式为 “节点名称#节点文件地址#节点md5校验码地址#子节点数,"
	s := root.Name + "#" + root.Address + "#" + root.Md5Address + "#" + strconv.Itoa(len(root.ChildNodes)) + ","
	for _,child := range root.ChildNodes {
		s += serialize(child)
	}
	return s
}

// 反序列化目录树
func deserialize(s string) (*Tree){
	nodes := strings.Split(s,",")
	nodes = nodes[:len(nodes)-1]
	queue := list.New()
	for _,node := range nodes {
		queue.PushBack(node)
	}
	return _deserialize(nil,queue)
}

func _deserialize(parent *Tree,queue *list.List) (*Tree){
	cur := queue.Front()
	queue.Remove(cur)
	if cur.Value.(string) == "N" {
		return nil
	}
	parts := strings.Split(cur.Value.(string),"#")
	size,_ := strconv.Atoi(parts[3])
	node := &Tree{Name: parts[0], ParentNode:parent, Address:parts[1], Md5Address:parts[2]}
	for i:=0 ; i<size ; i++ {
		node.ChildNodes = append(node.ChildNodes,_deserialize(node,queue))
	}
	return node
}

// 文件容器的反序列化读取
func ReadFileContainer(filename string) (*FileContainer,error){
	if data,err := FileRead(filename);err == nil {
		container := &FileContainer{}
		container.TotalSpace, _ = strconv.ParseInt(data[0],10,64)
		container.UsedSpace, _ = strconv.ParseInt(data[1],10,64)
		container.Authority = strings.Split(data[2],",")
		container.DirectoryTree = deserialize(data[3])

		return container,nil
	} else {
		return nil,errors.New("Read FileContainer Fails")
	}
}

// 文件容器的序列化写入
func WriteFileContainer(container *FileContainer,filename string) (error){
	s := strconv.FormatInt(container.TotalSpace,10) + "\n" + strconv.FormatInt(container.UsedSpace,10) + "\n" +
		strings.Join(container.Authority,",") + "\n" + serialize(container.DirectoryTree)
	data := []byte(s)
	return FileWrite(filename,data)
}

// Share容器的反序列化读取
func ReadShareContainer(filename string) ([]*Tree,error){
	var trees []*Tree
	if data,err := FileRead(filename);err == nil {
		for _,line := range data[:len(data)-1]{
			trees = append(trees,deserialize(line))
		}
		return trees,nil
	} else {
		return nil,errors.New("Read ShareContainer Fails")
	}
}
// Share容器的反序列化写入
func WriteShareContainer(tree *Tree,filename string) (error){
	s := serialize(tree) + "\n"
	data := []byte(s)
	return FileAppend(filename,data)
}

// 寻找文件在树中的位置
func FindFile(container *FileContainer,filePath string) (*Tree) {

	root := container.DirectoryTree	//根目录

	filepathParts := strings.Split(filePath,"/")
	filepathParts = filepathParts[1:]		// 根目录即为文件容器

	for i:=0 ; i < len(filepathParts) ; i++ {
		for _,child := range root.ChildNodes {
			if child.Name == filepathParts[i]{
				root = child
				break
			}
		}
	}

	return root
}

// 添加文件(对文件目录信息的操作，不涉及实际写磁盘)
func AddFile(container *FileContainer,filePath,address,md5address string) (error) {
	filepathParts := strings.Split(filePath,"/")
	fileName := filepathParts[len(filepathParts)-1]		//文件名
	new_filePath := strings.Join(filepathParts[:len(filepathParts)-1],"/")		//文件路径
	pos := FindFile(container,new_filePath)

	if len(filepathParts) < 2 || pos.Name != filepathParts[len(filepathParts)-2] {
		return errors.New("filePath is wrong")
	}

	for _,child := range pos.ChildNodes {
		if child.Name == fileName {
			return errors.New("File Has Been Existed")
		}
	}
	pos.ChildNodes = append(pos.ChildNodes,&Tree{Name:fileName,ParentNode:pos,Address:address,
		Md5Address:md5address})
	return nil
}

// 删除文件(对文件目录信息的操作，不涉及实际写磁盘)
func DeleteFile(container *FileContainer,filePath string) (error) {
	filepathParts := strings.Split(filePath,"/")
	fileName := filepathParts[len(filepathParts)-1] //文件名
	pos := FindFile(container,filePath)

	if pos.Name != fileName {
		return errors.New("filePath is wrong")
	}

	parent := pos.ParentNode
	index := -1
	for i,child := range parent.ChildNodes {
		if child.Name == fileName {
			index = i
			break
		}
	}
	parent.ChildNodes = append(parent.ChildNodes[:index],parent.ChildNodes[index+1:]...)
	return nil
}

// 拷贝文件(实质上只要修改文件目录信息表即可，无需移动磁盘文件)
func CopyFile(container *FileContainer,srcPath,dscPath string) (error) {
	pos := FindFile(container,srcPath)
	dscPos := FindFile(container,dscPath)

	filepathParts := strings.Split(srcPath,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return errors.New("srcPath is wrong")
	}
	filepathParts = strings.Split(dscPath,"/")
	_dirname = filepathParts[len(filepathParts)-1]
	if dscPos.Name != _dirname {
		return errors.New("dscPath is wrong")
	}
	if dscPos.Address != "" {
		return errors.New("dscPath is not a dir")
	}
	return AddFile(container,dscPath + "/" + pos.Name,pos.Address,pos.Md5Address)
}

// 移动文件(与拷贝文件相同，只要修改文件目录信息表，无需移动磁盘文件)
func MoveFile(container *FileContainer,srcPath,dscPath string) (error) {
	pos := FindFile(container,srcPath)
	dscPos := FindFile(container,dscPath)

	filepathParts := strings.Split(srcPath,"/")
	_dirname := filepathParts[len(filepathParts)-1]
	if pos.Name != _dirname {
		return errors.New("srcPath is wrong")
	}
	filepathParts = strings.Split(dscPath,"/")
	_dirname = filepathParts[len(filepathParts)-1]
	if dscPos.Name != _dirname {
		return errors.New("dscPath is wrong")
	}
	if dscPos.Address != "" {
		return errors.New("dscPath is not a dir")
	}
	if err := AddFile(container,dscPath + "/" + pos.Name,pos.Address,pos.Md5Address);err == nil{
		return DeleteFile(container,srcPath)
	} else {
		return err
	}
}

// 检测是否是目录
func IsDir(container *FileContainer,filePath string) (bool) {
	pos := FindFile(container,filePath)
	filepathParts := strings.Split(filePath,"/")
	if pos.Name != filepathParts[len(filepathParts)-1] {
		return false
	}
	if pos != nil && pos.Address == "" {
		return true
	}
	return false
}

// *** File 结构体操作 ***
// 磁盘文件读取
func FileReadFromDisk(path string) (*File,error) {
	if data,err := FileRead(path);err == nil {
		file := &File{}
		file.PreAddress = data[0]
		for _,d := range data[1:len(data)-1] {
			file.Content = append(file.Content,[]byte(d+"\n")...)
		}
		file.Content = append(file.Content,[]byte(data[len(data)-1])...)
		return file,nil
	}	else {
		return nil,err
	}
}

// 磁盘文件写入
func FileWriteToDisk(file *File,path string) (error) {
	s := file.PreAddress + "\n"
	data := []byte(s)
	if FileWrite(path,data) == nil && FileAppend(path,file.Content) == nil{
		return nil
	} else {
		return errors.New("FileWriteToDisk Fails , File Name is " + path)
	}
}

// 磁盘文件移动 因为做了小文件处理 直接读写文件即可
// 与FileContainer的操作一起使用
// 在文件目录树中已经判断了是否存在文件 因此这个操作只涉及读写 不用再判断
func FileMoveInDisk(src,dsc string) (error) {
	if file,err := FileReadFromDisk(src);err == nil {
		if err2 := FileWriteToDisk(file,dsc);err2 == nil {
			return os.Remove(src)
		} else {
			return err2
		}
	} else {
		return err
	}
}

// 磁盘文件复制
func FileCopyInDisk(src,dsc string) (error) {
	if file,err := FileReadFromDisk(src);err == nil {
			return FileWriteToDisk(file,dsc)
	} else {
		return err
	}
}

// *** 磁盘文件操作 ***
// 磁盘文件切分
func FileDivide(src,dsc string,size int64) (error) {
	fileInfo,_ := os.Stat(src)
	file,err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileNum,totalNum int64
	fileNum = 0
	totalNum = fileInfo.Size() / size + 1
	_totalNum := strconv.FormatInt(totalNum,10)
	// 按size分割文件后，要将最后多出的size删除
	leftSize := totalNum * size - fileInfo.Size()

	key,err := GetKey(config.DESKEY)

	if err != nil {
		return err
	}

	for ;fileNum < totalNum; fileNum++{
		// 基本信息 1.编号 2.文件名 3.总编号
		baseInfo := strconv.FormatInt(fileNum,10) + "\n" + fileInfo.Name() + "\n" + _totalNum + "\n"
		baseInfoByte := []byte(baseInfo)

		buf := make([]byte,size)
		if fileNum == totalNum -1 {
			buf = make([]byte,size - leftSize)
		}

		_,err := file.Read(buf)
		if err != nil {
			return err
		}
		data := append(baseInfoByte,buf...)
		cryptedData,err := encrypt.DESEncrypt(data,key)
		if err != nil {
			return err
		}
		if err := FileWrite(dsc + strconv.FormatInt(fileNum,10) + ".tmp",cryptedData);err != nil {
			return err
		}
	}
	return nil
}

// 磁盘文件合并
func FileMerge(src,dsc string) (error) {
	rd,err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	var fileName string
	var totalNum int = 0

	key,err := GetKey(config.DESKEY)

	if err != nil {
		return err
	}
	//按次序重命名文件 并删去文件元信息
	for _,f := range rd {
		cryptedData,err := FileReadRaw(src + f.Name())
		if err != nil {
			return err
		}
		originData,err := encrypt.DESDecrypt(cryptedData,key)
		if err != nil {
			return err
		}
		err = FileWrite(src + f.Name(),originData)
		if err != nil {
			return err
		}
		if data,err := FileRead(src + f.Name());err != nil{
			return err
		} else {

			fileNum := data[0]
			fileName = data[1]
			totalNum,_ = strconv.Atoi(data[2])
			if err := FileDeleteMeta(src + f.Name(),3);err != nil{
				return err
			}
			if err := os.Rename(src + f.Name(),src + fileNum + ".tmp");err != nil {
				return err
			}
		}
	}

	if err := FileCreate(dsc + fileName);err != nil{
		return err
	}

	for i:=0;i < totalNum;i++ {
		data,err := FileReadRaw(src + strconv.Itoa(i) + ".tmp")
		if err != nil {
			return err
		} else {
			if err := FileAppend(dsc + fileName,data);err != nil{
				return err
			}
		}
	}

	return nil
}

// *** Tree 结构体操作 ***
// 还原绝对路径
func GetPath(root *Tree) (string) {
	var s string
	for {
		if root != nil {
			s = root.Name + "/" + s
			root = root.ParentNode
		} else {
			break
		}
	}
	return s
}

// *** 密钥操作 ***
func GetKey(filename string) ([]byte,error) {
	key,err := FileReadRaw(config.DESKEY)
	if err != nil {
		return nil,err
	}
	return key,nil
}