package repository

import (
	todo "do-app"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"regexp"
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
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
	}

	type mockBehavior func(args args, lists []todo.TodoList)

	testTable := []struct {
		name         string
		lists        []todo.TodoList
		args         args
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			lists: []todo.TodoList{
				{
					Title:       "test Title1 ",
					Description: "test description 1",
				},
				{
					Title: "test Title 2",
				},
				{
					Title:       "test Title3",
					Description: "test description 3",
				},
			},
			args: args{
				userId: 1,
			},
			mockBehavior: func(args args, lists []todo.TodoList) {

				row := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(lists[0].Id, lists[0].Title, lists[0].Description).
					AddRow(lists[1].Id, lists[1].Title, lists[1].Description).
					AddRow(lists[2].Id, lists[2].Title, lists[2].Description)

				mock.ExpectQuery(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl`).
					WithArgs(args.userId).WillReturnRows(row)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Error",
			args: args{
				userId: 1,
			},
			mockBehavior: func(args args, lists []todo.TodoList) {

				mock.ExpectQuery(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl`).
					WithArgs(args.userId).WillReturnError(assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.lists)

			got, err := r.GetAll(testCase.args.userId)

			testCase.wantErr(t, err)
			assert.Equal(t, testCase.lists, got)
		})
	}
}

func TestTodoListPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
	}

	type mockBehavior func(args args, list todo.TodoList)

	testTable := []struct {
		name         string
		list         todo.TodoList
		args         args
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			list: todo.TodoList{
				Title:       "test",
				Id:          2,
				Description: "test desc",
			},
			args: args{
				userId: 1,
				listId: 2,
			},
			mockBehavior: func(args args, list todo.TodoList) {

				row := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(list.Id, list.Title, list.Description)

				mock.ExpectQuery(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl`).
					WithArgs(args.userId, args.listId).WillReturnRows(row)

			},
			wantErr: assert.NoError,
		},
		{
			name: "Error",
			args: args{
				userId: 1,
				listId: 0,
			},
			mockBehavior: func(args args, list todo.TodoList) {

				mock.ExpectQuery(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl`).
					WithArgs(args.userId, args.listId).WillReturnError(assert.AnError)

			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args, testCase.list)

			got, err := r.GetById(testCase.args.userId, testCase.list.Id)

			assert.Equal(t, testCase.list, got)
			testCase.wantErr(t, err)
		})
	}
}

func TestTodoListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func(args args) {

				mock.ExpectExec(`DELETE FROM todo_lists`).
					WithArgs(args.userId, args.listId).WillReturnResult(sqlmock.NewResult(1, 0))

			},
			wantErr: assert.NoError,
		},
		{
			name: "Failed",
			args: args{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func(args args) {

				mock.ExpectExec(`DELETE FROM todo_lists`).
					WithArgs(args.userId, args.listId).WillReturnError(assert.AnError)

			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args)

			err = r.Delete(testCase.args.userId, testCase.args.listId)

			testCase.wantErr(t, err)
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
		input  todo.UpdateListInput
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				listId: 1,
				input: todo.UpdateListInput{
					Title:       stringPointer("title test"),
					Description: stringPointer("desc test"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE todo_lists tl SET title=$1, description=$2 
												FROM users_lists ul WHERE tl.id = ul.list_id AND ul.list_id=$3 
												                        AND ul.user_id=$4`)).
					WithArgs(args.input.Title, args.input.Description, args.listId, args.userId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: assert.NoError,
		},
		{
			name: "Empty Description",
			args: args{
				userId: 1,
				listId: 1,
				input: todo.UpdateListInput{
					Title: stringPointer("title test"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE todo_lists tl SET title=$1
												FROM users_lists ul WHERE tl.id = ul.list_id AND ul.list_id=$2 
												                        AND ul.user_id=$3`)).
					WithArgs(args.input.Title, args.listId, args.userId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: assert.NoError,
		},
		{
			name: "Empty title",
			args: args{
				userId: 1,
				listId: 1,
				input: todo.UpdateListInput{
					Description: stringPointer("desc test"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE todo_lists tl SET description=$1
												FROM users_lists ul WHERE tl.id = ul.list_id AND ul.list_id=$2 
												                        AND ul.user_id=$3`)).
					WithArgs(args.input.Description, args.listId, args.userId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: assert.NoError,
		},
		{
			name: "Failed",
			args: args{
				userId: 1,
				listId: 1,
				input: todo.UpdateListInput{
					Title:       stringPointer("title test"),
					Description: stringPointer("desc test"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE todo_lists tl SET title=$1, description=$2 
												FROM users_lists ul WHERE tl.id = ul.list_id AND ul.list_id=$3 
												                        AND ul.user_id=$4`)).
					WithArgs(args.input.Title, args.input.Description, args.listId, args.userId).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = r.Update(testCase.args.userId, testCase.args.listId, testCase.args.input)

			testCase.wantErr(t, err)
		})
	}
}
