#云存储的私有文件保护系统
文件基本操作已完成  
>1.文件容器信息表操作 [x]  
>2.文件结构体操作  [x]  
>3.磁盘文件操作 [x]
---
加密模块  
在服务器无法知晓用户文件的情况下，可以不使用RSA来保护用户的DES密钥，如果服务器需要知晓，可采用RSA密钥来进行保护
，相关的操作可直接修改file/file_operation.go中的GetKey函数来实现
>1.为熟悉Golang 自己写了DES  
>2.RSA可采用Golang已有的库  
>
---
服务端 
操作集合  
>cd path 移动到path路径下  
>ls 显示当前目录信息  
>push filename dsc 上传文件到云端  本地绝对路径 -> 云端绝对路径  
>pull filename dsc 从云端获取文件  云端相对路径 -> 本地绝对路径  
>mv src dsc 移动文件  
>cp src dsc 复制文件  
>del filename 删除文件  
>deldir dirname 删除目录及目录下所有文件
>mkdir dirname 创建目录  绝对路径目录  
>share filename username  
>info 查看自己文件容器信息  