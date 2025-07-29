package handlers

import (
	"backend/src/services"
)

type Handle struct {
	DM *services.DataManager
}

func NewHandle(root string) *Handle {
	dm := services.NewDataManager(root)
	dm.SetupServices()

	return &Handle{DM: dm}
}

func (handle *Handle) SetupAPI() {
	handle.DM.ConcurrentJobScanner()
}
