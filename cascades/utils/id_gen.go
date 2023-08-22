package utils

type Generator[T any] interface {
	Generate() T
}
