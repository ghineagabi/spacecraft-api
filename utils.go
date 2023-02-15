package main

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

const SUCCESSFUL = "Successful operation"
const NOTENOUGHPARAMS = "Not enough parameters provided to fulfill the request."

var Db *sql.DB
var Cred FileCredentials

type FileCredentials struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Dbname   string `json:"dbname" binding:"required"`
}

// PopulateConfig Attempts to open a json file and return its context to the global Cred variable
func PopulateConfig(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, &Cred)
	if err != nil {
		panic(err)
	}
}

// ValidatePassword is currently unused, but the main idea is to use it as a custom validator for the password field
func ValidatePassword(s string) bool {
	letters := 0
	var number, upper, sevenOrMore, lower, special bool
	for _, c := range s {
		switch {
		case unicode.IsLower(c):
			lower = true
			letters++
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			letters++
		case c > unicode.MaxASCII || c == ' ':
			return false
		}
	}
	sevenOrMore = letters >= 8 && letters < 32
	if sevenOrMore && number && upper && lower && special {
		return true
	}
	return false
}

// OnlyUnicode also returns false if the string is empty
func OnlyUnicode(s string) bool {
	if s = strings.TrimSpace(s); s == "" {
		return false
	}
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// Returns the SHA512 equivalent of the provided string
func SHA512(text string) string {
	h := sha512.New512_256()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// Wrapper for specific errors regarding fields that need validation from the request
type InvalidFieldsError struct {
	AffectedField string
	Reason        string
	Location      string
}

func (m *InvalidFieldsError) Error() string {
	return fmt.Sprintf("Cannot process <%s> field: <%s>. Reason: <%s>", m.Location, m.AffectedField, m.Reason)
}

// DecodeAuth decodes the Authorization Header (from base64 to string format) and returns an UserCredentials
func DecodeAuth(auth string) (UserCredentials, error) {
	if strings.HasPrefix(auth, "Basic ") {
		sDec, err := b64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return UserCredentials{}, err
		}
		name, pass, found := bytes.Cut(sDec, []byte{58}) // Separate by ":"
		if found {
			return UserCredentials{Email: string(name), Password: string(pass)}, nil
		}
		return UserCredentials{}, &InvalidFieldsError{"Authorization", "Invalid format. Missing colon ", "Basic auth"}
	}
	return UserCredentials{}, &InvalidFieldsError{"Authorization", "Can only process Basic Authentication", "Basic auth"}
}
