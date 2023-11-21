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