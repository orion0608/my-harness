package lifecycle

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	defaultPortMin = 50000
	defaultPortMax = 60000
)

func pickPort(host string, portMin, portMax int) (int, error) {
	if portMin <= 0 {
		portMin = defaultPortMin
	}
	if portMax <= portMin {
		portMax = defaultPortMax
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	span := portMax - portMin + 1
	tries := span
	if tries > 200 {
		tries = 200
	}

	for i := 0; i < tries; i++ {
		port := portMin + rng.Intn(span)
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			continue
		}
		_ = ln.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no free port in range %d-%d on %s", portMin, portMax, host)
}
