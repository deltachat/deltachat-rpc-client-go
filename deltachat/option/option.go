package option

import (
	"encoding/json"
)

type Option[T any] struct {
	Value *T
}

func Some[T any](value T) Option[T] {
	return Option[T]{Value: &value}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (self Option[T]) IsSome() bool {
	return self.Value != nil
}

func (self Option[T]) IsNone() bool {
	return self.Value == nil
}

func (self Option[T]) Unwrap() T {
	return *self.Value
}

func (self Option[T]) UnwrapOr(or T) T {
	if self.IsNone() {
		return or
	}
	return *self.Value
}

func (self Option[T]) MarshalJSON() ([]byte, error) {
	val, err := json.Marshal(self.Value)
	return val, err
}

func (self *Option[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.Value)
}
