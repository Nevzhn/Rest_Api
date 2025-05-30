package repository

import (
	todo "do-app"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestTodoListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		list   todo.TodoList
	}

	type mockBehavior func(args args, listId int)

	testTable := []struct {
		name         string
		listId       int
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:   "OK",
			listId: 1,
			args: args{
				userId: 1,
				list: todo.TodoList{
					Title: "test title",
				},
			},
			mockBehavior: func(args args, listId int) {
				mock.ExpectBegin()

				row := sqlmock.NewRows([]string{"id"}).AddRow(listId)
				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).WillReturnRows(row)

				mock.ExpectExec(`INSERT INTO users_lists`).
					WithArgs(args.userId, listId).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:   "Begin error",
			listId: 1,
			args: args{
				userId: 1,
				list: todo.TodoList{
					Title: "test title",
				},
			},
			mockBehavior: func(args args, listId int) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name:   "1st query empty",
			listId: 1,
			args: args{
				userId: 1,
				list: todo.TodoList{
					Title: "test title",
				},
			},
			mockBehavior: func(args args, listId int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).WillReturnError(fmt.Errorf("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:   "2nd query empty",
			listId: 1,
			args: args{
				userId: 1,
				list: todo.TodoList{
					Title: "test title",
				},
			},
			mockBehavior: func(args args, listId int) {
				mock.ExpectBegin()

				row := sqlmock.NewRows([]string{"id"}).AddRow(listId)
				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).WillReturnRows(row)

				mock.ExpectExec(`INSERT INTO users_lists`).
					WithArgs(args.userId, listId).WillReturnError(fmt.Errorf("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args, testCase.listId)

			got, err := r.Create(testCase.args.userId, testCase.args.list)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.listId, got)
			}
		})
	}
}

func TestTodoListPostgres_GetAll(t *testing.T) {

}
