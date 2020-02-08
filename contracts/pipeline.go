package contracts

type Listener interface {
	Listen() error
}

type Handler interface {
	Handle(*Article) error
}

type Finalizer interface {
	Finalize() error
}
