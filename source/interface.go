package source

type SourceInterface[E any] interface {
	ReadMessage(max int) <-chan E
}
