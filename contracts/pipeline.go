package contracts

import "errors"

type Handler interface {
	Handle(*Article)
}

type Finalizer interface {
	Finalize() error
}

var ErrDroppedArticle = errors.New("dropped article")
