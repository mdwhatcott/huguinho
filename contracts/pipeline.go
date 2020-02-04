package contracts

type Handler interface {
	Listen()
	Finalize()
}
