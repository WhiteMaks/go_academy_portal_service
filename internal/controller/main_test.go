package controller

import (
	"os"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if err := util.RegisterBindingValidators(); err != nil {
		panic(err)
	}
	if err := model.RegisterAthleteBindingValidatorsFromGin(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
