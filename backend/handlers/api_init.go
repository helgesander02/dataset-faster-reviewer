package handlers

import (
	"backend/src/services"
)

type Handle struct {
	UserServices  *services.UserServices
	JointServices *services.JointServices
}

func NewHandle(root string, backupDir string) *Handle {
	services.ImageRoot = root
	services.BackupDir = backupDir
	us := services.NewUserServices()
	js := services.NewJointServices()

	services.CheckServicesState(us, js)
	return &Handle{UserServices: us, JointServices: js}
}
