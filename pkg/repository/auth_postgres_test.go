package repository

import (
	todo "do-app"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		input todo.User
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		args         args
		id           int
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				input: todo.User{
					Name:     "Test name",
					Username: "username test",
					Password: "qwerty",
				},
			},
			id: 1,
			mockBehavior: func(args args, id int) {

				row := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(args.input.Name, args.input.Username, args.input.Password).
					WillReturnRows(row)

			},
			wantErr: assert.NoError,
		},
		{
			name: "Failed",
			mockBehavior: func(args args, id int) {

				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(args.input.Name, args.input.Username, args.input.Password).
					WillReturnError(assert.AnError)

			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.CreateUser(testCase.args.input)

			testCase.wantErr(t, err)
			assert.Equal(t, testCase.id, got)
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username string
		password string
	}

	type mockBehavior func(args args, user todo.User)

	testTable := []struct {
		name         string
		args         args
		user         todo.User
		mockBehavior mockBehavior
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				username: "test",
				password: "qwerty",
			},
			user: todo.User{
				Id:       1,
				Name:     "testt",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(args args, user todo.User) {

				row := sqlmock.NewRows([]string{"id", "name", "username", "password"}).
					AddRow(user.Id, user.Name, user.Username, user.Password)
				mock.ExpectQuery(`SELECT id FROM users`).
					WithArgs(args.username, args.password).WillReturnRows(row)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Failed",
			args: args{
				username: "test",
				password: "qwerty",
			},
			mockBehavior: func(args args, user todo.User) {

				mock.ExpectQuery(`SELECT id FROM users`).
					WithArgs(args.username, args.password).WillReturnError(assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.user)

			got, err := r.GetUser(testCase.args.username, testCase.args.password)

			testCase.wantErr(t, err)
			assert.Equal(t, testCase.user, got)
		})
	}
}
