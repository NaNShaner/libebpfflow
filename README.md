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


### 使用及输出样例
```sh
$ sudo ./fmtflow -h
Usage of ./fmtflow [ OPTIONS ]:
  -c string
    	请输入ebpflowexport可执行文件的绝对路径
Note: please run as root 
```
输出结果如下：
```json
{
  "event_time": "1647334436.451203",
  "ifname": "enp0s3",
  "packet_action": "Sent",
  "proto": "IPv4/TCP",
  "task_info": {
    "Pid": "3131",
    "Tid": "3096",
    "FullTaskPath": "/usr/local/bin/kube-proxy --config=/var/lib/kube-proxy/config.conf",
    "Uid": "0",
    "Gid": "0"
  },
  "father_task_info": {
    "Pid": "3077",
    "Tid": "0",
    "FullTaskPath": "/usr/bin/containerd-shim-runc-v2-namespace",
    "Uid": "0",
    "Gid": "0"
  },
  "connect_info": {
    "Saddr": "192.168.3.182:53438",
    "Daddr": "192.168.3.148:6443"
  },
  "container_info": {
    "container_id": "fdf3cfd11d8658e6ba93b7eb48e535c35e7b295075977d27c1d390c4",
    "docker_name": "k8s_kube-proxy_kube-proxy-bhrbb_kube-system_6f0aa01c-75b9-48e4-86fb-3f68f50b7023_7",
    "kube_name": "kube-proxy",
    "pod_name": "kube-proxy-bhrbb",
    "kube_name_space": "kube-system\n"
  },
  "connect_status": "CONNECT",
  "latency": "0.29"
}
```