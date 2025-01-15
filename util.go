package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"time"

	"github.com/pterm/pterm"
)

var (
	Config = Configuration{}
)

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

func secondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%d:%02d", minutes, seconds)
	return str
}

func countdown(timeToWait int) {
	countdown, _ := pterm.DefaultArea.WithRemoveWhenDone().Start(pterm.Gray("Waiting for " + secondsToMinutes(timeToWait) + " seconds..."))
	for i := timeToWait; i >= 0; i-- {
		countdown.Update(pterm.Gray("Waiting for " + secondsToMinutes(i) + " seconds..."))
		time.Sleep(time.Second)
	}
	countdown.Stop()
}
