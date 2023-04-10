package tasks

import (
	"time"

	"github.com/BimaAdi/ginGormBoilerplate/models"
	"github.com/BimaAdi/ginGormBoilerplate/repository"
	"github.com/BimaAdi/ginGormBoilerplate/settings"
)

func CreateSuperUser(envPath string, email string, username string, password string) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)

	// Initiate Database connection
	models.Initiate()

	now := time.Now()
	repository.CreateUser(models.DBConn, username, email, password, true, true, now, &now)
}
