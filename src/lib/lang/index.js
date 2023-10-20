import { writable } from 'svelte/store'

export default class Env {
  constructor() {
    this.stack = writable([])
    this.code = writable([])
    this.err = writable('')

    this.cmds = {
      '.': () => this.f((a, b) => +(a + '.' + b)),
      _: () => this.f(a => -a),
      '+': () => this.f((a, b) => a - -b),
      '++': () => this.f((a, b) => '' + a + b),
      '-': () => this.f((a, b) => a - b),
      '*': () => this.f((a, b) => a * b),
      '/': () => this.f((a, b) => a / b),
      '%': () => this.f((a, b) => a % b),
      '**': () => this.f((a, b) => a ** b),
      dup: () => this.fs(a => [a, a]),
      pop: () => this.pop(),
      swap: () => this.fs((a, b) => [b, a]),
      rot: () => this.fs((a, b, c) => [b, c, a]),
      rot_: () => this.fs((a, b, c) => [c, a, b]),
    }
  }

  step(f) {
    this.err.set('')
    try {
      f()
    } catch (e) {
      this.err.set(e.message)
    }
  }

  static from(s) {
    let env = new Env()
    env.stack.set(s)
    return env
  }

  showStack() {
    return 'STACK:\n\n' + this.stack$.map(x => JSON.stringify(x)).join`\n`
  }

  f(f) {
    let n = f.length
    this.checkLen(n)
    this.push(f(...this.last(n)))
  }

  fs(f) {
    let n = f.length
    this.checkLen(n)
    this.push(...f(...this.last(n)))
  }

  push(...a) {
    this.stack.update(xs => xs.concat(a))
  }

  pop() {
    let a
    this.checkLen(1)
    this.stack.update(xs => ((a = xs.pop()), xs))
    return a
  }

  at(i) {
    this.checkLen(i)
    return this.stack$.at(i)
  }

  checkLen(n) {
    let l = this.len
    if (l < n) throw new Error(`stack len ${l} < ${n}`)
  }

  get len() {
    return this.stack$.length
  }

  get stack$() {
    return Env.getSub(this.stack)
  }

  get code$() {
    return Env.getSub(this.code)
  }

  static getSub(s) {
    let a
    s.subscribe(xs => {
      a = xs
    })()
    return a
  }

  last(n) {
    let a
    this.stack.update(xs => ((a = xs.splice(-n, n)), xs))
    return a
  }
}
