package handler

import (
	"do-app/pkg/service"
	mock_service "do-app/pkg/service/mocks"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectStatusCode    int
		expectResponsesBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectStatusCode:    200,
			expectResponsesBody: "1",
		},
		{
			name:                "Empty header",
			headerName:          "",
			headerValue:         "Bearer token",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectStatusCode:    401,
			expectResponsesBody: `{"message":"empty auth header"}`,
		},
		{
			name:                "Invalid header",
			headerName:          "Authorization",
			headerValue:         "Bearer",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectStatusCode:    401,
			expectResponsesBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Service failrule",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("invalid parse token"))
			},
			expectStatusCode:    401,
			expectResponsesBody: `{"message":"invalid parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(userCtx)
				ctx.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectResponsesBody, w.Body.String())
		})
	}
}
