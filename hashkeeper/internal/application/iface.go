package application

import (
	"hashkeeper/internal/domain/calculate"
	"hashkeeper/internal/domain/find"
)

type App interface {
	CalculateHandler() calculate.CalculateHandler
	FindHandler() find.FindHandler

	Shutdown()
}
