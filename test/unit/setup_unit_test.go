package unit

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	err := os.Chdir("../../") //because config file in the root folder
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)
}
