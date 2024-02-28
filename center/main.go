package main

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"

	"github.com/emirpasic/gods/stacks/arraystack"
)

func main() {
	// c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	// s, e := serial.OpenPort(c)
	// if e != nil {
	// 	log.Fatal(e)
	// }

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	env := Env{big.NewInt(0), arraystack.New(), map[byte]any{}}

	fmt.Print("\033[H\033[2J")
	env.show()

	// go func() {
	// 	for {
	// 		buf := make([]byte, 128)
	// 		n, e := s.Read(buf)
	// 		fmt.Print("\033[H\033[2J")
	// 		if e != nil {
	// 			log.Fatal(e)
	// 		}
	//
	// 		for _, x := range string(buf[:n]) {
	// 			env.kext(x)
	// 		}
	//
	// 		env.show()
	// 	}
	// }()

	ch := make([]byte, 1)
	for {
		os.Stdin.Read(ch)
		fmt.Print("\033[H\033[2J")
		env.kint(ch[0])
		env.show()
	}
}
