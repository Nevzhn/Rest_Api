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

func TestHandler_createItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem)

	testTable := []struct {
		name              string
		inputBody         string
		inputItem         todo.TodoItem
		userId            int
		listId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"test"}`,
			inputItem: todo.TodoItem{
				Title: "test",
			},
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem) {
				s.EXPECT().Create(userId, listId, input).Return(1, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"id":1}`,
		},
		{
			name:      "No User",
			inputBody: `{"title":"test"}`,
			inputItem: todo.TodoItem{
				Title: "test",
			},
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:      "No List",
			inputBody: `{"title":"test"}`,
			inputItem: todo.TodoItem{
				Title: "test",
			},
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid list id param"}`,
		},
		{
			name:              "invalid input",
			userId:            1,
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Error service",
			inputBody: `{"title":"test"}`,
			inputItem: todo.TodoItem{
				Title: "test",
			},
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, listId int, input todo.TodoItem) {
				s.EXPECT().Create(userId, listId, input).Return(1, fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItems(c)
			testCase.mockBehavior(item, testCase.userId, testCase.listId, testCase.inputItem)

			services := &service.Service{TodoItems: item}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/:id/items", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.createItem)

			var a string
			switch testCase.listId {
			case 0:
				a = "/id/items"
			default:
				a = fmt.Sprintf("/%d/items", testCase.listId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", a,
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAllItems(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItems, userId, listId int)

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
			mockBehavior: func(s *mock_service.MockTodoItems, userId, listId int) {
				s.EXPECT().GetAll(userId, listId).Return([]todo.TodoItem{
					{
						Title: "test",
					},
				}, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `[{"id":0,"title":"test","description":"","done":false}]`,
		},
		{
			name:              "No User",
			listId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, listId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No List",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, listId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid list id param"}`,
		},
		{
			name:   "Service error",
			userId: 1,
			listId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, listId int) {
				s.EXPECT().GetAll(userId, listId).Return([]todo.TodoItem{}, fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItems(c)
			testCase.mockBehavior(item, testCase.userId, testCase.listId)

			services := &service.Service{TodoItems: item}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/:id/items", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.getAllItems)

			var a string
			switch testCase.listId {
			case 0:
				a = "/id/items"
			default:
				a = fmt.Sprintf("/%d/items", testCase.listId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", a, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getItemById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItems, userId, itemId int)

	testTable := []struct {
		name              string
		userId            int
		itemId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			itemId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int) {
				s.EXPECT().GetById(userId, itemId).Return(todo.TodoItem{
					Title: "test",
				}, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"id":0,"title":"test","description":"","done":false}`,
		},
		{
			name:              "No Item",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid item id param"}`,
		},
		{
			name:              "No User",
			itemId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:   "Service Error",
			userId: 1,
			itemId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int) {
				s.EXPECT().GetById(userId, itemId).Return(todo.TodoItem{}, fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItems(c)
			testCase.mockBehavior(item, testCase.userId, testCase.itemId)

			services := &service.Service{TodoItems: item}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.getItemById)

			var a string
			switch testCase.itemId {
			case 0:
				a = "/id"
			default:
				a = fmt.Sprintf("/%d", testCase.itemId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", a, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
			assert.Equal(t, testCase.expectStatusCode, w.Code)
		})
	}
}

func TestHandler_updateItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput)

	testTable := []struct {
		name              string
		inputBody         string
		inputItem         todo.UpdateItemInput
		userId            int
		itemId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{}`,
			inputItem: todo.UpdateItemInput{},
			userId:    1,
			itemId:    1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput) {
				s.EXPECT().Update(userId, itemId, input).Return(nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"status":"ok"}`,
		},
		{
			name:              "No User",
			inputBody:         `{}`,
			inputItem:         todo.UpdateItemInput{},
			itemId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No Item",
			inputBody:         `{}`,
			inputItem:         todo.UpdateItemInput{},
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid item id param"}`,
		},
		{
			name:              "Invalid Body",
			inputItem:         todo.UpdateItemInput{},
			userId:            1,
			itemId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service error",
			inputBody: `{}`,
			inputItem: todo.UpdateItemInput{},
			userId:    1,
			itemId:    1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int, input todo.UpdateItemInput) {
				s.EXPECT().Update(userId, itemId, input).Return(fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItems(c)
			testCase.mockBehavior(item, testCase.userId, testCase.itemId, testCase.inputItem)

			services := &service.Service{TodoItems: item}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.updateItem)

			var a string
			switch testCase.itemId {
			case 0:
				a = "/id"
			default:
				a = fmt.Sprintf("/%d", testCase.itemId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", a,
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItems, userId, itemId int)

	testTable := []struct {
		name              string
		userId            int
		itemId            int
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			itemId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int) {
				s.EXPECT().Delete(userId, itemId).Return(nil)
			},
			expectStatusCode:  200,
			expectRequestBody: `{"status":"ok"}`,
		},
		{
			name:              "No User",
			itemId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int) {},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:              "No Item",
			userId:            1,
			mockBehavior:      func(s *mock_service.MockTodoItems, userId, itemId int) {},
			expectStatusCode:  400,
			expectRequestBody: `{"message":"invalid item id param"}`,
		},
		{
			name:   "Service Error",
			userId: 1,
			itemId: 1,
			mockBehavior: func(s *mock_service.MockTodoItems, userId, itemId int) {
				s.EXPECT().Delete(userId, itemId).Return(fmt.Errorf("service failure"))
			},
			expectStatusCode:  500,
			expectRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItems(c)
			testCase.mockBehavior(item, testCase.userId, testCase.itemId)

			services := &service.Service{TodoItems: item}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/:id", func(ctx *gin.Context) {
				if testCase.userId == 0 {
					return
				}
				ctx.Set(userCtx, testCase.userId)
			}, handler.deleteItem)

			var a string
			switch testCase.itemId {
			case 0:
				a = "/id"
			default:
				a = fmt.Sprintf("/%d", testCase.itemId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", a, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectStatusCode, w.Code)
			assert.Equal(t, testCase.expectRequestBody, w.Body.String())
		})
	}
}
