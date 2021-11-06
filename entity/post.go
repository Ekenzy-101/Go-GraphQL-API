package entity

type Post struct {
	ID        string      `json:"id,omitempty"`
	Content   string      `json:"content,omitempty"`
	CreatedAt interface{} `json:"createdAt,omitempty"`
	Title     string      `json:"title,omitempty"`
	User      *User       `json:"user,omitempty"`
	UserID    string      `json:"userId,omitempty"`
	UpdatedAt interface{} `json:"updatedAt,omitempty"`
}

func (p *Post) SetContent(value string) *Post {
	p.Content = value
	return p
}

func (p *Post) SetUser(value *User) *Post {
	p.User = &User{
		ID:   value.ID,
		Name: value.Name,
	}
	return p
}

func (p *Post) SetTitle(value string) *Post {
	p.Title = value
	return p
}

func (p *Post) SetUpdatedAt(value interface{}) *Post {
	p.UpdatedAt = value
	return p
}
