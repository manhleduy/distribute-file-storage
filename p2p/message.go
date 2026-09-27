package p2p

import "net"

//message present any arbitary data that is being sent
//over the each trasnport between two noades in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
