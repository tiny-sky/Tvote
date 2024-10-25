package tools

import (
	"log"
	"net"
	"strconv"
	"strings"
)

func ErrToPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func FigureOutListen(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 1 {
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return ""
		}
		return ":" + listenOn
	}

	host, port, err := net.SplitHostPort(listenOn)
	if err != nil {
		log.Fatalf(err.Error())
		return listenOn

	}

	if len(host) > 0 && host != "0.0.0.0" {
		return listenOn
	}

	ip := InternalIp()
	if len(ip) == 0 {
		return listenOn
	}
	return strings.Join(append([]string{ip}, port), ":")
}

func InternalIp() string {
	infs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inf := range infs {
		if isEthDown(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

func isEthDown(f net.Flags) bool {
	return f&net.FlagUp != net.FlagUp
}

func isLoopback(f net.Flags) bool {
	return f&net.FlagLoopback == net.FlagLoopback
}
