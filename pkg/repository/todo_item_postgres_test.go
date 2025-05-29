package repository

import (
	todo "do-app"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				row := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(row)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			id: 2,
		},
		{
			name: "Empty fields",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnError(fmt.Errorf("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Empty fields",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				row := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(row)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).
					WillReturnError(fmt.Errorf("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Begin fields",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}

	type mockBehavior func(args args, items []todo.TodoItem)

	testTable := []struct {
		name         string
		args         args
		items        []todo.TodoItem
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				userId: 1,
			},
			items: []todo.TodoItem{
				{
					Title:       "test title",
					Description: "test description",
					Id:          1,
					Done:        false,
				},
			},
			mockBehavior: func(args args, items []todo.TodoItem) {

				row := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(items[0].Id, items[0].Title, items[0].Description, items[0].Done)
				mock.ExpectQuery(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti`).
					WithArgs(args.listId, args.userId).WillReturnRows().WillReturnRows(row)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.items)

			got, err := r.GetAll(testCase.args.userId, testCase.args.listId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.items, got)
			}
		})
	}
}
