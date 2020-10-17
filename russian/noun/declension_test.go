package declension

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dshipenok/gomorphos/russian/cases"
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
		Cases map[cases.Case]string
	}{
		{
			Word: "прохожий",
			Cases: map[cases.Case]string{
				cases.Imenit:  "прохожий",
				cases.Rodit:   "прохожего",
				cases.Dat:     "прохожему",
				cases.Vinit:   "прохожего",
				cases.Tvorit:  "прохожим",
				cases.Predloj: "прохожем",
			},
		},
		{
			Word: "лошадь",
			Cases: map[cases.Case]string{
				cases.Imenit:  "лошадь",
				cases.Rodit:   "лошади",
				cases.Dat:     "лошади",
				cases.Vinit:   "лошадь",
				cases.Tvorit:  "лошадью",
				cases.Predloj: "лошади",
			},
		},
		{
			Word: "коридор",
			Cases: map[cases.Case]string{
				cases.Imenit:  "коридор",
				cases.Rodit:   "коридора",
				cases.Dat:     "коридору",
				cases.Vinit:   "коридор",
				cases.Tvorit:  "коридором",
				cases.Predloj: "коридоре",
			},
		},
		{
			Word: "кухня",
			Cases: map[cases.Case]string{
				cases.Imenit:  "кухня",
				cases.Rodit:   "кухни",
				cases.Dat:     "кухне",
				cases.Vinit:   "кухню",
				cases.Tvorit:  "кухней",
				cases.Predloj: "кухне",
			},
		},
		{
			Word: "бремя",
			Cases: map[cases.Case]string{
				cases.Imenit:  "бремя",
				cases.Rodit:   "бремени",
				cases.Dat:     "бремени",
				cases.Vinit:   "бремя",
				cases.Tvorit:  "бременем",
				cases.Predloj: "бремени",
			},
		},
		{
			Word: "путь",
			Cases: map[cases.Case]string{
				cases.Imenit:  "путь",
				cases.Rodit:   "пути",
				cases.Dat:     "пути",
				cases.Vinit:   "путь",
				cases.Tvorit:  "путем",
				cases.Predloj: "пути",
			},
		},
		{
			Word: "розетка",
			Cases: map[cases.Case]string{
				cases.Imenit:  "розетка",
				cases.Rodit:   "розетки",
				cases.Dat:     "розетке",
				cases.Vinit:   "розетку",
				cases.Tvorit:  "розеткой",
				cases.Predloj: "розетке",
			},
		},
		{
			Word: "лоджия",
			Cases: map[cases.Case]string{
				cases.Imenit:  "лоджия",
				cases.Rodit:   "лоджии",
				cases.Dat:     "лоджии",
				cases.Vinit:   "лоджию",
				cases.Tvorit:  "лоджией",
				cases.Predloj: "лоджии",
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
