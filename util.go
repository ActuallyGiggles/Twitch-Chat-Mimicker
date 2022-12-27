package main

import (
	"math/rand"
	"os"
	"os/exec"
	"time"
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

func RandomNumber(min int, max int) (num int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	num = r.Intn(max-min) + min
	return num
}

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
