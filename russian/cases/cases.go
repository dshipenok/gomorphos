package cases

import (
	"strings"

	"github.com/dshipenok/gomorphos/str"
)

type Case int //string

const (
	Imenit  Case = 0 //"imenit"
	Rodit        = 1 //"rodit"
	Dat          = 2 //"dat"
	Vinit        = 3 //"vinit"
	Tvorit       = 4 //"tvorit"
	Predloj      = 5 //"predloj"
)

type Cases map[Case]string

func NewCases() map[Case]string {
	return make(map[Case]string, int(Predloj))
}

func NewCasesWord(w str.Word) map[Case]string {
	s := w.String()
	return map[Case]string{
		Imenit:  s,
		Rodit:   s,
		Dat:     s,
		Vinit:   s,
		Tvorit:  s,
		Predloj: s,
	}
}

/**
 * @param string $case
 * @return string
 * @throws \Exception
 */
func CanonizeCase(wCase string) Case {
	wCase = strings.ToLower(wCase)
	switch wCase {
	//  case Imenit:
	case "именительный", "именит", "и":
		return Imenit

		//  case Cases::RODIT:
	case "родительный", "родит", "р":
		return Rodit

		//  case Cases::DAT:
	case "дательный", "дат", "д":
		return Dat

		//  case Cases::VINIT:
	case "винительный", "винит", "в":
		return Vinit

		//  case Cases::TVORIT:
	case "творительный", "творит", "т":
		return Tvorit

		//  case Cases::PREDLOJ:
	case "предложный", "предлож", "п":
		return Predloj

	//  case Cases::LOCATIVE:
	//      return Cases::LOCATIVE;

	//  default:
	//      return \morphos\CasesHelper::canonizeCase($case);
	default:
		return Imenit
	}
}
