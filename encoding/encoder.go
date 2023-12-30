package encoding

type Decoder[T any] interface {
	Decode(T) error
}
