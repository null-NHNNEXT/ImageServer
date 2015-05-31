package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const jwt_secret = "mysecret"

func getToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, keyFn)
	if err == nil {
		token, err = validate(token)
	}
	return token, err
}

func validate(token *jwt.Token) (*jwt.Token, error) {
	var err error
	if !token.Valid {
		err = fmt.Errorf("invalid token")
	}
	return token, err
}

func keyFn(token *jwt.Token) (interface{}, error) {
	return []byte(jwt_secret), nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head></head><body><h1>NHN NEXT Human Design Project!</h1><h2>Image Server by Go language</h2></body></html>")
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get called")

	vars := mux.Vars(r)
	userId := vars["userId"]

	fmt.Println("get called")
	http.ServeFile(w, r, "/home/null/Images/Profiles/"+userId+".jpg")
}

func UploadProfileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload called")

	vars := mux.Vars(r)
	userId := vars["userId"]

	// Check var Id is same as JWT id
	token := r.Header.Get("Authorization")
	parsedToken, err := getToken(token)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if userId != parsedToken.Claims["uuid"] {
		fmt.Fprintln(w, "No rights")
		return
	}

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	out, err := os.Create("/home/null/Images/Profiles/" + userId + ".jpg")
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		fmt.Println(err)
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
}

func main() {
	fmt.Println("main called")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/images/profile/{userId}", GetProfileHandler).Methods("GET")
	router.HandleFunc("/images/profile/{userId}", UploadProfileHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
