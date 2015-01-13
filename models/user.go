package models

import (
	"crypto/rand"
	"fmt"
	"io"
)

type User struct {
	ID    string
	Key   string
	Admin bool
}

func GetUser(id string) (*User, error) {
	r, err := DB.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	u := User{}
	for r.Next() {
		err = r.Scan(&u.ID, &u.Key, &u.Admin)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func CreateUser() (*User, error) {
	usr := User{
		ID:    generateID(),
		Key:   generateKey(),
		Admin: false,
	}

	stmt, err := DB.Prepare("INSERT INTO users (id, key, admin) VALUES ($1, $2, FALSE)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(usr.ID, usr.Key)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func DeleteUser(id string) error {
	// Delete
	stmt, err := DB.Prepare("delete from users where id=$1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
		fmt.Println(err)
	}
	return nil
}

func ListUsers() (*[]User, error) {
	users := []User{}
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		return &users, err
	}

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Key, &u.Admin)
		if err != nil {
			return &users, err
		}
		users = append(users, u)
	}
	return &users, nil
}

func generateKey() string {
	b := make([]byte, 10)
	io.ReadFull(rand.Reader, b)
	for i := range b {
		if b[i] > 0x7E {
			b[i] = b[i] & 0x7E
		}
		if b[i] < 0x21 {
			b[i] = b[i] | 0x21
		}
	}
	return string(b)
}
