package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
)

func GenerateUUID() (id string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	ids, _ := uuid.NewV4()
	id = ids.String()
	return
}
