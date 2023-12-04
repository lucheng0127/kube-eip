# kube-eip

**Kube-eip**用于给kubevirt vmi提供业务网eip功能。

![网络架构图](./architecture/architecture.png)

kube-eip通过iptables做snat和dnat完成vmi ip和eip的转换。

ipset定义pod，service ip cidr和已经做了eip binding的eip

和vmi ip在策略路由和iptables中，通过新增KUBE-EIP-PREROUTING，

KUBE-EIP-POSTROUTING链对ingress和egress流量进行劫持，并

进行nat转换。对于访问k8s内部的流量(访问pod ip和service ip)

不做nat转换。eip路由通过bgp向eip router同步eip ipv4路由，

nexthop指向hyper的管理网地址。可开启arp-poisoning，来做arp

代答（不建议开启），eip_agent如果down则无法正常应答eip的arp请求。

## eip的绑定与解绑

自定义EipBinding CRD。一个EipBinding就代表一个eip绑定给一个vmi。

通过operator监听，EipBinding和kubevirt的VirtualMachineInstance

的变化，如果eip绑定的后端vmi的ip发生变化，则rpc调用eip_agent的grpc

接口来更新后端的eip绑定规则。

## 组件

Kube-eip由Operator和EipAgent组成。Operator用于监听EipBinding，

当有新增，更新，删除EipBinding事件时，通过grpc调用EipAgent实现eip

绑定，更新，及解绑。

EipAgent daemonset形式运行在每个节点，暴露6127端口对外提供绑定和

解绑eip的接口。

## 使用

稳定版的eipbinding operator和eip agent镜像

* quay.io/shawnlu0127/eipbinding_operator:20231130
* quay.io/shawnlu0127/eip_agent:20231204

根据实际情况修改eip-agent-cm configmap中的配置(conf/agent/eip_agent.yaml)

```
# 部署
IMG=quay.io/shawnlu0127/eipbinding_operator:20231130 make deploy
make deploy-agent

# 清理
make undeploy
make undeploy-agent
```

**编译并部署你自己的镜像版本**

```
# Build eipbinding operator and eip agent
IMG=quay.io/shawnlu0127/eipbinding_operator:20231130 make docker-build-operator
IMG=quay.io/shawnlu0127/eip_agent:20231204 make docker-build-agent

# Push your image
IMG=quay.io/shawnlu0127/eipbinding_operator:20231130 make docker-push
IMG=quay.io/shawnlu0127/eip_agent:20231204 make docker-push

# Deploy eipbinding operator and eip agent
IMG=quay.io/shawnlu0127/eipbinding_operator:20231130 make deploy
make deploy-agent

# Undeploy
make undeploy
make undeploy-agent
```

**eipctl的使用**

eipctl是一个命令行工具用于绑定或者解绑eip给vmi

```
root@shawn-server:~/workspace/kube-eip# eipctl -h
NAME:
   eipctl - A new cli application

USAGE:
   eipctl [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --target value  rpc server address, default 127.0.0.1:6127 (default: "127.0.0.1:6127")
   --eip-ip value  eip ip address
   --vmi-ip value  vmi ip address
   --action value  action, bind or unbind
   --help, -h      show help
```

**eip_agent**

eip_agent以daemonset形式运行在每个节点，下面是eip_agent的帮助信息

```
root@shawn-server:~/workspace/kube-eip# eip_agent -h
NAME:
   EipAgent - A new cli application

USAGE:
   EipAgent [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value                                   agent port that rpc listen on (default: 6127)
   --log-level value                              log level, default info (default: "info")
   --internal-net value [ --internal-net value ]  networks that exclude from nat
   --gateway-ip value                             externel network gateway ip address
   --gateway-dev value                            externel network gateway device
   --bgp-type value                               bgp manager type, default is none, gobgp is avaliable (default: "none")
   --eip-cidr value                               eip network cidr
   --arp-poisoning                                whether use arp poisoning to make a arp reply for eip, default is false, when enable will not add eip to external network gateway device (default: false)
   --help, -h                                     show help
```

**开始使用吧**

```
root@shawn-server:~/workspace/kube-eip# kubectl get all -n kube-eip-system
Warning: kubevirt.io/v1 VirtualMachineInstancePresets is now deprecated and will be removed in v2.
NAME                                               READY   STATUS    RESTARTS   AGE
pod/kube-eip-controller-manager-5df4d7d5fd-8qtr6   2/2     Running   0          2m11s

NAME                                                  TYPE        CLUSTER-IP        EXTERNAL-IP   PORT(S)    AGE
service/kube-eip-controller-manager-metrics-service   ClusterIP   192.168.223.251   <none>        8443/TCP   2m11s

NAME                                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/kube-eip-controller-manager   1/1     1            1           2m11s

NAME                                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/kube-eip-controller-manager-5df4d7d5fd   1         1         1       2m11s
root@shawn-server:~/workspace/kube-eip# kubectl get all -n kube-eip-agent
Warning: kubevirt.io/v1 VirtualMachineInstancePresets is now deprecated and will be removed in v2.
NAME                  READY   STATUS    RESTARTS   AGE
pod/eip-agent-2dkbb   1/1     Running   0          2m34s
pod/eip-agent-klcpd   1/1     Running   0          2m34s

NAME                       DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
daemonset.apps/eip-agent   2         2         2       2            2           <none>          2m34s
root@shawn-server:~/workspace/kube-eip# kubectl get eipbinding
NAME     AGE
cirros   45s
```

