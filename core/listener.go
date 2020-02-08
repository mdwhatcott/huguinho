package core

import (
	"errors"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Listener struct {
	input   chan contracts.Article
	output  chan contracts.Article
	handler contracts.Handler
}

func NewListener(input, output chan contracts.Article, handler contracts.Handler) *Listener {
	return &Listener{
		input:   input,
		output:  output,
		handler: handler,
	}
}

func (this *Listener) Listen() error {
	defer close(this.output)
	defer this.drain()

	for article := range this.input {
		err := this.handler.Handle(&article)
		if err == ErrDropArticle {
			continue
		}
		if err != nil {
			return err
		}
		this.output <- article
	}

	finalizer, ok := this.handler.(contracts.Finalizer)
	if ok {
		return finalizer.Finalize()
	}

	return nil
}

func (this *Listener) drain() {
	for range this.input {
	}
}

var ErrDropArticle = errors.New("do not continue dispatching to handlers")
