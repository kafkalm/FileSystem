package file
type Tree struct {
	Name       string // 文件/目录名
	ParentNode *Tree
	ChildNodes []*Tree
	Address    string // 文件地址 没有文件地址即可判断这个是个目录
	Md5Address string // md5校验码地址
}

// 文件容器
type FileContainer struct {
	TotalSpace    int64    //总空间
	UsedSpace     int64    //已用空间
	Authority     []string //权限信息
	DirectoryTree *Tree    //文件目录结构树
}

// 文件
type File struct {
	Content []byte			//内容
	PreAddress string		//上一个文件位置
}

