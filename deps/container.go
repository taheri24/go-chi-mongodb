package deps

import "github.com/goava/di"

var MainContainer *di.Container

func SetDiContainer(container *di.Container) {
	MainContainer = container
}
