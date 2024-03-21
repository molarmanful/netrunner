package main

type KeyMap[A comparable, B any] struct {
	m map[A]B
	k []A
}

func NewKeyMap[A comparable, B any]() *KeyMap[A, B] {
	return &KeyMap[A, B]{map[A]B{}, make([]A, 0)}
}

func (km *KeyMap[A, B]) Get(k A) (B, bool) {
	v, ok := km.m[k]
	return v, ok
}

func (km *KeyMap[A, B]) Set(k A, v B) {
	if _, ok := km.m[k]; !ok {
		km.k = append(km.k, k)
	}
	km.m[k] = v
}

type EachFn[A comparable, B any] func(B, A)

func (km *KeyMap[A, B]) Each(f EachFn[A, B]) {
	for _, k := range km.k {
		f(km.m[k], k)
	}
}
