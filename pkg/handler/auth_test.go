package handler

import (
	"bytes"
	todo "do-app"
	"do-app/pkg/service"
	mock_service "do-app/pkg/service/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_SingUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name              string
		inputBody         string
		inputUser         todo.User
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"id":1}`,
		},
		{
			name:              "No Pole",
			inputBody:         `{"username":"test","password":"qwerty"}`,
			mockBehavior:      func(s *mock_service.MockAuthorization, user todo.User) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user signInInput)

	testTable := []struct {
		name              string
		inputBody         string
		inputUser         signInInput
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"qwerty"}`,
			inputUser: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user signInInput) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("1", nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"token":"1"}`,
		},
		{
			name:              "No Pole",
			inputBody:         `{"username":"test"}`,
			mockBehavior:      func(s *mock_service.MockAuthorization, user signInInput) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Failed in service",
			inputBody: `{"username":"test", "password":"qwerty"}`,
			inputUser: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user signInInput) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("1", errors.New("error generate token"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"error generate token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}
