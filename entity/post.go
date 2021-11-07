package entity

import "time"

type Post struct {
	ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
	Content   string     `json:"content,omitempty" bson:"content,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Title     string     `json:"title,omitempty" bson:"title,omitempty"`
	User      *User      `json:"user,omitempty" bson:"user,omitempty"`
	UserID    string     `json:"userId,omitempty" bson:"userId,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (p *Post) SetContent(value string) *Post {
	p.Content = value
	return p
}

func (p *Post) SetCreatedAt(value time.Time) *Post {
	p.CreatedAt = &value
	return p
}

func (p *Post) SetID(value string) *Post {
	p.ID = value
	return p
}

func (p *Post) SetUser(value *User) *Post {
	p.User = &User{
		ID:   value.ID,
		Name: value.Name,
	}
	p.UserID = ""
	return p
}

func (p *Post) SetTitle(value string) *Post {
	p.Title = value
	return p
}

func (p *Post) SetUpdatedAt(value time.Time) *Post {
	p.UpdatedAt = &value
	return p
}
