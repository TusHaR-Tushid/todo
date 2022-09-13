package helper

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-todo/database"
	"go-todo/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func FetchPasswordAndId(userMail string) (models.UserCredentials, error) {

	SQL := `SELECT id, password
          FROM   users
          WHERE  email =$1
          AND   archived_at IS NULL`

	var userCredentials models.UserCredentials

	err := database.UserTodoDb.Get(&userCredentials, SQL, userMail)
	if err != nil {
		log.Printf("FetchPassword:%v", err)
		return userCredentials, err
	}

	return userCredentials, nil
}

func Logout(userID int) error {
	SQL := `UPDATE sessions
            SET    expires_at=now()
            WHERE  user_id=$1`

	_, err := database.UserTodoDb.Exec(SQL, userID)
	if err != nil {
		logrus.Printf("Logout: cannot do logout:%v", err)
		return err
	}

	return nil
}

func FetchUserID(sessionID uuid.UUID) (int, error) {
	SQL := `SELECT user_id
            FROM   sessions
            WHERE  id=$1`

	var userID int
	err := database.UserTodoDb.Get(&userID, SQL, sessionID)
	if err != nil {
		logrus.Printf("FetchUserID: cannot get userID:%v", err)
		return userID, err
	}

	return userID, nil
}

func CreateSession(claims *models.Claim) (uuid.UUID, error) {
	SQL := `INSERT INTO sessions(user_id)
            VALUES   ($1)
            RETURNING id`
	var sessionID uuid.UUID
	err := database.UserTodoDb.Get(&sessionID, SQL, claims.Id)
	if err != nil {
		logrus.Printf("CreateSession: cannot create user session:%v", err)
		return sessionID, err
	}

	return sessionID, nil
}

func CheckSession(userID int) (uuid.UUID, error) {
	SQL := `SELECT id
           FROM    sessions
           WHERE   expires_at IS NULL
           AND     user_id=$1`
	var sessionID uuid.UUID

	err := database.UserTodoDb.Get(&sessionID, SQL, userID)
	if err != nil {
		logrus.Printf("CheckSession: session expired:%v", err)
		return sessionID, err
	}

	return sessionID, nil
}

func CreateUser(userDetails models.Users) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {

		log.Printf("CreateUser: error is:%v", err)
	}
	//language=SQL
	SQL := `INSERT INTO users(name, email, password, age, gender, address)
          VALUES($1, $2, $3, $4, $5, $6)
          RETURNING id`
	var id int
	userDetails.Email = strings.ToLower(userDetails.Email)

	err = database.UserTodoDb.Get(&id, SQL, userDetails.Name, userDetails.Email, hash, userDetails.Age, userDetails.Gender, userDetails.Address)
	if err != nil {
		log.Printf("CreateUser:%v", err)
		return id, err
	}

	return id, nil
}

func CreateTodo(todoDetails *models.TodoInput, userID int) error {
	SQL := `INSERT INTO todo(created_by, title, description, expiring_at)
          VALUES($1, $2, $3, $4)`
	_, err := database.UserTodoDb.Exec(SQL, userID, todoDetails.Title, todoDetails.Description, todoDetails.ExpiringAt)
	if err != nil {
		log.Printf("CreateTodo: unable to create todo: %v", err)
		return err
	}

	return nil
}

func GetAllTodo(conditionCheck models.ConditionCheck, userID int) ([]models.UsersTodoList, error) {
	//language=sql
	SQL := `SELECT 
                   id,
                   title,    
                   description,
                   created_by,
                   created_at,
                   expiring_at,
                   is_completed,
                   is_active
          FROM todo
          WHERE created_by=$6
          AND ($1 OR is_active=$2)
          AND ($3 or title ilike '%' || $4 || '%')
          AND is_completed=$5
          ORDER BY title
          LIMIT  $7
          OFFSET $8`
	todoSlice := make([]models.UsersTodoList, 0)
	err := database.UserTodoDb.Select(&todoSlice, SQL, !conditionCheck.IsStatus, conditionCheck.IsActive, !conditionCheck.IsSearched, conditionCheck.SearchedName, conditionCheck.IsCompleted, userID, conditionCheck.Limit, conditionCheck.Limit*conditionCheck.Page)
	if err != nil {
		log.Printf("GetAllTodo:uable to get todo %v", err)
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
         FROM users JOIN todo ON users.id = todo.created_by
         WHERE users.archived_at IS NULL
               AND todo.is_completed IS TRUE`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL)
	if err != nil {
		log.Printf("GetCompleted: cannot %v", err)
		return nil, err
	}

	return todoSlice, nil
}

func GetUpcoming(userID int) ([]models.UsersTodo, error) {
	SQL := `SELECT   users.id, 
                   name,
                   archived_at,
                   title,    
                   description,
                   todo.created_at,
                   expiring_at,
                   is_completed
          FROM users JOIN todo ON users.id = todo.created_by
          WHERE archived_at IS NULL 
                AND users.id=$1
                AND expiring_at < now()`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL, userID)
	if err != nil {
		log.Printf("GetUpcoming: cannot add upcoming todo to slice:%v", err)
		return nil, err
	}

	return todoSlice, nil
}

func GetExpired(userID int) ([]models.UsersTodo, error) {
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
                AND users.id=$1
                AND expiring_at > now()`
	todoSlice := make([]models.UsersTodo, 0)

	err := database.UserTodoDb.Select(&todoSlice, SQL, userID)
	if err != nil {
		log.Printf("GetExpired:cannot add expired todo to slice:%v", err)
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
		log.Printf("UpdateTodo:cannot update todo : %v", err)
		return err
	}

	return nil
}

func MarkCompleted(id int, todoVar models.Todo) error {
	SQL := `UPDATE todo
           SET    is_completed=$1
           WHERE  id=$2`

	_, err := database.UserTodoDb.Exec(SQL, todoVar.IsCompleted, id)
	if err != nil {
		log.Printf("MarkCompleted: cannot marke todo as completed:%v", err)
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
		log.Printf("DeleteTodo: user cannot be delete todo: %v", err)
		return err
	}

	return nil
}
