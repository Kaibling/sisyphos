package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 11

func NewULID() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

// func GetContext(key reqctx.String, r *http.Request) interface{} {
// 	parameter := r.Context().Value(key)
// 	// if parameter == nil {
// 	// 	panic(apierrors.NewClientError(errors.New("context parameter '" + key + "' missing")))
// 	// }
// 	return parameter
// }

func PrettyJSON(i interface{}) {
	a, _ := json.MarshalIndent(i, "", " ")
	fmt.Println(string(a))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Printf("error: %v\n", err)
	return err == nil
}

func ToPointer[T any](d T) *T {
	return &d
}

func ReadPtr[T any](d *T) T {
	if d == nil {
		return *new(T)
	}
	return *d
}

func PtrRead[T any](d *T) T {
	if d == nil {
		return *new(T)
	}
	return *d
}
