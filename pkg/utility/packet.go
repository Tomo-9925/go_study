package utility

import (
	gnq "github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket/layers"
)

//CheckSrcPort returns source port 
func CheckSrcPort(p gnq.NFPacket) int {
if tcpLayer := p.Packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
  // Get actual TCP data from this layer
	tcp, _ := tcpLayer.(*layers.TCP)
	return (int)(tcp.SrcPort)
}else if udpLayer :=p.Packet.Layer(layers.LayerTypeUDP); udpLayer !=nil {
	// Get actual UDP data from this layer
	udp,_:=udpLayer.(*layers.UDP)
	return (int)(udp.SrcPort)
}
return -1
}
