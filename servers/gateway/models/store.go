package models

import "errors"

// interface for our sql store

//ErrUserNotFound is returned when the user can't be found
var ErrUserNotFound = errors.New("user not found")

//Store represents a store for Users
type Store interface {
	//GetByID returns the User with the given ID
	GetByID(id int64) (*User, error)

	//GetByEmail returns the User with the given email
	GetByEmail(email string) (*User, error)

	//GetByUserName returns the User with the given Username
	GetByUserName(username string) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the DBMS-assigned ID
	Insert(user *User) (*User, error)

	//Update applies UserUpdates to the given user ID
	//and returns the newly-updated user
	Update(id int64, updates *Updates) (*User, error)

	//Delete deletes the user with the given ID
	Delete(id int64) error

	//inserts a login into Login table
	InsertUserLog(userLog *UserLog) (*UserLog, error)

	//returns the userLog with the given ID
	GetUserLogByID(id int64) (*UserLog, error)

	//deletes userlog with given ID
	DeleteUserLog(id int64) error

	//make a new Pokemon team
	MakeNewTeam(id int64) (*Team, error)

	//delete a pokemon team
	DeleteTeam(id int64) error
}
