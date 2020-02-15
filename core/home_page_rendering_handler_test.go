package core

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestHomePageRenderingHandlerFixture(t *testing.T) {
	gunit.Run(new(HomePageRenderingHandlerFixture), t)
}

type HomePageRenderingHandlerFixture struct {
	*gunit.Fixture

	handler  *HomePageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *HomePageRenderingHandlerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewHomePageRenderingHandler(this.disk, this.renderer)
}

func (this *HomePageRenderingHandlerFixture) Test() {
}
