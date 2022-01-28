package genvar_test

import (
	"fmt"
	"os"

	"github.com/tbhartman/genvar"
)

func ExampleNewOs() {
	myos := genvar.NewOs()
	myos.Setenv("myenv", "1")
	fmt.Println(os.Getenv("myenv"))
	// output:
	// 1
}
