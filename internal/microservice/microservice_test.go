package microservice

import (
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/util"
)

func TestNewMicroserviceRegistersAuthRoutes(t *testing.T) {
	ms := NewMicroservice(util.Microservice{
		Port:          8080,
		TokenKey:      "test-token-key",
		TokenLifeTime: time.Minute,
	}, nil)

	registeredRoutes := make(map[string]bool)
	for _, route := range ms.Router.Routes() {
		registeredRoutes[route.Method+" "+route.Path] = true
	}

	for _, route := range []string{
		"GET /api/microservice/v1",
		"POST /api/user/token/v1",
		"POST /api/user/admin/v1",
		"GET /api/user/v1",
		"POST /api/user/v1",
	} {
		if !registeredRoutes[route] {
			t.Errorf("expected route %q to be registered", route)
		}
	}
}
