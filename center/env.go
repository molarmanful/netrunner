package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type Env struct {
	cur       *big.Int
	stack     *arraystack.Stack
	mode      int
	vars      map[byte]any
	macros    map[byte]string
	macro_cur ModeMacro
}

type ModeMacro struct {
	name  byte
	macro string
}

var _1 = big.NewInt(1)
var _10 = big.NewInt(10)

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

func (env *Env) kint(x rune, m string) {
tco:
	fmt.Println("KEY:", string(x))

	// TODO: make sure ch inputs are captured in macros
	switch env.mode {
	case 1:
		env.macro_cur.macro += string(x)
	}

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
			a.DivMod(a, b, b)
			env.stack.Push(a)
			env.stack.Push(b)
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
		c, s := env.readch(m)
		m = s
		a, _ := env.stack.Pop()
		env.vars[c] = a

	case ':':
		c, s := env.readch(m)
		m = s
		if a, ok := env.vars[c]; ok {
			env.stack.Push(a)
		} else {
			log.Println("undef var -", string(c))
		}

	case ',':
		if env.mode != 1 {
			env.mode = 1
			c, s := env.readch(m)
			m = s
			env.macro_cur.name = c
			env.macro_cur.macro = ""
		} else {
			env.mode = 0
			env.macros[env.macro_cur.name] = env.macro_cur.macro[:len(env.macro_cur.macro)-1]
		}

	case '.':
		c, s := env.readch(m)
		m = s
		if m1, ok := env.macros[c]; ok {
			m = m1 + m
		} else {
			log.Println("undef macro -", string(c))
		}

	default:
		log.Println("undef key -", x)
	}

	if m == "" {
		return
	}
	x = rune(m[0])
	m = m[1:]
	goto tco
}

func (env *Env) readch(m string) (byte, string) {
	if m == "" {
		c := make([]byte, 1)
		os.Stdin.Read(c)
		if env.mode == 1 {
			env.macro_cur.macro += string(c[0])
		}
		return c[0], m
	}
	return m[0], m[1:]
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

	fmt.Println("CUR:", env.cur)

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
		fmt.Println("\nRECORDING ", string(env.macro_cur.name), ":", strconv.Quote(env.macro_cur.macro))
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
