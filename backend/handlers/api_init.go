package handlers

import (
	"backend/src/services"
)

type Handle struct {
	UserServices  *services.UserServices
	JointServices *services.JointServices
}

func NewHandle(root string) *Handle {
	services.ImageRoot = root
	us := services.NewUserServices()
	js := services.NewJointServices()
	services.CheckServicesStart(us, js)

	return &Handle{UserServices: us, JointServices: js}
}

func (handle *Handle) SetupAPI() {
	handle.JointServices.ConcurrentJobScanner()
}
