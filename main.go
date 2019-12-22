package main

//import (
//	. "FileSystem/file"
//	"fmt"
//)
//
//func main() {
//	container := &FileContainer{TotalSpace: 100000, UsedSpace:0, Authority:[]string{"kafkal"}}
//	tree := &Tree{Name: "root", Address:"/", Md5Address:"/"}
//	tree.ChildNodes = append(tree.ChildNodes,&Tree{Name: "son1", ParentNode:tree, Address:"son1", Md5Address:"md5son1"},
//		&Tree{Name: "son2", ParentNode:tree, Address:"son2", Md5Address:"md5son2"})
//	container.DirectoryTree = tree
//	AddFile(container,"root/son1/test1.txt")
//	WriteFileContainer(container,"C:\\Users\\10517\\Desktop\\test.txt")
//	fmt.Println("写入成功")
//	_container,_ := ReadFileContainer("C:\\Users\\10517\\Desktop\\test.txt")
//	DeleteFile(_container,"root/son1/test1.txt")
//	AddFile(_container,"root/son1/test1.txt")
//	MoveFile(_container,"root/son1/test1.txt","root/son2/test1.txt")
//	MoveFile(_container,"root/son1/test1.txt","root/son2/test1.txt")
//	WriteFileContainer(_container,"C:\\Users\\10517\\Desktop\\test.txt")
//	fmt.Println(_container)
//}

//import (
//	. "FileSystem/encrypt"
//	"fmt"
//)
//
//func main() {
//	//s := Divide_64("test")
//	key := "1234567"
//	s := "test"
//	ciphertext := Encrypt(s,key)
//	fmt.Println(ciphertext)
//	plaintext := Decrypt(ciphertext,key)
//	fmt.Println(plaintext)
//}

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Join([]string{"1","2","3"},","))
	//if err := FileDivide("C:\\Users\\10517\\Desktop\\毕设数据结构.txt","C:\\Users\\10517\\Desktop\\tmp\\",40);err != nil{
	//	fmt.Println("Fail")
	//} else {
	//	fmt.Println("Success")
	//}
	//file,err := FileReadFromDisk("D:\\FStmp\\1bca46b45f441562f9d54736f165e5e1")
	//if err != nil {
	//	return
	//}
	//fmt.Println(file)
	//if err := FileMerge("C:\\Users\\10517\\Desktop\\tmp\\","C:\\Users\\10517\\Desktop\\tmp\\");err != nil{
	//	fmt.Println("Fail")
	//} else {
	//	fmt.Println("Success")
	//}
	//data,err := file.FileRead("C:Users\\10517\\Desktop\\test.txt")
	//if err != nil {
	//	return
	//}
	//op := strings.Split(" "," ")
	//fmt.Println(op)
	//fmt.Println(data)
	//if _,err:= FileRead("‪C:\\Users\\10517\\Desktop\\tmp\\1.tmp");err != nil{
	//	fmt.Println("Fail")
	//}
	//tree := &Tree{Name: "root", Address:"/", Md5Address:"/"}
	//tree.ChildNodes = append(tree.ChildNodes,&Tree{Name: "son1", ParentNode:tree, Address:"son1", Md5Address:"md5son1"},
	//	&Tree{Name: "son2", ParentNode:tree, Address:"son2", Md5Address:"md5son2"})
	//cur := tree.ChildNodes[0]
	//fmt.Println(GetPath(cur))
}

//import (
//	"fmt"
//)
//
//func main() {
//	fmt.Println(runtime.GOOS)
//}