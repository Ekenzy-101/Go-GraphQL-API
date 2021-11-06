package entity

type User struct {
	ID        string      `json:"id,omitempty"`
	CreatedAt interface{} `json:"createdAt,omitempty"`
	Name      string      `json:"name,omitempty"`
	Email     string      `json:"email,omitempty"`
	Password  string      `json:"password,omitempty"`
	Posts     []Post      `json:"posts,omitempty"`
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
