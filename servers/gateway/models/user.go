package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// This .go file will handle any potential interactions within trainer models, our database, and API. This will likely
// be very important for our authorization and authentication.

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

var ErrInsertingUser = errors.New("can't insert user into db")

//User representing a user account in the database
type User struct {
	TrainerID int64  `json:"trainerID"`
	Email     string `json:"-"`
	PassHash  []byte `json:"-"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// If we want to include logs in the future
/* type UserLog struct {
	LogInID int64  `json:"loginID"`
	LogDate string `json:"loginDate"`
	IpAddr  string `json:"ipAddr"`
} */

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("password needs to be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("password and passwordConf don't match")
	}
	if len(nu.UserName) == 0 {
		return fmt.Errorf("can't have a 0-length username")
	}
	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("username can't have spaces")
	}
	return nil
}

func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	validUser := User{}
	validUser.Email = strings.Trim(nu.Email, " ")
	h := md5.New()
	h.Write([]byte(strings.ToLower(strings.Trim(nu.Email, " "))))
	hash := hex.EncodeToString(h.Sum(nil))
	url := []string{gravatarBasePhotoURL, hash}
	validUser.PhotoURL = strings.Join(url, "")

	validUser.FirstName = strings.Trim(nu.FirstName, " ")
	validUser.LastName = strings.Trim(nu.LastName, " ")
	validUser.UserName = nu.UserName
	err2 := validUser.SetPassword(nu.Password)
	if err2 != nil {
		return nil, err2
	}

	return &validUser, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = pass
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if updates == nil {
		return fmt.Errorf("update is nil")
	}
	first := strings.Trim(updates.FirstName, " ")
	last := strings.Trim(updates.LastName, " ")
	if len(first) == 0 && len(last) == 0 {
		return fmt.Errorf("update is empty")
	}
	u.FirstName = first
	u.LastName = last
	return nil
}
