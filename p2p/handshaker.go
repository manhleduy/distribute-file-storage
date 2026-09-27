package p2p

import "errors"

//err when the handshake is between the local and remote could not be formed
var ErrInvalidHandshake = errors.New("invalid handshake")
// handshakeFunc is the signal for the connection is accepted or not
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
