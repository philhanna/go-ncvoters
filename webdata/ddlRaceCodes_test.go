package webdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRaceCodesDDL(t *testing.T) {
	tests := []struct {
		name string
		raceCodes map[string]string
		want string
	}{
		{
			"just2",
			map[string]string{
				"A": "ASIAN",
				"O": "OTHER",
			},
			`
BEGIN TRANSACTION;
DROP TABLE IF EXISTS race_codes;
CREATE TABLE race_codes (
	race TEXT,
	description TEXT
);
INSERT INTO race_codes VALUES('A','ASIAN');
INSERT INTO race_codes VALUES('O','OTHER');
COMMIT;
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := NormalizeSQL(tt.want)
			have := NormalizeSQL(CreateRaceCodesDDL(tt.raceCodes))
			assert.Equal(t, want, have)
		})
	}
}
