package encoding

type Encoder[T any] interface {
	Encode(T) error
}
