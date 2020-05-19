package utility

import (
	gnq "github.com/AkihiroSuda/go-netfilter-queue"
)

//CheckProtocol return Protocol TCP:1 UDP:2 Others:0
func CheckProtocol(packet gnq.NFPacket) int {
	protocol := uint16(packet.Packet.Layers()[0].LayerContents()[9])
	switch protocol {
	case 6:
		return 1
	case 16:
		return 2
	default:
		return 0
	}
}
