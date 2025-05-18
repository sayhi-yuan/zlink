package source

type SourceInterface[E any] interface {
	ReadMessage() <-chan E
}
