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
	stacks    *KeyMap[rune, *arraystack.Stack]
	cmds      string
	cmds_mu   sync.Mutex
	mode      int
	vars      *KeyMap[rune, any]
	macros    *KeyMap[rune, string]
	macro_rec ModeMacro
}

func NewEnv() *Env {
	env := &Env{
		stack:     arraystack.New(),
		stacks:    NewKeyMap[rune, *arraystack.Stack](),
		cmds:      "",
		mode:      0,
		vars:      NewKeyMap[rune, any](),
		macros:    NewKeyMap[rune, string](),
		macro_rec: ModeMacro{0, ""},
	}
	env.stacks.Set('0', env.stack)
	return env
}

type ModeMacro struct {
	name  rune
	macro string
}

var _0 = big.NewInt(0)
var _1 = big.NewInt(1)
var _10 = big.NewInt(10)

func (env *Env) Loop() {
	for {
		x := env.waitch()
		Clr()
		env.KInt(x)
		env.Show()
	}
}

func (env *Env) KInt(x rune) {
	fmt.Println("KEY:", string(x))

	switch x {

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		env.stack.Push(big.NewInt(int64(x - 48)))

	case ' ':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			c := new(big.Int)
			c.SetString(a.Text(10)+b.Text(10), 10)
			env.stack.Push(c)
		})

	case '_':
		env.Arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.stack.Push(a.Neg(a))
		})

	case '+':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Add(a, b))
		})

	case '-':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Sub(a, b))
		})

	case '*':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Mul(a, b))
		})

	case '/':
		env.Arg(2, func(xs []any) {
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
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.stack.Push(a.Exp(a, b, nil))
		})

	case 10:
		env.Arg(1, func(xs []any) {
			a := xs[0]
			b := new(big.Int)
			b.Set(a.(*big.Int))
			env.stack.Push(a)
			env.stack.Push(b)
		})

	case 127:
		env.stack.Pop()

	case '\\':
		env.Arg(2, func(xs []any) {
			env.stack.Push(xs[0])
			env.stack.Push(xs[1])
		})

	case '@':
		env.Arg(3, func(xs []any) {
			env.stack.Push(xs[1])
			env.stack.Push(xs[0])
			env.stack.Push(xs[2])
		})

	case 'c':
		env.stack.Clear()

	case '=':
		x = env.waitch()
		a, _ := env.stack.Pop()
		env.vars.Set(x, a)

	case ':':
		x = env.waitch()
		if a, ok := env.vars.Get(x); ok {
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
			env.macros.Set(env.macro_rec.name, env.macro_rec.macro[:len(env.macro_rec.macro)-1])
		}

	case '.':
		x = env.waitch()
		if m1, ok := env.macros.Get(x); ok {
			env.cmds_mu.Lock()
			env.cmds = m1 + env.cmds
			env.cmds_mu.Unlock()
		} else {
			log.Println("undef macro -", string(x))
		}

	case '#':
		x = env.waitch()
		if m1, ok := env.macros.Get(x); ok {
			env.Arg(1, func(xs []any) {
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

	case '[':
		x = env.waitch()
		env.stacks.Set(x, arraystack.New())
		*env.stacks.m[x] = *env.stack
		*env.stack = *arraystack.New()

	case ']':
		x = env.waitch()
		if stack, ok := env.stacks.Get(x); ok {
			it := stack.Iterator()
			for it.End(); it.Prev(); {
				b := new(big.Int)
				b.Set(it.Value().(*big.Int))
				env.stack.Push(b)
			}
		} else {
			log.Println("undef stack -", string(x))
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

type FnArg func([]any)

func (env *Env) Arg(n int, f FnArg) {
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

func (env *Env) Show() {

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
		fmt.Println("\nRECORDING", string(env.macro_rec.name), ":", strconv.Quote(env.macro_rec.macro))
	}

	fmt.Println("\nVARS:")
	env.vars.Each(func(v any, k rune) {
		fmt.Println(string(k), ":=", v)
	})

	fmt.Println("\nMACROS:")
	env.macros.Each(func(v string, k rune) {
		fmt.Println(string(k), ":=", strconv.Quote(v))
	})

	env.stacks.Each(func(v *arraystack.Stack, k rune) {
		fmt.Println("\nSTACK", string(k), ":")
		it := v.Iterator()
		for it.End(); it.Prev(); {
			fmt.Print(it.Value(), " ")
		}
		fmt.Println("")
	})
}

func Clr() {
	fmt.Print("\033[H\033[2J")
}
