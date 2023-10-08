package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgtype"
)

// RandomString generateRandomName
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// RandomEmail generateRandomEmail
func RandomEmail() string {
	return RandomString(10) + "@example.com"
}

// RandomUser  generateRandomOwner
func RandomUser() string {
	return RandomString(10)
}

// RandomInt generateRandomInt
func RandomInt() int {
	return rand.Int()
}

func RandomDepartment() string {
	allDepartments := []string{"ACCOUNTS", "PROCUREMENT", "FINANCE", "ADMIN"}
	return allDepartments[rand.Intn(len(allDepartments))]
}

// IntToPgNumeric int to pg numeric
func IntToPgNumeric(i int) pgtype.Numeric {
	var n pgtype.Numeric
	n.Set(i)
	return n
}

// RandomTime random time and time after
func RandomTime() time.Time {
	return time.Now().Add(time.Duration(rand.Intn(1000)) * time.Hour)
}

func FloatToPgNumeric(val float64) pgtype.Numeric {
	strVal := fmt.Sprintf("%.10f", val)
	var numericValue pgtype.Numeric
	err := numericValue.Set(strVal)
	if err != nil {
		panic(err)
	}

	return numericValue
}
