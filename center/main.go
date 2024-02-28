package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, e := serial.OpenPort(c)
	if e != nil {
		log.Fatal(e)
	}

	cur := 0
	stack := arraystack.New()

	for {
		buf := make([]byte, 128)
		n, e := s.Read(buf)
		fmt.Print("\033[H\033[2J")
		if e != nil {
			log.Fatal(e)
		}

		for _, x := range string(buf[:n]) {
			switch x {
			case '0':
				cur++
				cur %= 10
			case '1':
				stack.Push(cur)
			case '2':
				if b, ok := stack.Pop(); ok {
					if a, ok := stack.Pop(); ok {
						sa := strconv.Itoa(a.(int))
						sb := strconv.Itoa(b.(int))
						s, e := strconv.Atoi(sa + sb)
						if e != nil {
							log.Fatal(e)
						} else {
							stack.Push(s)
						}
					}
				}
			case '3':
				if b, ok := stack.Pop(); ok {
					if a, ok := stack.Pop(); ok {
						stack.Push(a.(int) + b.(int))
					}
				}
			case '4':
				if a, ok := stack.Pop(); ok {
					stack.Push(-a.(int))
				}
			default:
			}
		}

		fmt.Println("CUR:", cur)
		fmt.Println("STACK:")
		it := stack.Iterator()
		for it.Next() {
			fmt.Println(it.Value())
		}
	}
}
