/*
 *
 * @desc 用于格式化 ebpflowexport 的输出结果 用json的格式输出
 *
 */
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Ebpfflow 用于整体存储ebpflowexport的输出结果
type Ebpfflow struct {
	EventTime      string        `json:"event_time"`
	Ifname         string        `json:"ifname"`
	PacketAction   string        `json:"packet_action"`
	Proto          string        `json:"proto"`
	TaskInfo       TaskInfo      `json:"task_info"`
	FatherTaskInfo TaskInfo      `json:"father_task_info"`
	ConnectInfo    ConnectInfo   `json:"connect_info"`
	ContainerInfo  ContainerInfo `json:"container_info,omitempty"`
	ConnectStatus  string        `json:"connect_status"`
	Latency        string        `json:"latency,omitempty"`
}

// TaskInfo 存储当前进程以及其父进程的信息
type TaskInfo struct {
	Pid          string
	Tid          string
	FullTaskPath string
	Uid          string
	Gid          string
}

// ConnectInfo 存储请求的源目地址
type ConnectInfo struct {
	Saddr string
	Daddr string
}

// ContainerInfo 存储容器请求的源目地址
//[containerID: 45c29b0ed9600fb0ded3dac02ad54c97b4351a0389edd60ca31b87d8]
//[docker_name: k8s_orders_orders-6ddd96f65-8hf86_px-sock-shop_4bb15661-8804-4bf0-9ed0-2007e03fdb79_611]
//[kube_name: orders]
//[kube_pod: orders-6ddd96f65-8hf86]
//[kube_ns: px-sock-shop]
type ContainerInfo struct {
	ContainerID   string `json:"container_id,omitempty"`
	DockerName    string `json:"docker_name,omitempty"`
	KubeName      string `json:"kube_name,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	KubeNameSpace string `json:"kube_name_space,omitempty"`
}

func main() {
	_cmdEbpflowexport := flag.String("c", "", "请输入ebpflowexport可执行文件的绝对路径")
	flag.Parse()
	_, err := os.Stat(*_cmdEbpflowexport)
	if err != nil {
		log.Fatalf("%s文件不存在", *_cmdEbpflowexport)
	}

	cmd := exec.Command("sh", "-c", *_cmdEbpflowexport)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		listSplitReslut, _ := getFieldToStuct(line)
		// checkErr(err, 3)
		fmt.Printf("%s\n", PrintResultJson(listSplitReslut))

	}

	err = cmd.Wait()
	if err != nil {
		return
	}
	return
}

// getFieldToStuct 对于每一行输出进行格式化
func getFieldToStuct(f string) (Ebpfflow, error) {
	var ebpfflow Ebpfflow
	var flowSplit []string
	line := f
	if strings.Contains(line, "pid") {
		lineSplit := strings.Split(line, "][")
		for _, v := range strings.Split(lineSplit[0], " [") {
			flowSplit = append(flowSplit, v)
		}
		for i := 1; i < len(lineSplit); i++ {
			flowSplit = append(flowSplit, lineSplit[i])
		}
		fmt.Printf("%s\n", line)

		ebpfflow.EventTime = flowSplit[0]
		ebpfflow.Ifname = flowSplit[1]
		ebpfflow.PacketAction = flowSplit[2]
		ebpfflow.Proto = flowSplit[3]
		task, err := splitLine(&flowSplit[4])
		checkErr(err, -1)
		ebpfflow.TaskInfo = task
		fatherTask, err := splitLine(&flowSplit[5])
		checkErr(err, 1)
		ebpfflow.FatherTaskInfo = fatherTask
		connetInfo, err := connetInfoSplit(&flowSplit[6])
		checkErr(err, 2)
		ebpfflow.ConnectInfo = connetInfo

		if len(flowSplit) > 6 && len(flowSplit) <= 9 {
			if len(flowSplit) == 6 || strings.Contains(flowSplit[3], "UDP") {
				return ebpfflow, nil
			} else {
				ebpfflow.ConnectStatus = strings.Replace(flowSplit[7], "]", "", -1)
				if strings.Contains(line, "latency") {
					latencySplit := strings.Split(flowSplit[8], " ")
					ebpfflow.Latency = latencySplit[1]
				}
			}
		} else if len(flowSplit) > 9 {
			if strings.Contains(flowSplit[3], "TCP") && strings.Contains(line, "CONNECT") && !strings.Contains(line, "CONNECT_FAILED") {
				//for i, s := range flowSplit {
				//	fmt.Printf("flowSplit TCP CONNECT !CONNECT_FAILED: %d --> %s\n", i, s)
				//}
				ebpfflow.ConnectStatus = strings.Replace(flowSplit[7], "]", "", -1)
				latencySplit := strings.Split(flowSplit[8], " ")
				ebpfflow.Latency = latencySplit[1]
				containerID := containerInfoSplit(flowSplit[9])
				dockerName := containerInfoSplit(flowSplit[10])
				kubeName := containerInfoSplit(flowSplit[11])
				kubePod := containerInfoSplit(flowSplit[12])
				kubeNs := containerInfoSplit(flowSplit[13])

				ebpfflow.ContainerInfo = ContainerInfo{
					ContainerID:   containerID,
					DockerName:    dockerName,
					KubeName:      kubeName,
					KubeNameSpace: strings.Replace(kubeNs, "]", "", -1),
					PodName:       kubePod,
				}
			} else if strings.Contains(flowSplit[3], "TCP") && strings.Contains(line, "containerID") && !strings.Contains(flowSplit[7], "containerID:") {

				containerID := containerInfoSplit(flowSplit[8])
				dockerName := containerInfoSplit(flowSplit[9])
				kubeName := containerInfoSplit(flowSplit[10])
				kubePod := containerInfoSplit(flowSplit[11])
				kubeNs := containerInfoSplit(flowSplit[12])

				ebpfflow.ContainerInfo = ContainerInfo{
					ContainerID:   containerID,
					DockerName:    dockerName,
					KubeName:      kubeName,
					KubeNameSpace: strings.Replace(kubeNs, "]", "", -1),
					PodName:       kubePod,
				}
			} else if strings.Contains(flowSplit[3], "UDP") && strings.Contains(line, "containerID") {
				containerID := containerInfoSplit(flowSplit[7])
				dockerName := containerInfoSplit(flowSplit[8])
				kubeName := containerInfoSplit(flowSplit[9])
				kubePod := containerInfoSplit(flowSplit[10])
				kubeNs := containerInfoSplit(flowSplit[11])

				ebpfflow.ContainerInfo = ContainerInfo{
					ContainerID:   containerID,
					DockerName:    dockerName,
					KubeName:      kubeName,
					KubeNameSpace: strings.Replace(kubeNs, "]", "", -1),
					PodName:       kubePod,
				}
			}
		} else {
			for i, s := range flowSplit {
				fmt.Printf("else: %d --> %s\n", i, s)
			}
			// 如有未命中所有判断条件的情况下，报错退出，并输出当前解析报错行的内容
			log.Fatalf("%q\n", flowSplit)
		}

		return ebpfflow, nil
	} else {
		return Ebpfflow{}, fmt.Errorf("文本行解析失败，%s", line)
	}

}

// checkErr 通用错误解析
func checkErr(err error, errCode int) {
	if err != nil {
		fmt.Printf("解析文件异常，退出码为%d,异常信息：%v", errCode, err.Error())
		os.Exit(errCode)
	}
}

// splitLine 截取当前进程的信息，并填充结构体
func splitLine(s *string) (TaskInfo, error) {
	var st []string
	for _, v := range strings.Split(*s, "], ") {
		for id, vs := range strings.Split(v, " [") {
			vslSplit := strings.Split(vs, ": ")
			if id%2 == 0 {
				for _, vsl := range strings.Split(vslSplit[1], "/") {
					st = append(st, vsl)
				}
			} else {
				st = append(st, vs)
			}

		}
	}

	return TaskInfo{
		Pid:          st[0],
		Tid:          st[1],
		FullTaskPath: st[2],
		Uid:          st[3],
		Gid:          st[4],
	}, nil

}

// connetInfoSplit 针对请求的截取源目地址按" <-> "进行截取
func connetInfoSplit(s *string) (ConnectInfo, error) {

	var strogeValue []string
	if !strings.Contains(*s, " <-> ") {
		return ConnectInfo{}, fmt.Errorf("源目地址按\" <-> \"进行截取失败，传入的字符串为%s", *s)
	} else {
		for _, v := range strings.Split(*s, " <-> ") {
			strogeValue = append(strogeValue, v)
		}
		return ConnectInfo{
			Saddr: strings.Replace(strogeValue[0], "addr: ", "", -1),
			Daddr: strogeValue[1],
		}, nil
	}
}

// containerInfoSplit 针对请求的字符串按": "进行截取
func containerInfoSplit(s string) string {
	var st []string
	for _, v := range strings.Split(s, ": ") {
		st = append(st, v)
	}
	return st[1]
}

// PrintResultJson 解析结果，并输出json
func PrintResultJson(s Ebpfflow) []byte {
	// 字典格式化为json
	//data, err := json.Marshal(s)
	//if err != nil {
	//	fmt.Printf("JSON marshaling failed: %s", err)
	//	return nil
	//}

	// 针对json增加人类的可读性
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	return data
}
