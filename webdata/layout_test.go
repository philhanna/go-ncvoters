package webdata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLayout(t *testing.T) {
	path := filepath.Join("..", "testdata", "layout_ncvoter.txt")
	blob, err := os.ReadFile(path)
	assert.Nil(t, err)
	body := string(blob)
	reader := strings.NewReader(body)
	layout := NewLayout(reader)
	assert.Equal(t, 67, len(layout.GetColumns()))
	if true {
		for i, column := range layout.GetColumns() {
			fmt.Printf("%d: %v\n", i, column)
		}
	}
}
