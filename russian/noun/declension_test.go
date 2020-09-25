package declension

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dshipenok/gomorphos/str"
)

func Test_GetDeclension(t *testing.T) {
	decls := GetDeclension(str.Word("лень"), false)

	assert.EqualValues(t, 3, decls)
}

func Test_GetCase(t *testing.T) {
	casedStr := GetCase(str.Word("кухня"), "винительный", false)

	assert.EqualValues(t, "кухню", casedStr)
}

func Test_GetCases(t *testing.T) {
	tests := []struct {
		Word  string
		Cases map[Case]string
	}{
		{
			Word: "коридор",
			Cases: map[Case]string{
				Imenit:  "коридор",
				Rodit:   "коридора",
				Dat:     "коридору",
				Vinit:   "коридор",
				Tvorit:  "коридором",
				Predloj: "коридоре",
			},
		},
		{
			Word: "кухня",
			Cases: map[Case]string{
				Imenit:  "кухня",
				Rodit:   "кухни",
				Dat:     "кухне",
				Vinit:   "кухню",
				Tvorit:  "кухней",
				Predloj: "кухне",
			},
		},
		{
			Word: "бремя",
			Cases: map[Case]string{
				Imenit:  "бремя",
				Rodit:   "бремени",
				Dat:     "бремени",
				Vinit:   "бремя",
				Tvorit:  "бременем",
				Predloj: "бремени",
			},
		},
		{
			Word: "путь",
			Cases: map[Case]string{
				Imenit:  "путь",
				Rodit:   "пути",
				Dat:     "пути",
				Vinit:   "путь",
				Tvorit:  "путем",
				Predloj: "пути",
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.Word, func(t *testing.T) {
			w := str.Word(tst.Word)

			cases := GetCases(w, false)

			assert.Equal(t, tst.Cases, cases)
		})
	}
}
