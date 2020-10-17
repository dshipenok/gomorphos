package adjective

import (
	"testing"

	"github.com/dshipenok/gomorphos/russian/cases"
	"github.com/dshipenok/gomorphos/russian/gender"
	"github.com/dshipenok/gomorphos/str"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetCases(t *testing.T) {
	tests := []struct {
		Word        string
		Gender      gender.Gender
		Animateness bool
		Cases       cases.Cases
		Disabled    bool
	}{
		{
			Word:   "красный",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "красный",
				cases.Rodit:   "красного",
				cases.Dat:     "красному",
				cases.Vinit:   "красный",
				cases.Tvorit:  "красным",
				cases.Predloj: "красном",
			},
		},

		{
			Word:   "синяя",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "синяя",
				cases.Rodit:   "синей",
				cases.Dat:     "синей",
				cases.Vinit:   "синюю",
				cases.Tvorit:  "синей",
				cases.Predloj: "синей",
			},
		},

		{
			Word:   "зеленое",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "зеленое",
				cases.Rodit:   "зеленого",
				cases.Dat:     "зеленому",
				cases.Vinit:   "зеленое",
				cases.Tvorit:  "зеленым",
				cases.Predloj: "зеленом",
			},
		},

		{
			Word:        "каленый",
			Gender:      gender.Invalid,
			Animateness: true,
			Cases: cases.Cases{
				cases.Imenit:  "каленый",
				cases.Rodit:   "каленого",
				cases.Dat:     "каленому",
				cases.Vinit:   "каленого",
				cases.Tvorit:  "каленым",
				cases.Predloj: "каленом",
			},
		},

		{
			Word:        "каленый",
			Gender:      gender.Invalid,
			Animateness: false,
			Cases: cases.Cases{
				cases.Imenit:  "каленый",
				cases.Rodit:   "каленого",
				cases.Dat:     "каленому",
				cases.Vinit:   "каленый",
				cases.Tvorit:  "каленым",
				cases.Predloj: "каленом",
			},
		},

		{
			Word:   "хороший",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "хороший",
				cases.Rodit:   "хорошего",
				cases.Dat:     "хорошему",
				cases.Vinit:   "хороший",
				cases.Tvorit:  "хорошим",
				cases.Predloj: "хорошем",
			},
		},

		{
			Word:   "свежий",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "свежий",
				cases.Rodit:   "свежего",
				cases.Dat:     "свежему",
				cases.Vinit:   "свежий",
				cases.Tvorit:  "свежим",
				cases.Predloj: "свежем",
			},
		},

		{
			Word:   "горячий",
			Gender: gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "горячий",
				cases.Rodit:   "горячего",
				cases.Dat:     "горячему",
				cases.Vinit:   "горячий",
				cases.Tvorit:  "горячим",
				cases.Predloj: "горячем",
			},
		},

		{
			Disabled: true, // притяжательные не поддерживаются
			Word:     "волчий",
			Gender:   gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "волчий",
				cases.Rodit:   "волчьего",
				cases.Dat:     "волчьему",
				cases.Vinit:   "волчий",
				cases.Tvorit:  "волчьим",
				cases.Predloj: "волчьем",
			},
		},
		{
			Disabled: true, // притяжательные не поддерживаются (но для таких подходит склонение по правилам существительных)
			Word:     "папин",
			Gender:   gender.Invalid,
			Cases: cases.Cases{
				cases.Imenit:  "папин",
				cases.Rodit:   "папиного",
				cases.Dat:     "папиному",
				cases.Vinit:   "папин",
				cases.Tvorit:  "папиным",
				cases.Predloj: "папином",
			},
		},
	}

	for _, tst := range tests {
		if tst.Disabled {
			continue
		}
		t.Run(tst.Word, func(t *testing.T) {
			w := str.Word(tst.Word)

			cases, err := GetCases(w, tst.Animateness, tst.Gender)

			require.NoError(t, err)
			assert.Equal(t, tst.Cases, cases)
		})
	}
}
