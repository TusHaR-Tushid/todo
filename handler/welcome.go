package handler

import (
	"context"
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
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Printf("ParseErr : %v", err1)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userID := claims.Id

		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, err := w.Write([]byte(fmt.Sprintf("Hello,%s", claims.Issuer)))
		if err != nil {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		log.Printf("decoder error %v", err)
		return
	}
	idFromUser, err1 := helper.CreateUser(userDetails)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(idFromUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return
	}

	var todoDetails models.Todo
	err := json.NewDecoder(r.Body).Decode(&todoDetails)
	if err != nil {
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
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return
	}
	status := r.URL.Query().Get("status")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		log.Printf("err is :%v", err)
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Printf("err is :%v", err)
	}
	if status == "active" {
		isStatus = true
		isActive = true
	} else if status == "draft" {
		isStatus = true
		isActive = false
	}
	todos, todoErr := helper.GetAllTodo(isStatus, isActive, userID, page, limit)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetCompleted(w http.ResponseWriter, r *http.Request) {
	todos, todoErr := helper.GetCompleted()
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetUpcoming(w http.ResponseWriter, r *http.Request) {
	todos, todoErr := helper.GetUpcoming()
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetExpired(w http.ResponseWriter, r *http.Request) {
	todos, todoErr := helper.GetExpired()
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		log.Printf("parsing error :%v", err)
	}
	var usersTodo models.Todo
	err = json.NewDecoder(r.Body).Decode(&usersTodo)
	if err != nil {
		log.Printf("decoder error %v", err)
		return
	}

	updatedTodo := helper.UpdateTodo(id, usersTodo)
	if updatedTodo != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Printf("parsing error :%v", err)
	}
	todoErr := helper.DeleteTodo(id)
	if todoErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	var userDetails models.UsersLoginDetails
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userCredentials, checkErr := helper.FetchPassword(userDetails.Email)
	if checkErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the expected password

	//If a password exists for the given UserCred
	//And, if it is the same as the password we received, then we can move ahead
	//if NOT, then we return an "Unauthorized" status
	if userCredentials.Password != userDetails.Password {
		w.WriteHeader(http.StatusUnauthorized)
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
			Issuer:    userDetails.Name,
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(tokenString))
	if err != nil {
		log.Printf("encoder error %v", err)
		return
	}

}
