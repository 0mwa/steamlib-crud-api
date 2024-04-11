package repository

import "database/sql"

type Users struct {
	Db *sql.DB
}

func (u Users) CheckUser(login string, passwd string) (int, error) {
	var err error
	var result *sql.Rows
	result, err = u.Db.Query("SELECT id FROM users WHERE login = $1 AND passwd = $2", login, passwd)
	if err != nil {
		return 0, err
	}
	if result.Next() {
		var userId int
		err = result.Scan(&userId)
		if err != nil {
			return 0, err
		}
		return userId, nil
	}
	return 0, nil
}
