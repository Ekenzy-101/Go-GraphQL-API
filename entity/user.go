package entity

import "time"

type User struct {
	ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Name      string     `json:"name,omitempty" bson:"name,omitempty"`
	Email     string     `json:"email,omitempty" bson:"email,omitempty"`
	Password  string     `json:"password,omitempty" bson:"password,omitempty"`
	Posts     []Post     `json:"posts,omitempty" bson:"posts,omitempty"`
}

func (u *User) SetCreatedAt(value time.Time) *User {
	u.CreatedAt = &value
	return u
}

func (u *User) SetID(value string) *User {
	u.ID = value
	return u
}

func (u *User) SetPassword(value string) *User {
	u.Password = value
	return u
}

// Accepts interfaces, returns struct
