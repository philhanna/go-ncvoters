package ncvoters

import "testing"

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
			if got := CreateInsertSQL(tt.cols); got != tt.want {
				t.Errorf("CreateInsertSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
