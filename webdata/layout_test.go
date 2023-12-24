package webdata

import (
	"fmt"
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
	layout, err := ParseLayoutFile(path)
	assert.Nil(t, err)

	assert.Equal(t, NCOLUMNS, len(layout.AllColumns))
	
	if false {
		for i, column := range layout.AllColumns {
			fmt.Printf("%d: %v\n", i, column)
		}
	}
}
