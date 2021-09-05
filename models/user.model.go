package models

// create a user model
type UserResource struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone,omitempty"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const User = `
CREATE TABLE IF NOT EXISTS users (
    user_id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    phone VARCHAR
)
`
