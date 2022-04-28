package util

type Optional[T any] struct {
	val     T
	defined bool
}

func New[T any](val T) Optional[T] {
	return Optional[T]{val, true}
}

func (o Optional[A]) IsDefined() bool {
	return o.defined
}

func (o Optional[A]) Get() A {
	if !o.defined {
		panic("Optional undefined")
	}
	return o.val
}
