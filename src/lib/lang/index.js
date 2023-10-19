import { writable } from 'svelte/store'

export default class Env {
  constructor() {
    this.stack = writable([])

    this.cmds = {
      _: () => this.f(a => -a),
      '+': () => this.f((a, b) => a + b),
      '-': () => this.f((a, b) => a - b),
      '*': () => this.f((a, b) => a * b),
      '/': () => this.f((a, b) => a / b),
      '%': () => this.f((a, b) => a % b),
      '**': () => this.f((a, b) => a ** b),
      dup: () => (console.log(this.stack$), this.push(this.at(-1))),
      pop: () => this.pop(),
    }
  }

  static from(s) {
    let env = new Env()
    env.stack.set(s)
    return env
  }

  f(f) {
    let n = f.length
    this.checkLen(n)
    this.push(f(...this.last(n)))
  }

  push(a) {
    this.stack.update(xs => [...xs, a])
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
    let a
    this.stack.subscribe(xs => (a = xs))()
    return a
  }

  last(n) {
    let a
    this.stack.update(xs => ((a = xs.splice(-n, n)), xs))
    return a
  }
}
