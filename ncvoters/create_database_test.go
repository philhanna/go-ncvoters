package ncvoters

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInsertSQL(t *testing.T) {
	tests := []struct {
		name string
		cols []string
		want string
	}{
		{"Three columns", []string{"name", "rank", "serial_number"}, "INSERT INTO voters (name,rank,serial_number) VALUES (?,?,?)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			have := CreateInsertSQL(tt.cols)
			assert.Equal(t, want, have)
		})
	}
}

func TestGetSelectedIndices(t *testing.T) {
	all := []string{
		"Larry",
		"Curly",
		"Moe",
	}
	tests := []struct {
		name   string
		subset []string
		want   []int
	}{
		{"Find all", all, []int{0, 1, 2}},
		{"Empty list", []string{}, []int{}},
		{"Just one", []string{"Larry"}, []int{0}},
		{"Just two", []string{"Curly", "Moe"}, []int{1, 2}},
		{"Bogus", []string{"Bogus"}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			have := GetSelectedIndices(all, tt.subset)
			assert.Equal(t, want, have)
		})
	}
}

func TestGetZipEntry(t *testing.T) {
	zipFileName := filepath.Join("..", "testdata", "stooges.zip")

	// Good entry
	entryName := "stooges.csv"
	fp, err := GetZipEntry(zipFileName, entryName)
	assert.Nil(t, err)
	assert.NotNil(t, fp)
	assert.Equal(t, entryName, fp.FileHeader.Name)

	// Bad entry
	entryName = "bogus"
	fp, err = GetZipEntry(zipFileName, entryName)
	assert.NotNil(t, err)
}
