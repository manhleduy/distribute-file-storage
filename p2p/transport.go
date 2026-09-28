package p2p
//peer is an interface which present the remote node
type Peer interface{
	Close() error

}
//transport is anything that handles the communication
// between the nodes in these network, this can be of the 
// form (TCP, UDP, Websockets)
type Transport interface{
	Dial(string) error
	ListenAndAccept() error
	Consume() <- chan RPC
	Close() error
}
