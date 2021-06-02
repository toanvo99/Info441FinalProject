package models

import "database/sql"

// Query handler for any requests on our database.

type SQLStore struct {
	Database *sql.DB
}

// initializes new SQLStore struct
func NewSQLStore(database *sql.DB) *SQLStore {
	return &SQLStore{
		Database: database,
	}
}

// grabs a user by their id
func (ss *SQLStore) GetByID(id int64) (*User, error) {
	rows, err := ss.Database.Query("SELECT * FROM users WHERE ID=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &User{}

	for rows.Next() {
		if err2 := rows.Scan(
			&user.TrainerID,
			&user.Email,
			&user.PassHash,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.PhotoURL,
		); err2 != nil {
			return nil, err2
		}
	}
	if err3 := rows.Err(); err3 != nil {
		return nil, err3
	}
	return user, nil
}

// grabs a user by their email
func (ss *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err := ss.Database.Query("SELECT * FROM users WHERE Email=?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &User{}

	for rows.Next() {
		if err2 := rows.Scan(
			&user.TrainerID,
			&user.Email,
			&user.PassHash,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.PhotoURL,
		); err2 != nil {
			return nil, err2
		}
	}
	if err3 := rows.Err(); err3 != nil {
		return nil, err3
	}
	return user, nil
}

// returns a user by their username
func (ss *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err := ss.Database.Query("SELECT * FROM users WHERE UserName=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &User{}

	for rows.Next() {
		if err2 := rows.Scan(
			&user.TrainerID,
			&user.Email,
			&user.PassHash,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.PhotoURL,
		); err2 != nil {
			return nil, err2
		}
	}
	if err3 := rows.Err(); err3 != nil {
		return nil, err3
	}
	return user, nil
}

// inserts a new user into ss.Database and returns user with autofilled id
func (ss *SQLStore) Insert(user *User) (*User, error) {
	insertQuery := "INSERT into users (Email, PassHash, UserName, FirstName, LastName, PhotoURL) VALUES (?,?,?,?,?,?)"
	response, err := ss.Database.Exec(insertQuery,
		user.Email,
		user.PassHash,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.PhotoURL,
	)

	if err != nil {
		return nil, err
	}

	id, err2 := response.LastInsertId()

	if err2 != nil {
		return nil, err2
	}

	user.TrainerID = id
	return user, nil
}

/* updates a user's (found by id) first and/or last name from updates object and returns user with updates,
   returns error if it doesn't work for some reason
*/
func (ss *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	insertQuery := "UPDATE users SET FirstName = ?, LastName = ? WHERE ID=?"
	_, err := ss.Database.Exec(insertQuery,
		updates.FirstName,
		updates.LastName,
		id,
	)
	if err != nil {
		return nil, err
	}
	user, err2 := ss.GetByID(id)
	if err2 != nil {
		return nil, err2
	}
	return user, nil
}

// deletes user by id from ss.Database, returns error if it doesn't work
func (ss *SQLStore) Delete(id int64) error {
	insertQuery := "DELETE FROM users WHERE ID=?"
	_, err := ss.Database.Exec(insertQuery,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SQLStore) InsertUserLog(uLog *UserLog) (*UserLog, error) {
	insertQuery := "INSERT into loginTable (ID, LoginDateTime, clientIP) VALUES (?,?,?)"
	response, err := ss.Database.Exec(insertQuery,
		uLog.LogInID,
		uLog.LogDate,
		uLog.IpAddr,
	)

	if err != nil {
		return nil, err
	}

	id, err2 := response.LastInsertId()

	if err2 != nil {
		return nil, err2
	}

	uLog.LogInID = id
	return uLog, nil
}

func (ss *SQLStore) GetUserLogByID(id int64) (*UserLog, error) {
	rows, err := ss.Database.Query("SELECT * FROM loginTable WHERE LoginID=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uLog := &UserLog{}

	for rows.Next() {
		if err2 := rows.Scan(
			&uLog.LogInID,
			&uLog.ID,
			&uLog.LogDate,
			&uLog.IpAddr,
		); err2 != nil {
			return nil, err2
		}
	}
	if err3 := rows.Err(); err3 != nil {
		return nil, err3
	}
	return uLog, nil
}

func (ss *SQLStore) DeleteUserLog(id int64) error {
	insertQuery := "DELETE FROM loginTable WHERE LoginID=?"
	_, err := ss.Database.Exec(insertQuery,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
