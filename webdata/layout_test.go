package webdata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadLayout(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"bogus", "BOGUS", "", true},
		{"good", URL, "voter_layout.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			have, err := DownloadLayout(tt.url)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.True(t, strings.HasSuffix(have, tt.want))
			}
		})
	}
}

func TestNewLayout(t *testing.T) {
	const NCOLUMNS = 67
	path := filepath.Join("..", "testdata", "layout_ncvoter.txt")
	blob, err := os.ReadFile(path)
	assert.Nil(t, err)

	body := string(blob)
	reader := strings.NewReader(body)
	layout, err := NewLayout(reader)
	assert.Nil(t, err)

	assert.Equal(t, NCOLUMNS, len(layout.GetColumns()))
	
	if false {
		for i, column := range layout.GetColumns() {
			fmt.Printf("%d: %v\n", i, column)
		}
	}
}
