package hashing

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var testTime time.Time

func TestMain(m *testing.M) {
	fmt.Println("Set up something for test")
	testTime = time.Now()
	exitVal := m.Run()
	fmt.Println("Cleaning after testing is finished")
	os.Exit(exitVal)
}
func Test_checkPassword(t *testing.T) {
	fmt.Println(testTime)
	password := "really long and strong password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	checkResult := CheckPassword(password, hash)
	if !checkResult {
		t.Error("password and hash didn't pass check")
	}
}
