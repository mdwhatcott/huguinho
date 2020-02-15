package core

import "github.com/mdwhatcott/huguinho/contracts"

type HomePageRenderingHandler struct {
}

func NewHomePageRenderingHandler(disk contracts.WriteFile, renderer contracts.Renderer) *HomePageRenderingHandler {
	return &HomePageRenderingHandler{}
}
