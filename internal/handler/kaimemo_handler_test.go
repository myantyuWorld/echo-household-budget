package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoRouterSetup(t *testing.T) {
	tests := []struct {
		name           string
		setupRouter    func() *echo.Echo
		expectedRoutes []string
	}{
		{
			name: "router with basic middleware",
			setupRouter: func() *echo.Echo {
				e := echo.New()
				return e
			},
			expectedRoutes: []string{"/", "/health"},
		},
		{
			name: "router with custom error handler",
			setupRouter: func() *echo.Echo {
				e := echo.New()
				e.HTTPErrorHandler = func(err error, c echo.Context) {
					c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
				return e
			},
			expectedRoutes: []string{"/", "/health"},
		},
		{
			name: "router with custom binder",
			setupRouter: func() *echo.Echo {
				e := echo.New()
				e.Binder = &echo.DefaultBinder{}
				return e
			},
			expectedRoutes: []string{"/", "/health"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setupRouter()

			for _, route := range tt.expectedRoutes {
				req := httptest.NewRequest(http.MethodGet, route, nil)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				assert.NotEqual(t, http.StatusNotFound, rec.Code, "Route %s should be registered", route)
			}

			nonExistentRoute := "/non-existent"
			req := httptest.NewRequest(http.MethodGet, nonExistentRoute, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code, "Non-existent route should return 404")
		})
	}
}
