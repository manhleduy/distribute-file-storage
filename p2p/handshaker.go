package p2p

//handshakeFunc is the signal for the connection is accepted or not
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {return nil}

