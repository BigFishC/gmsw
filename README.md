# 说明
---
## 功能
* 加密
     * 将给定的字符串以AES的方式进行加密，并把加密之后的字符串写入对应的文件中

* 传输
     * 将本地文件传输至目标端

* 执行命令
     * 实现kill service功能
     * 实现本地执行命令功能

## 具体使用方法

+ 帮助  
[root@localhost]# gmsf -h

+ 加密  
[root@localhost]# gmsf encrypt --tpwd='1234'

+ 远程执行命令  
[root@localhost]# gmsf cmd -t user@ip 'cmdline'

+ 远程执行命令(非22端口)  
[root@localhost]# gmsf cmd -P 端口 -t user@ip 'cmdline'

+ 传输文件  
[root@localhost]# gmsf cmd -T 文件完成路径 -t user@ip '目标段文件完成路径'   //默认5个并发传输

+ 传输文件（非22端口）  
[root@localhost]# gmsf cmd -P 端口 -T 文件完成路径 -t user@ip '目标段文件完成路径' //默认5个并发传输

+ 关闭服务  
[root@localhost]# gmsf kill servicename

+ 检测服务  
[root@localhost]# gmsf check servicename

+ 启动服务  
[root@localhost]# gmsf start -d '启动目录' '启动命令'

