package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Users struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
}

type Todo struct {
	CreatedBy   int       `json:"createdBy"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ExpiringAt  time.Time `json:"expiringAt"`
	IsCompleted bool      `json:"isCompleted"`
}

type UsersTodo struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	Age         int       `json:"age" db:"age"`
	Gender      string    `json:"gender" db:"gender"`
	Address     string    `json:"address" db:"address"`
	CreatedBy   int       `json:"createdBy" db:"created_by"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	ExpiringAt  time.Time `json:"expiringAt" db:"expiring_at"`
	IsCompleted bool      `json:"isCompleted"  db:"is_completed"`
	ArchivedAt  time.Time `json:"archivedAt"  db:"archived_at"`
}

type UsersLoginDetails struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claim struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

type UserCredentials struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

//Storing the users information as in memory map
