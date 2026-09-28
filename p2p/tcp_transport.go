package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	//conn is the underlying connection of the peer
	conn net.Conn

	//if we dial and retrieve a conn  ->outbound == true
	//if we accept and retrieve a conn -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
//remoteAddr implements the Peer interface and will return the 
// remote address of the underlying connection.
func (p *TCPPeer) RemoteAddr() net.Addr{
	return p.conn.RemoteAddr()
}
// Close implement the peer interface
func (p *TCPPeer) Close() error{

	return p.conn.Close()
} 

func (p *TCPPeer) Send(b []byte) error{
	_, err := p.conn.Write(b)
	return err
}
type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer		  func(Peer) error

}
type TCPTransport struct {
	TCPTransportOpts
	listener 	net.Listener
	rpcch 		chan RPC
}
//close implements the Transport interface
func (t *TCPTransport) Close()error{

	return t.listener.Close()
}
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),

	}
}
//consume the inplementf the Transport interface, which will retrun read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <- chan RPC {

	return t.rpcch
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed){
			return 
		}

		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}
		

		go t.handleConn(conn,false)
	}

}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)


	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	peer := NewTCPPeer(conn, outbound)

	defer func(){
		fmt.Printf("dropping peer connection: %s", err )
		conn.Close()
	}()
	if err := t.HandshakeFunc(peer); err != nil {
		
		return
	}
	if t.OnPeer != nil{
		if err = t.OnPeer(peer); err!= nil{
			return 
		}
	}

	// Read Loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
	
		if err != nil {
			return 
		}
		rpc.From = conn.RemoteAddr()

		t.rpcch <- rpc
		
	}

}
//Dial implement the transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil{
		return err
	}
	
	go t.handleConn(conn, true)


	return nil
}

