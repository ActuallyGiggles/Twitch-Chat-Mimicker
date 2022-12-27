package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"os/exec"
)

var (
	Config = Configuration{}
)

func addUser(user string) {
	u := User{
		Name: user,
	}
	Users = append(Users, u)
}

func RandomNumber(min, max int) int {
	var result int
	switch {
	case min > max:
		// Fail with error
		return result
	case max == min:
		result = max
	case max > min:
		maxRand := max - min
		b, err := rand.Int(rand.Reader, big.NewInt(int64(maxRand)))
		if err != nil {
			return result
		}
		result = min + int(b.Int64())
	}
	return result
}

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
