package helper

import (
	"go-todo/database"
	"go-todo/models"
	"log"
)

func FetchPassword(userMail string) (models.UserCredentials, error) {
	SQL := `SELECT id, password
          FROM   users
          WHERE  email =$1`

	var userCredentials models.UserCredentials
	err := database.UserTodoDb.Get(&userCredentials, SQL, userMail)
	if err != nil {
		return userCredentials, err
	}
	return userCredentials, nil
}

func CreateUser(userDetails models.Users) (int, error) {
	SQL := `INSERT INTO users(name,email,password,age,gender,address)
          VALUES($1,$2,$3,$4,$5,$6)
          RETURNING id`
	var id int
	err := database.UserTodoDb.Get(&id, SQL, userDetails.Name, userDetails.Email, userDetails.Password, userDetails.Age, userDetails.Gender, userDetails.Address)
	if err != nil {
		return id, err
	}
	return id, nil
}

func CreateTodo(todoDetails models.Todo, userID int) error {
	SQL := `INSERT INTO todo(created_by,title, description, expiring_at)
          VALUES($1,$2,$3,$4)`
	_, err := database.UserTodoDb.Exec(SQL, userID, todoDetails.Title, todoDetails.Description, todoDetails.ExpiringAt)
	if err != nil {
		log.Printf("CreateTodo: unable to create todo details because of %v", err)
		return err
	}
	return nil
}

func GetAllTodo(isStatus, isActive bool, userID, page, limit int) ([]models.UsersTodo, error) {

	SQL := `SELECT users.id,   
                   name,
                   title,    
                   description,
                   todo.created_at,
                   expiring_at,
                   is_completed
          FROM users
          JOIN todo ON users.id = todo.created_by
          WHERE archived_at IS NULL
          AND ($1 OR is_active=$2)
          AND  users.id=$3
          ORDER BY users.id
          LIMIT  $4
          OFFSET $5`
	todoSlice := make([]models.UsersTodo, 0)
	err := database.UserTodoDb.Select(&todoSlice, SQL, !isStatus, isActive, userID, limit, limit*page)
	if err != nil {
		log.Printf("cannot add todo to slice:%v", err)
		return nil, err
	}
	return todoSlice, nil

}

func GetCompleted() ([]models.UsersTodo, error) {
	SQL := `SELECT   users.id, 
                   users.name,
                   users.archived_at,
                   todo.title,    
                   todo.description,
                   todo.created_at,
                   todo.expiring_at,
                   todo.is_completed
          FROM users
          JOIN todo ON users.id = todo.created_by
          WHERE users.archived_at IS NULL
                AND todo.is_completed IS TRUE`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL)
	if err != nil {
		log.Printf("cannot add todo to slice:%v", err)
		return nil, err
	}
	return todoSlice, nil
}
func GetUpcoming() ([]models.UsersTodo, error) {
	SQL := `SELECT   users.id, 
                   name,
                   archived_at,
                   title,    
                   description,
                   todo.created_at,
                   expiring_at,
                   is_completed
          FROM users
          JOIN todo ON users.id = todo.created_by
          WHERE archived_at IS NULL 
                AND expiring_at < now()`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL)
	if err != nil {
		log.Printf("cannot add todo to slice:%v", err)
		return nil, err
	}
	return todoSlice, nil
}

func GetExpired() ([]models.UsersTodo, error) {
	SQL := `SELECT   users.id, 
                   name,
                   archived_at,
                   title,    
                   description,
                   todo.created_at,
                   expiring_at,
                   is_completed
          FROM users
          JOIN todo ON users.id = todo.created_by
          WHERE archived_at IS NULL 
                AND expiring_at > now()`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL)
	if err != nil {
		log.Printf("cannot add todo to slice:%v", err)
		return nil, err
	}
	return todoSlice, nil
}
func UpdateTodo(id int, todoVar models.Todo) error {

	SQL := `UPDATE todo
          SET    title=$1,
                 description=$2,
                 expiring_at=$3,
                 is_completed=$4,
                 is_active=$5
          WHERE  id=$6`

	_, err := database.UserTodoDb.Exec(SQL, todoVar.Title, todoVar.Description, todoVar.ExpiringAt, todoVar.IsCompleted, todoVar.IsActive, id)
	if err != nil {
		log.Printf("cannot update todo due to %v", err)
		return err
	}
	return nil
}
func DeleteTodo(id int) error {
	SQL := `UPDATE users
          SET    archived_at=now()
          WHERE  id=$1`

	_, err := database.UserTodoDb.Exec(SQL, id)
	if err != nil {
		log.Printf("user cannot be deleted %v", err)
		return err
	}
	return nil
}
