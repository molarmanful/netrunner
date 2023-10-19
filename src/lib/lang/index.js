export default class {
  constructor() {
    this.stack = []

    this.cmds = ['+', '-', '*', '/', '%', '**']
  }

  f(f) {
    let a = f.length
    this.stack.push(f(...this.stack.splice(-a, a)))
  }

  ['+']() {
    return this.f((a, b) => a + b)
  }

  ['-']() {
    return this.f((a, b) => a - b)
  }

  ['*']() {
    return this.f((a, b) => a * b)
  }

  ['/']() {
    return this.f((a, b) => a / b)
  }

  ['%']() {
    return this.f((a, b) => a % b)
  }

  ['**']() {
    return this.f((a, b) => a ** b)
  }
}
