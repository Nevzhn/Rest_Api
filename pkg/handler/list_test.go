package handler

import (
	"bytes"
	todo "do-app"
	"do-app/pkg/service"
	mock_service "do-app/pkg/service/mocks"
	"fmt"
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
		{
			name:      "No User",
			inputBody: `{"title":"test", "description":"testst"}`,
			inputList: todo.TodoList{
				Title:       "test",
				Description: "testst",
			},
			mockBehavior:      func(s *mock_service.MockTodoLists, list todo.TodoList, userId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No required pole",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, list todo.TodoList, userId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Server failure",
			inputBody: `{"title":"test", "description":"testst"}`,
			inputList: todo.TodoList{
				Title:       "test",
				Description: "testst",
			},
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, list todo.TodoList, userId int) {
				s.EXPECT().Create(userId, list).Return(1, fmt.Errorf("server failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"server failure"}`,
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
				if testCase.userId == 0 {
					return
				}
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

func TestHandler_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoLists, userId int)

	testTable := []struct {
		name              string
		userId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId int) {
				s.EXPECT().GetAll(userId).Return([]todo.TodoList{
					{
						Title:       "test",
						Description: "testdesc",
						Id:          1,
					},
				}, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"data":[{"id":1,"title":"test","description":"testdesc"}]}`,
		},
		{
			name:              "No User",
			mockBehavior:      func(s *mock_service.MockTodoLists, userId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId int) {
				s.EXPECT().GetAll(userId).Return(nil, fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoLists(c)
			testCase.mockBehavior(list, testCase.userId)

			services := &service.Service{TodoLists: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.getAllLists)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getListById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoLists, userId, listId int)

	testTable := []struct {
		name              string
		userId            int
		listId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int) {
				s.EXPECT().GetById(userId, listId).Return(todo.TodoList{
					Title:       "test",
					Description: "testdesc",
					Id:          1,
				}, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"id":1,"title":"test","description":"testdesc"}`,
		},
		{
			name:              "No user",
			userId:            0,
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "Invalid Param",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:   "Error Service",
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int) {
				s.EXPECT().GetById(userId, listId).Return(todo.TodoList{}, fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoLists(c)
			testCase.mockBehavior(list, testCase.userId, testCase.listId)

			services := &service.Service{TodoLists: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.getListById)

			var a string
			if testCase.listId != 0 {
				a = fmt.Sprintf("/%d", testCase.listId)
			} else {
				a = "/id"
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", a, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput)

	testTable := []struct {
		name              string
		inputUpdate       string
		updateList        todo.UpdateListInput
		userId            int
		listId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:        "OK",
			inputUpdate: `{}`,
			updateList:  todo.UpdateListInput{},
			userId:      1,
			listId:      1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput) {
				s.EXPECT().Update(userId, listId, input).Return(nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"status":"ok"}`,
		},
		{
			name:              "No User",
			inputUpdate:       `{}`,
			updateList:        todo.UpdateListInput{},
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No List",
			inputUpdate:       `{}`,
			updateList:        todo.UpdateListInput{},
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:        "Service Error",
			inputUpdate: `{}`,
			updateList:  todo.UpdateListInput{},
			userId:      1,
			listId:      1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput) {
				s.EXPECT().Update(userId, listId, input).Return(fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
		{
			name:              "Invalid Body",
			userId:            1,
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int, input todo.UpdateListInput) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoLists(c)
			testCase.mockBehavior(list, testCase.userId, testCase.listId, testCase.updateList)

			services := &service.Service{TodoLists: list}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.updateList)

			var a string
			switch testCase.listId {
			case 0:
				a = "/id"
			default:
				a = fmt.Sprintf("/%d", testCase.listId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", a,
				bytes.NewBufferString(testCase.inputUpdate))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoLists, userId, listId int)

	testTable := []struct {
		name              string
		userId            int
		listId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int) {
				s.EXPECT().Delete(userId, listId).Return(nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"status":"ok"}`,
		},
		{
			name:              "No User",
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No List",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoLists, userId, listId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:   "Service Error",
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoLists, userId, listId int) {
				s.EXPECT().Delete(userId, listId).Return(fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoLists(c)
			testCase.mockBehavior(list, testCase.userId, testCase.listId)

			services := &service.Service{TodoLists: list}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.deleteList)

			var a string
			switch testCase.listId {
			case 0:
				a = "/id"
			default:
				a = fmt.Sprintf("/%d", testCase.listId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", a, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}
