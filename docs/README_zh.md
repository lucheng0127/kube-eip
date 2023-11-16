# kube-eip

**Kube-eip**用于给kubevirt vmi提供业务网eip功能。

![网络架构图](./architecture/architecture.png)

kube-eip通过iptables做snat和dnat完成vmi ip和eip的转换。ipset定义pod cidr和service ip cidr，在策略路由和iptables中，对访问pod和service ip的流量不做snat，并且继续走cni0。

eip路由通过bgp向eip router同步eip ipv4路由，nexthop指向hyper的管理网地址。

## eip的绑定与解绑

自定义EipBinding CRD。一个EipBinding就代表一个eip绑定给一个vmi。通过operator监听，EipBinding和kubevirt的VirtualMachineInstance的变化，从而实现对eip绑定到vmi后端的动态变化。

## 组件

Kube-eip由Operator和EipAgent组成。Operator用于监听EipBinding，当有新增，更新，删除EipBinding事件时，通过grpc调用EipAgent实现eip绑定，更新，及解绑。