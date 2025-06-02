package repository

import (
	"database/sql"
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
					Title:       "test title 1",
					Description: "test description",
					Id:          1,
					Done:        false,
				},
				{
					Title: "test title 2",
					Id:    2,
					Done:  true,
				},
				{
					Title:       "test title 3",
					Id:          3,
					Description: "jopa",
					Done:        true,
				},
			},
			mockBehavior: func(args args, items []todo.TodoItem) {

				row := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(items[0].Id, items[0].Title, items[0].Description, items[0].Done).
					AddRow(items[1].Id, items[1].Title, items[1].Description, items[1].Done).
					AddRow(items[2].Id, items[2].Title, items[2].Description, items[2].Done)
				mock.ExpectQuery(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti`).
					WithArgs(args.listId, args.userId).WillReturnRows().WillReturnRows(row)
			},
		},
		{
			name: "Select Error",
			args: args{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(args args, items []todo.TodoItem) {

				mock.ExpectQuery(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti`).
					WithArgs(args.listId, args.userId).WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.items)

			got, err := r.GetAll(testCase.args.userId, testCase.args.listId)

			if testCase.wantErr {
				assert.Error(t, err)
				log.Println(err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.items, got)
			}
		})
	}
}

func TestTodoItemPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
	}

	type mockBehavior func(args args, item todo.TodoItem)

	testTable := []struct {
		name         string
		item         todo.TodoItem
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			item: todo.TodoItem{
				Id:          1,
				Title:       "test title",
				Description: "test description",
				Done:        true,
			},
			args: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func(args args, item todo.TodoItem) {
				row := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(item.Id, item.Title, item.Description, item.Done)
				mock.ExpectQuery(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti`).
					WithArgs(args.itemId, args.userId).WillReturnRows(row)
			},
		},
		{
			name: "Error",
			args: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func(args args, item todo.TodoItem) {
				mock.ExpectQuery(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti`).
					WithArgs(args.itemId, args.userId).WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args, testCase.item)

			got, err := r.GetById(testCase.args.userId, testCase.args.itemId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.item, got)
			}
		})
	}
}

func TestTodoItemPostgres_Delete(t *testing.T) {

	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(`DELETE FROM todo_items ti USING lists_items li, users_lists ul`).
					WithArgs(args.userId, args.itemId).WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name: "Not Found",
			args: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(`DELETE FROM todo_items ti USING lists_items li, users_lists ul`).
					WithArgs(args.userId, args.itemId).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args)

			err = r.Delete(testCase.args.userId, testCase.args.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

func TestTodoItemPostgres_update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
		input  todo.UpdateItemInput
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				itemId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("title test"),
					Description: stringPointer("description test"),
					Done:        boolPointer(true),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul 
                    						 WHERE (.+)`).
					WithArgs(args.input.Title, args.input.Description, args.input.Done, args.userId, args.itemId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name: "OK_WithoutDescription",
			args: args{
				userId: 1,
				itemId: 1,
				input: todo.UpdateItemInput{
					Title: stringPointer("title test"),
					Done:  boolPointer(true),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul 
                    						 WHERE (.+)`).
					WithArgs(args.input.Title, args.input.Done, args.userId, args.itemId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name: "Empty input",
			args: args{
				userId: 1,
				itemId: 1,
				input:  todo.UpdateItemInput{},
			},
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE todo_items ti SET FROM lists_items li, users_lists ul 
                    						 WHERE (.+)`).
					WithArgs(args.userId, args.itemId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args)

			err = r.Update(testCase.args.userId, testCase.args.itemId, testCase.args.input)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
