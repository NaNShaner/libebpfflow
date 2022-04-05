# libebpfflow
Traffic visibility library based on eBPF

### 详情请跳转
[libebpfflow 官方仓库](https://github.com/ntop/libebpfflow)


### 本仓库功能
针对ebpflowexport非结构化的输出，进行接续解析为json进行输出
 
### 编译
You need a modern eBPF-enabled Linux distribution.

On Ubuntu 16.04/18.04/20.04 Server LTS you can install the prerequisites (we assume that the compiler is already installed) as follows:
```sh
$ cd /opt/
$ sudo git clone https://github.com/NaNShaner/libebpfflow.git
$ cd libebpfflow/go/fmtflow/ && go build .
```

### 使用方式
`ebpflowexport` 不可以带参数
```sh
$ ./fmtflow -c /opt/libebpfflow/libebpfflow-master/ebpflowexport
```


### Testing
```sh
$ sudo ./fmtflow -h
Usage of ./fmtflow [ OPTIONS ]:
  -c string
    	请输入ebpflowexport可执行文件的绝对路径
Note: please run as root 
```