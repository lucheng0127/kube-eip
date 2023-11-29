package arp

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
)

type ArpCracker struct {
	handle  *pcap.Handle
	targets map[string]int
	iface   *net.Interface
}

func NewArpCracker(dev string) (*ArpCracker, error) {
	ctx := ctx.NewTraceContext()
	cracker := new(ArpCracker)
	cracker.targets = make(map[string]int)

	iface, err := net.InterfaceByName(dev)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("get interface %s %s", dev, err.Error()))
		return nil, err
	}
	cracker.iface = iface

	handle, err := pcap.OpenLive(dev, 1600, true, pcap.BlockForever)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("monitor interface %s traffic %s", dev, err.Error()))
		return nil, err
	} else if err := handle.SetBPFFilter("arp"); err != nil {
		logger.Error(ctx, fmt.Sprintf("set bpf filter %s", err.Error()))
		return nil, err
	}

	cracker.handle = handle

	return cracker, nil
}

func (cracker *ArpCracker) doArpReply(dstHwAddr, dstProtAddr, srcprotAdPr []byte) {
	ctx := ctx.NewTraceContext()
	ether := layers.Ethernet{
		EthernetType: layers.EthernetTypeARP,
		SrcMAC:       cracker.iface.HardwareAddr,
		DstMAC:       dstHwAddr,
	}
	arpReply := layers.ARP{
		AddrType:  layers.LinkTypeEthernet,
		Protocol:  layers.EthernetTypeIPv4,
		Operation: layers.ARPReply,

		HwAddressSize:   6,
		ProtAddressSize: 4,

		SourceHwAddress:   cracker.iface.HardwareAddr,
		SourceProtAddress: srcprotAdPr,

		DstHwAddress:   dstHwAddr,
		DstProtAddress: dstProtAddr,
	}

	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}, &ether, &arpReply)

	if err != nil {
		logger.Error(ctx, fmt.Sprintf("build arp reply %s", err.Error()))
	}

	if err = cracker.handle.WritePacketData(buf.Bytes()); err != nil {
		logger.Error(ctx, fmt.Sprintf("send arp reply %s", err.Error()))
	}
}

func (cracker *ArpCracker) handleArp(pkt gopacket.Packet) {
	arpLayer := pkt.Layer(layers.LayerTypeARP)
	if arpLayer == nil {
		return
	}

	arpPkt := arpLayer.(*layers.ARP)
	if arpPkt.Operation == layers.ARPReply {
		return
	}

	target := net.IP(arpPkt.DstProtAddress).String()
	if _, ok := cracker.targets[target]; !ok {
		// Arp request dst not in eip targets, do nothing
		return
	}

	cracker.doArpReply(arpPkt.SourceHwAddress, arpPkt.SourceProtAddress, arpPkt.DstProtAddress)
}

func (cracker *ArpCracker) AddTarget(target string) {
	ctx := ctx.NewTraceContext()
	logger.Info(ctx, fmt.Sprintf("add %s to arp poisoning targets", target))

	cracker.targets[target] = 0
}

func (cracker *ArpCracker) DeleteTarget(target string) {
	ctx := ctx.NewTraceContext()
	logger.Info(ctx, fmt.Sprintf("remove %s from arp poisoning targets", target))

	delete(cracker.targets, target)
}

func (cracker *ArpCracker) Poisoning() {
	packetSource := gopacket.NewPacketSource(cracker.handle, cracker.handle.LinkType())
	for pkt := range packetSource.Packets() {
		cracker.handleArp(pkt)
	}
}
