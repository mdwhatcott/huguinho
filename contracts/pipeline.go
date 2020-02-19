package contracts

import "errors"

type Handler interface {
	Handle(*Article)
}

type Finalizer interface {
	Finalize() error
}

var ErrDropArticle = errors.New("do not continue dispatching to handlers")
