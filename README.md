# kube-eip

[中文](./docs/README_zh.md)

**Kube-eip** is an elastic ip management add-on for kubevirt. The aim is to provide an access to kubevirt vmi through elastic ip.

## Architecture

![Architecture](./docs/architecture/architecture.png)

At it's core, kube-eip use the rules of iptables implement the transform between eip and vm instance ip. By default we access the pod and service network via vmi pod ip, others by eip, if make a eipBinding to a vmi pod. For convenience we use ipset to manage the networks of pods and service, and alse this's a policy route for each eipBinding, make sure that we connect to other pods with vmi pod ip.

Then for eip routes, we use bgp to declare that the nexthop to access vmi pod that binded eip. There ware a series bgp server we can choise, but we can use gobgp as the bgp library native.

## Lifecycle of eip

Kube-eip extends kubeernets by adding a eipBinding CRD. An eipBinding represent a eip binded to a vmi pod. And an eip can be create and bind or destoryed along with the eipBinding. Also operator will watch EipBinding and kubevirt VirtualMachineInstance, and handle create update and delete event.

## Modules

There are two compose Operator and EipAgent. Operator watch the EipBinding and VirtualMachineInstance create, update and delete event. Then call EipAgent to build the rules on hyper, via grpc.

## Usage

The stable eipbinding operator and eip agent docker images

* quay.io/shawnlu0127/eipbinding_operator:20231130
* quay.io/shawnlu0127/eip_agent:20231204

*TODO(user): modify configmap eip-agent-cm (config/agent/eip_agent.yaml)*

```
# Modify it before deploy eip_agent kube-eip/config/agent/eip_agent.yaml
...

# TODO(user): change content of config map
apiVersion: v1
kind: ConfigMap
metadata:
  name: eip-agent-cm
  namespace: kube-eip-agent
data:
  svc_net: 192.168.223.0/24      # K8s service ip cidr
  pod_net: 10.244.0.0/16         # K8s pod ip cidr
  eip_net: 192.168.18.0/24       # The public network cidr
  eip_gw_ip: 192.168.18.1        # The public network gateway
  eip_gw_dev: enp2s0             # The interface on each hyper, that access public netwrok interface(If interface name not same, add a linux bridge with the same name(such as br-pub) and add interface to linux bridge)
  log_level: debug

...
```

```
# Deploy
IMG={your own image name and tag} make deploy
make deploy-agent

# Undeploy
make undeploy
make undeploy-agent
```

**Build your own image, push and deploy**

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

**eipctl**

eipctl is an command line tool to bind or unbind eip to or from vmi

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

eip_agent run as daemonset in pod, there is the help info of it

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

**Have fun and enjoy it ٩(๑>◡<๑)۶**

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
