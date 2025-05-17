package service

import (
	todo "do-app"
	"do-app/pkg/repository"
	"fmt"
)

type TodoItemService struct {
	repo     repository.TodoItems
	listRepo repository.TodoLists
}

func NewTodoItemService(repo repository.TodoItems, listRepo repository.TodoLists) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (i *TodoItemService) Create(userId, listId int, input todo.TodoItem) (int, error) {
	_, err := i.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, fmt.Errorf("Create service item: %w", err)
	}
	return i.repo.Create(listId, input)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
