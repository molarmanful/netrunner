package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type Env struct {
	Cur       string                           `json:"cur"`
	Stack     *arraystack.Stack                `json:"stack"`
	Stacks    *KeyMap[rune, *arraystack.Stack] `json:"stacks"`
	Cmds      string                           `json:"cmds"`
	cmdsMu    sync.Mutex
	Mode      int                   `json:"mode"`
	Vars      *KeyMap[rune, any]    `json:"vars"`
	Macros    *KeyMap[rune, string] `json:"macros"`
	Macro_rec ModeMacro             `json:"macro_rec"`
}

func NewEnv() *Env {
	env := &Env{
		Cur:       "",
		Stack:     arraystack.New(),
		Stacks:    NewKeyMap[rune, *arraystack.Stack](),
		Cmds:      "",
		Mode:      0,
		Vars:      NewKeyMap[rune, any](),
		Macros:    NewKeyMap[rune, string](),
		Macro_rec: ModeMacro{0, ""},
	}
	env.Stacks.Set('0', env.Stack)
	return env
}

type ModeMacro struct {
	Name  rune   `json:"name"`
	Macro string `json:"macro"`
}

var _0 = big.NewInt(0)
var _1 = big.NewInt(1)
var _10 = big.NewInt(10)

func BigBoolNOT(n *big.Int) *big.Int {
	if n.Sign() == 0 {
		return _1
	}
	return _0
}

func (env *Env) Loop() {
	for {
		x := env.WaitCh()
		Clr()
		env.KInt(x)
		env.Show()
	}
}

func (env *Env) KInt(x rune) {
	fmt.Println("KEY:", strconv.QuoteRune(x))

	switch x {

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		env.Cur += string(x)

	case ' ':
		a := new(big.Int)
		a.SetString(env.Cur, 10)
		env.Stack.Push(a)
		env.Cur = ""

	case '&':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			c := new(big.Int)
			c.SetString(a.String()+b.String(), 10)
			env.Stack.Push(c)
		})

	case '\t':
		env.Arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.Cur = a.String()
		})

	case '_':
		env.Arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.Stack.Push(a.Neg(a))
		})

	case '+':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.Stack.Push(a.Add(a, b))
		})

	case '-':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.Stack.Push(a.Sub(a, b))
		})

	case '*':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.Stack.Push(a.Mul(a, b))
		})

	case '/':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			if b.Cmp(_0) != 0 {
				a.DivMod(a, b, b)
				env.Stack.Push(a)
				env.Stack.Push(b)
			} else {
				log.Println("div by zero")
			}
		})

	case '^':
		env.Arg(2, func(xs []any) {
			a := xs[1].(*big.Int)
			b := xs[0].(*big.Int)
			env.Stack.Push(a.Exp(a, b, nil))
		})

	case '!':
		env.Arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.Stack.Push(BigBoolNOT(a))
		})

	case '?':
		env.Arg(1, func(xs []any) {
			a := xs[0].(*big.Int)
			env.Stack.Push(big.NewInt(int64(a.Sign())))
		})

	case 10:
		env.Arg(1, func(xs []any) {
			a := xs[0]
			b := new(big.Int)
			b.Set(a.(*big.Int))
			env.Stack.Push(a)
			env.Stack.Push(b)
		})

	case 'o':
		env.Arg(2, func(xs []any) {
			a := xs[1]
			b := xs[0]
			c := new(big.Int)
			c.Set(a.(*big.Int))
			env.Stack.Push(a)
			env.Stack.Push(b)
			env.Stack.Push(c)
		})

	case 127:
		env.Arg(1, func(xs []any) {})

	case '\\':
		env.Arg(2, func(xs []any) {
			env.Stack.Push(xs[0])
			env.Stack.Push(xs[1])
		})

	case '@':
		env.Arg(3, func(xs []any) {
			env.Stack.Push(xs[1])
			env.Stack.Push(xs[0])
			env.Stack.Push(xs[2])
		})

	case 'c':
		env.Stack.Clear()

	case '=':
		x = env.WaitCh()
		a, _ := env.Stack.Pop()
		env.Vars.Set(x, a)

	case ':':
		x = env.WaitCh()
		if a, ok := env.Vars.Get(x); ok {
			env.Stack.Push(a)
		} else {
			log.Println("undef var -", string(x))
		}

	case ',':
		if env.Mode != 1 {
			env.Mode = 1
			x = env.WaitCh()
			env.Macro_rec.Name = x
			env.Macro_rec.Macro = ""
		} else {
			env.Mode = 0
			env.Macros.Set(env.Macro_rec.Name, env.Macro_rec.Macro[:len(env.Macro_rec.Macro)-1])
		}

	case '.':
		x = env.WaitCh()
		if m1, ok := env.Macros.Get(x); ok {
			env.cmdsMu.Lock()
			env.Cmds = m1 + env.Cmds
			env.cmdsMu.Unlock()
		} else {
			log.Println("undef macro -", string(x))
		}

	case '#':
		x = env.WaitCh()
		if m1, ok := env.Macros.Get(x); ok {
			env.Arg(1, func(xs []any) {
				n := xs[0].(*big.Int)
				for n.Cmp(_0) > 0 {
					env.cmdsMu.Lock()
					env.Cmds = m1 + env.Cmds
					env.cmdsMu.Unlock()
					n.Sub(n, _1)
				}
			})
		} else {
			log.Println("undef macro -", string(x))
		}

	case '[':
		x = env.WaitCh()
		env.Stacks.Set(x, arraystack.New())
		*env.Stacks.m[x] = *env.Stack
		*env.Stack = *arraystack.New()

	case ']':
		x = env.WaitCh()
		if stack, ok := env.Stacks.Get(x); ok {
			it := stack.Iterator()
			for it.End(); it.Prev(); {
				b := new(big.Int)
				b.Set(it.Value().(*big.Int))
				env.Stack.Push(b)
			}
		} else {
			log.Println("undef stack -", string(x))
		}

	case 'w':
		x = env.WaitCh()
		if f, err := os.Create("env_" + fmt.Sprint(int(x)) + ".json"); err == nil {
			encoder := json.NewEncoder(f)
			if err := encoder.Encode(env); err != nil {
				log.Println(err)
			}
			defer f.Close()
		} else {
			log.Println(err)
		}

	case 'r':
		x = env.WaitCh()
		if f, err := os.Open("env_" + fmt.Sprint(int(x)) + ".json"); err == nil {
			decode := json.NewDecoder(f)
			if err := decode.Decode(&env); err != nil {
				log.Println(err)
			}
			defer f.Close()
		} else {
			log.Println(err)
		}

	default:
		log.Println("undef key -", x)
	}
}

