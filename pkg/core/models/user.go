package models

type User struct {
	ID     int64
	Email  string
	PWHash []byte
}
