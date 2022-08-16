package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"go-todo/database/helper"
	"go-todo/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("secret_key")

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("token")

		claims := models.Claim{}

		tkn, err1 := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err1 != nil {
			if err1 == jwt.ErrSignatureInvalid {
				log.Printf("Signature invalid:%v", err1)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Printf("ParseErr : %v", err1)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("token is invalid")
			return
		}
		userID := claims.Id

		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))

		_, err := w.Write([]byte(fmt.Sprintf("Hello,%s", claims.Issuer)))
		if err != nil {
			//statuscode
			w.WriteHeader(http.StatusInternalServerError)
			//error
			log.Printf("Write error:%v", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("decoder error %v", err)
		return
	}
	idFromUser, err1 := helper.CreateUser(userDetails)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("CreateUser:error is:%v", err1)
		return
	}
	err = json.NewEncoder(w).Encode(idFromUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Encodeing error:%v", err)
		return
	}

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("CreateTodo:QueryParam for userID:%v", ok)
		return
	}

	var todoDetails models.TodoInput
	err := json.NewDecoder(r.Body).Decode(&todoDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("decoder error %v", err)
		return
	}

	err1 := helper.CreateTodo(todoDetails, userID)
	if err1 != nil {
		log.Printf("CretewTodoErrror : %v", err1)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	var isStatus bool
	var isActive bool
	isSearched := false
	// implementing isCompleted
	isCompleted := false
	searchedName := r.URL.Query().Get("name")
	if searchedName != "" {
		isSearched = true
	}
	statusIsCompleted := r.URL.Query().Get("isCompleted")
	if statusIsCompleted == "true" {
		isCompleted = true
	}
	// implemented isCompleted
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("GetAllTodo:QueryParam for userID:%v", ok)
		return
	}

	status := r.URL.Query().Get("status")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Page:err is :%v", err)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Limit:err is :%v", err)
		return
	}

	if status == "active" {
		isStatus = true
		isActive = true
	} else if status == "draft" {
		isStatus = true
		isActive = false
	}

	todos, todoErr := helper.GetAllTodo(isCompleted, isStatus, isActive, isSearched, searchedName, userID, page, limit)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("GetAllTodo:error is:%v", todoErr)
		return
	}

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Encoding error:%v", err)
		return
	}
}

//func GetCompletedTodo(w http.ResponseWriter, r *http.Request) {
//	todos, todoErr := helper.GetCompleted()
//	if todoErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	err := json.NewEncoder(w).Encode(todos)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}

func GetUpcomingTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("GetUpcomingTodo:UserId Query Param error:%v", ok)
		return
	}
	todos, todoErr := helper.GetUpcoming(userID)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("GetUpcoming:%v", todoErr)
		return
	}
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Encoding error:%v", err)
		return
	}
}

func GetExpiredTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("GetExpiredTodo:UserId Query Param error:%v", ok)
		return
	}

	todos, todoErr := helper.GetExpired(userID)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("GetExpired:error is:%v", todoErr)
		return
	}

	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Encoding error:%v", err)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("parsing error :%v", err)
		return
	}

	var usersTodo models.Todo
	err = json.NewDecoder(r.Body).Decode(&usersTodo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("decoder error %v", err)
		return
	}

	updatedTodoErr := helper.UpdateTodo(id, usersTodo)
	if updatedTodoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("UpdateTodo:not able to update todo:%v", updatedTodoErr)
		return
	}

}

func MarkCompleted(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("parsing error :%v", err)
		return
	}

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("decoder error %v", err)
		return
	}

	MarkCompletedErr := helper.MarkCompleted(id, todo)
	if MarkCompletedErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("MarkCompleted: %v", MarkCompletedErr)
		return
	}

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	//userID, ok := r.Context().Value("userID").(int)
	//if !ok {
	//	return
	//}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("parsing error :%v", err)
		return
	}

	todoErr := helper.DeleteTodo(id)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DeleteTodo:error is :%v", todoErr)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	var userDetails models.UsersLoginDetails
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Decoder error:%v", err)
		return
	}

	userCredentials, checkErr := helper.FetchPasswordAndId(userDetails.Email)
	if checkErr != nil && checkErr != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("FetchPassword:error is:%v", err)
		return
	}
	// get the expected password

	//If a password exists for the given UserCred
	//And, if it is the same as the password we received, then we can move ahead
	//if NOT, then we return an "Unauthorized" status
	if userCredentials.Password != userDetails.Password {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("password misMatch")
		return
	}

	expiresAt := time.Now().Add(60 * time.Minute)

	//claims := models.Claim{
	//	Name: userDetails.Name,
	//	StandardClaims: jwt2.StandardClaims{
	//		ExpiresAt: expiresAt.Unix(),
	//	},
	//}
	claims := &models.Claim{
		Id: userCredentials.Id,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("TokenString:cannot create token string:%v", err)
		return
	}
	_, err = w.Write([]byte(tokenString))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("encoder error %v", err)
		return
	}

}
