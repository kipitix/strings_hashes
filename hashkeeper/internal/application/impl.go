package application

import (
	"hashkeeper/internal/domain/calculate"
	"hashkeeper/internal/domain/find"
)

type ShutdownCallback func()

type appImpl struct {
	calcHandler      calculate.CalculateHandler
	findHandler      find.FindHandler
	shutdownCallback ShutdownCallback
}

var _ App = (*appImpl)(nil)

func NewApp(calcHandler calculate.CalculateHandler, findHandler find.FindHandler, shutdownCallback ShutdownCallback) App {
	return &appImpl{
		calcHandler:      calcHandler,
		findHandler:      findHandler,
		shutdownCallback: shutdownCallback,
	}
}

func (a appImpl) CalculateHandler() calculate.CalculateHandler {
	return a.calcHandler
}

func (a appImpl) FindHandler() find.FindHandler {
	return a.findHandler
}

func (a appImpl) Shutdown() {
	a.shutdownCallback()
}
