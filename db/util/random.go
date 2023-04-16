package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// init func running import --> const --> var --> init()
func init() {

	fmt.Println("Init function running")
	t := time.Now().UnixNano()

	rand.New(rand.NewSource(t))

}

// RandomString generates a random string of length
func RandomString(n int) string {

	var sb strings.Builder
	l := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)

	}

	return sb.String()

}

// RandomOwner generate random owner name
func RandomOwner(l int) string {
	return RandomString(l)
}

// RandomEmail generate random email address
func RandomEmail(l int) string {

	return fmt.Sprintf("%s@gmail.com", RandomString(l))
}
