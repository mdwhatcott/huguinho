package core

import "github.com/mdwhatcott/huguinho/contracts"

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
		if err != nil {
			return err
		}
		this.output <- article
	}

	return nil
}

func (this *Listener) drain() {
	for range this.input {
	}
}
