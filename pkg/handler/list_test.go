package handler

import (
	"bytes"
	todo "do-app"
	"do-app/pkg/service"
	mock_service "do-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoLists, list todo.TodoList, userId int)

	testTable := []struct {
		name              string
		inputBody         string
		inputList         todo.TodoList
		userId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"test", "description":"testst"}`,
			inputList: todo.TodoList{
				Title:       "test",
				Description: "testst",
			},
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, list todo.TodoList, userId int) {
				s.EXPECT().Create(userId, list).Return(1, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"id":1}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoLists(c)
			testCase.mockBehavior(list, testCase.inputList, testCase.userId)

			services := &service.Service{TodoLists: list}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/", func(ctx *gin.Context) {
				ctx.Set(userCtx, testCase.userId)
			}, handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}
