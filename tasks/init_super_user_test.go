package tasks_test

import (
	"testing"

	"github.com/BimaAdi/ginGormBoilerplate/core"
	"github.com/BimaAdi/ginGormBoilerplate/migrations"
	"github.com/BimaAdi/ginGormBoilerplate/models"
	"github.com/BimaAdi/ginGormBoilerplate/settings"
	"github.com/BimaAdi/ginGormBoilerplate/tasks"
	"github.com/stretchr/testify/assert"
)

func TestCreateSuperUser(t *testing.T) {
	// Given
	settings.InitiateSettings("../.env")
	models.Initiate()
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
	models.ClearAllData()

	// When
	tasks.CreateSuperUser("../.env", "test@local.com", "test", "password")

	// Expect
	createdUser := models.User{}
	err := models.DBConn.Where("email = ? AND username = ?", "test@local.com", "test").First(&createdUser).Error
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.True(t, core.CheckPasswordHash("password", createdUser.Password))
}
