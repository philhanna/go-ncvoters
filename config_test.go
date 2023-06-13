package goncvoters

import (
	"fmt"
	"testing"
)

func Test_newConfiguration(t *testing.T) {
	configuration := *newConfiguration()
	fmt.Printf("%#v\n", configuration)
}
