package contracts

type Listener interface {
	Listen() error
}

type Handler interface {
	Handle(*Article) error
}
