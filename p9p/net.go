package p9p

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func parsedial(dial string) (network, laddr string, err error) {
	f := strings.Split(dial, "!")
	if len(f) < 2 {
		return "", "", errors.New("invalid dial string: " + dial)
	}
	if f[0] == "net" {
		f[0] = "tcp"
	}
	switch f[0] {
	case "unix":
		return "unix", f[1], nil
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		host := f[1]
		if len(f) > 2 {
			host = host + ":" + f[2]
		}
		return f[0], host, nil
	default:
		return "", "", errors.New("unsupported network " + f[0])
	}
}

/* Just like net.Dial but accepts a dial-string */
func Dial(dial string) (net.Conn, error) {
	network, laddr, err := parsedial(dial)
	if err != nil {
		return nil, err
	}
	return net.Dial(network, laddr)
}

/* Addr can be a dial-string or a service name */
func ListenSrv(addr string) (net.Listener, error) {
	if !strings.Contains(addr, "!") {
		addr = "unix!" + Namespace() + "/" + addr
	}
	network, laddr, err := parsedial(addr)
	if err != nil {
		return nil, err
	}
	return net.Listen(network, laddr)
}

/* In go, attempting to listen on a unix socket that already exists in the
 * filesystem fails with error "address already in use". The socket is
 * removed from the filesystem when the Listener is closed, but a
 * typical 9p server will run until it receives a fatal signal, bypassing
 * the usual cleanup code. Hence this convenience function exists to
 * ease the process of restarting servers. */
func CloseOnSignal(listener net.Listener) {
	sigs := make(chan os.Signal)
	go func() {
		ossig := <- sigs
		syssig, _ := ossig.(syscall.Signal)
		listener.Close()
		os.Exit(128 + int(syssig))
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}
