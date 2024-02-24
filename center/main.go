package main

import (
	"log"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, e := serial.OpenPort(c)
	if e != nil {
		log.Fatal(e)
	}

	for {
		buf := make([]byte, 128)
		n, e := s.Read(buf)
		if e != nil {
			log.Fatal(e)
		}
		log.Printf("%q", buf[:n])
	}
}
