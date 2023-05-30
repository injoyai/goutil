package excel

import (
	"fmt"
	"testing"
)

func TestToString(t *testing.T) {
	fmt.Println(ToString(27, 2))
	fmt.Println(ToString(4, 2))
	fmt.Println(ToString(4, 201))
	fmt.Println(ToString(27, 200))
}
