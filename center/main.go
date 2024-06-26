package main

import (
	"os"
	"os/exec"
)

func main() {
	// c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	// s, e := serial.OpenPort(c)
	// if e != nil {
	// 	log.Fatal(e)
	// }

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	env := NewEnv()
	Clr()
	env.Show()

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

	go func() {
		for {
			ch := make([]byte, 1)
			os.Stdin.Read(ch)
			env.cmdsMu.Lock()
			env.Cmds += string(ch[0])
			env.cmdsMu.Unlock()
		}
	}()

	env.Loop()
}
