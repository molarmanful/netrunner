package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type Env struct {
	cur       *big.Int
	stack     *arraystack.Stack
	cmds      string
	cmds_mu   sync.Mutex
	mode      int
	vars      map[rune]any
	macros    map[rune]string
	macro_rec ModeMacro
}

type ModeMacro struct {
	name  rune
	macro string
}

var _0 = big.NewInt(0)
var _1 = big.NewInt(1)
var _10 = big.NewInt(10)

func (env *Env) loop() {
	for {
		x := env.waitch()
		env.clr()
		env.kint(x)
		env.show()
	}
}

func (env *Env) kext(x rune) {
	switch x {

	case '0':
		env.cur.Add(env.cur, _1).Mod(env.cur, _10)

	case '1':
		a := new(big.Int)
		a.Set(env.cur)
		env.stack.Push(a)

	case '2':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			c := new(big.Int)
			c.SetString(a.Text(10)+b.Text(10), 10)
			env.stack.Push(c)
		})

	case '3':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Add(a, b))
		})

	case '4':
		env.arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.stack.Push(a.Neg(a))
		})

	default:
	}
}

func (env *Env) kint(x rune) {
	fmt.Println("KEY:", string(x))

	switch x {

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		env.stack.Push(big.NewInt(int64(x - 48)))

	case ' ':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			c := new(big.Int)
			c.SetString(a.Text(10)+b.Text(10), 10)
			env.stack.Push(c)
		})

	case '_':
		env.arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.stack.Push(a.Neg(a))
		})

	case '+':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Add(a, b))
		})

	case '-':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Sub(a, b))
		})

	case '*':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Mul(a, b))
		})

	case '/':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			if b.Cmp(_0) != 0 {
				a.DivMod(a, b, b)
				env.stack.Push(a)
				env.stack.Push(b)
			} else {
				log.Println("div by zero")
			}
		})

	case '^':
		env.arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Exp(a, b, nil))
		})

	case 10:
		env.arg(1, func(xs []any) {
			a := xs[0]
			b := new(big.Int)
			b.Set(a.(*big.Int))
			env.stack.Push(a)
			env.stack.Push(b)
		})

	case 127:
		env.stack.Pop()

	case '\\':
		env.arg(2, func(xs []any) {
			env.stack.Push(xs[0])
			env.stack.Push(xs[1])
		})

	case '@':
		env.arg(3, func(xs []any) {
			env.stack.Push(xs[1])
			env.stack.Push(xs[0])
			env.stack.Push(xs[2])
		})

	case '=':
		x = env.waitch()
		a, _ := env.stack.Pop()
		env.vars[x] = a

	case ':':
		x = env.waitch()
		if a, ok := env.vars[x]; ok {
			env.stack.Push(a)
		} else {
			log.Println("undef var -", string(x))
		}

	case ',':
		if env.mode != 1 {
			env.mode = 1
			x = env.waitch()
			env.macro_rec.name = x
			env.macro_rec.macro = ""
		} else {
			env.mode = 0
			env.macros[env.macro_rec.name] = env.macro_rec.macro[:len(env.macro_rec.macro)-1]
		}

	case '.':
		x = env.waitch()
		if m1, ok := env.macros[x]; ok {
			env.cmds_mu.Lock()
			env.cmds = m1 + env.cmds
			env.cmds_mu.Unlock()
		} else {
			log.Println("undef macro -", string(x))
		}

	case '#':
		x = env.waitch()
		if m1, ok := env.macros[x]; ok {
			env.arg(1, func(xs []any) {
				n := xs[0].(*big.Int)
				for n.Cmp(_0) > 0 {
					env.cmds_mu.Lock()
					env.cmds = m1 + env.cmds
					env.cmds_mu.Unlock()
					n.Sub(n, _1)
				}
			})
		} else {
			log.Println("undef macro -", string(x))
		}

	default:
		log.Println("undef key -", x)
	}
}

func (env *Env) waitch() rune {
	for env.cmds == "" {
	}
	c := rune(env.cmds[0])
	env.cmds_mu.Lock()
	env.cmds = env.cmds[1:]
	env.cmds_mu.Unlock()
	if env.mode == 1 {
		env.macro_rec.macro += string(c)
	}
	return c
}

type fn_arg func([]any)

func (env *Env) arg(n int, f fn_arg) {
	if env.stack.Size() < n {
		log.Println("need", n, "items")
	} else {
		xs := make([]any, n)
		i := 0
		for i < n {
			a, _ := env.stack.Pop()
			xs[i] = a
			i++
		}
		f(xs)
	}
}

func (env *Env) show() {

	fmt.Println("CMDS:", strconv.Quote(env.cmds))

	fmt.Println("MODE:", func() string {
		switch env.mode {
		case 1:
			return "MACRO"
		default:
			return "NORMAL"
		}
	}())

	switch env.mode {
	case 1:
		fmt.Println("\nRECORDING ", string(env.macro_rec.name), ":", strconv.Quote(env.macro_rec.macro))
	}

	fmt.Println("\nVARS:")
	for k, v := range env.vars {
		fmt.Println(string(k), ":=", v)
	}

	fmt.Println("\nMACROS:")
	for k, v := range env.macros {
		fmt.Println(string(k), ":=", strconv.Quote(v))
	}

	fmt.Println("\nSTACK:")
	it := env.stack.Iterator()
	for it.End(); it.Prev(); {
		fmt.Println(it.Value())
	}
}

func (env *Env) clr() {
	fmt.Print("\033[H\033[2J")
}