func (env *Env) WaitCh() rune {
	for env.Cmds == "" {
	}
	println(env.Cmds)
	c := rune(env.Cmds[0])
	env.cmdsMu.Lock()
	env.Cmds = env.Cmds[1:]
	env.cmdsMu.Unlock()
	if env.Mode == 1 {
		env.Macro_rec.Macro += string(c)
	}
	return c
}

type FnArg func([]any)

func (env *Env) Arg(n int, f FnArg) {
	if env.Stack.Size() < n {
		log.Println("need", n, "items")
	} else {
		xs := make([]any, n)
		i := 0
		for i < n {
			a, _ := env.Stack.Pop()
			xs[i] = a
			i++
		}
		f(xs)
	}
}

func (env *Env) Show() {

	fmt.Println("CMDS:", strconv.Quote(env.Cmds))

	fmt.Println("MODE:", func() string {
		switch env.Mode {
		case 1:
			return "MACRO"
		default:
			return "NORMAL"
		}
	}())

	switch env.Mode {
	case 1:
		fmt.Println("\nRECORDING", string(env.Macro_rec.Name), ":", strconv.Quote(env.Macro_rec.Macro))
	}

	fmt.Println("\nVARS:")
	env.Vars.Each(func(v any, k rune) {
		fmt.Println(string(k), ":=", v)
	})

	fmt.Println("\nMACROS:")
	env.Macros.Each(func(v string, k rune) {
		fmt.Println(string(k), ":=", strconv.Quote(v))
	})

	fmt.Println("\nCUR:", env.Cur)

	env.Stacks.Each(func(v *arraystack.Stack, k rune) {
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
