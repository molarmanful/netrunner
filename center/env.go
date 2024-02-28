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
		env.stack.Push(env.cur)

	case '2':
		if env.stack.Size() < 2 {
			log.Println("need 2 items")
		} else {
			b, _ := env.stack.Pop()
			a, _ := env.stack.Pop()
			c := new(big.Int)
			c.SetString(a.(*big.Int).Text(10)+b.(*big.Int).Text(10), 10)
			env.stack.Push(c)
		}

	case '3':
		if b, ok := env.stack.Pop(); ok {
			if a, ok := env.stack.Pop(); ok {
				env.stack.Push(a.(int) + b.(int))
			}
		}

	case '4':
		if a, ok := env.stack.Pop(); ok {
			env.stack.Push(-a.(int))
		}

	default:
	}
}

func (env *Env) kint(x rune) {
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
		ch := make([]byte, 1)
		os.Stdin.Read(ch)
		a, _ := env.stack.Pop()
		env.vars[ch[0]] = a

	case ':':
		ch := make([]byte, 1)
		os.Stdin.Read(ch)
		if a, ok := env.vars[ch[0]]; ok {
			env.stack.Push(a)
		} else {
			log.Println("undef var -", string(ch[0]))
		}

	case ',':
		if env.mode != 1 {
			env.mode = 1
			ch := make([]byte, 1)
			os.Stdin.Read(ch)
			env.macro_cur.name = ch[0]
			env.macro_cur.macro = ""
		} else {
			env.mode = 0
			env.macros[env.macro_cur.name] = env.macro_cur.macro[:len(env.macro_cur.macro)-1]
		}

	case '.':
		ch := make([]byte, 1)
		os.Stdin.Read(ch)
		if m, ok := env.macros[ch[0]]; ok {
			for _, x := range m {
				env.kint(x)
			}
		} else {
			log.Println("undef macro -", string(ch[0]))
		}

	default:
		log.Println("undef key -", x)
	}
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
