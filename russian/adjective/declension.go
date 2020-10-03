package adjective

import (
	"errors"

	"github.com/dshipenok/gomorphos/russian"
	"github.com/dshipenok/gomorphos/russian/cases"
	"github.com/dshipenok/gomorphos/russian/gender"
	"github.com/dshipenok/gomorphos/str"
)

/**
 * Class AdjectiveDeclension.
 *
 * Склонение прилагательных.
 *
 * Правила склонения:
 * - http://www.fio.ru/pravila/grammatika/sklonenie-prilagatelnykh-v-russkom-yazyke/
 *
 * @package morphos\Russian
 */
func IsMutable(w str.Word) bool {
	return false
}

/**
* @param string $adjective
* @param string $case
* @param bool   $animateness
* @param null   $gender
*
* @return string
* @throws \Exception
 */
func GetCase(w str.Word, wCase string, animateness bool, gendr gender.Gender) (string, error) {
	cCase := cases.CanonizeCase(wCase)

	if gendr == gender.Invalid {
		gendr = DetectGender(w, nil)
	}

	forms, err := GetCases(w, animateness, gendr)
	if err != nil {
		return w.String(), err
	}
	return forms[cCase], nil
}

/**
* @param string $adjective
*
* @param bool $isEmphasized
*
* @return string
 */
func DetectGender(w str.Word, isEmphasized *bool) gender.Gender {
	lastChars := w.Lower().LastChars(2)
	switch lastChars {
	case "ой", "ый", "ий":
		if isEmphasized != nil && lastChars == "ой" {
			*isEmphasized = true
		}
		return gender.Male
	case "ая", "яя":
		return gender.Female
	case "ое", "ее":
		return gender.Neuter
	}
	return gender.Invalid
}

/**
 * @param string $adjective
 * @param bool $animateness
 * @param null|string $gender
 *
 * @return string[]
 * @phpstan-return array<string, string>
 */
func GetCases(w str.Word, animateness bool, gendr gender.Gender) (cases.Cases, error) {
	if gendr == gender.Invalid {
		isEmphasized := false
		gendr = DetectGender(w, &isEmphasized)
		if gendr == gender.Invalid {
			return cases.NewCasesWord(w), errors.New("unable to detect adjective gender")
		}
	}

	lastConsonantVowel := w.SliceWord(-2, -1)
	baseType := GetAdjectiveBaseType(w)
	w = w.SliceWord(0, -2)

	switch baseType {
	case HardBase:
		return declinateHardAdjective(w, animateness, gendr, lastConsonantVowel)

	case SoftBase:
		return declinateSoftAdjective(w, animateness, gendr, lastConsonantVowel)

	case MixedBase:
		return declinateMixedAdjective(w, animateness, gendr, lastConsonantVowel)
	}

	return cases.NewCasesWord(w), errors.New("invalid adjective base")
}

/**
* @param string $adjective
*
* @return int
 */
func GetAdjectiveBaseType(w str.Word) int {
	w = w.Lower()

	substring := russian.FindLastPositionForOneOfChars(w, russian.ConsonantsAdj)
	lastConsonant := substring.SliceWord(0, 1)

	// г, к, х, ударное ш - признак смешанного прилагательно
	if lastConsonant.OneOf("г", "к", "х") ||
		(lastConsonant.String() == "ш" && substring.SliceWord(1, 2).OneOf("о", "а")) {
		return MixedBase
	}

	if russian.CheckBaseLastConsonantSoftness(substring) || substring.SliceWord(0, 2).String() == "шн" {
		return SoftBase
	}
	return HardBase
}

/**
* @param string $adjective
* @param bool   $animateness
* @param string $gender
* @param string $afterConsonantVowel
*
* @return string[]
* @phpstan-return array<string, string>
 */
func declinateHardAdjective(w str.Word, animateness bool, gendr gender.Gender, afterConsonantVowel str.Word) (cases.Cases, error) {
	var postfix string
	switch gendr {
	case gender.Male:
		postfix = afterConsonantVowel.Concat("й")

	case gender.Female:
		postfix = afterConsonantVowel.Concat("я")

	case gender.Neuter:
		postfix = afterConsonantVowel.Concat("е")

	default:
		return cases.NewCasesWord(w), errors.New("invalid gender in hard adjective")
	}

	var rodit, dat, vinit, tvorit, predloj string
	if gendr != gender.Female {
		rodit = w.Concat("о", "го")
		dat = w.Concat("о", "му")
		tvorit = w.Concat("ым")
		predloj = w.Concat("ом")
	} else {
		rodit = w.Concat("о", "й")
		dat = w.Concat("о", "й")
		vinit = w.Concat("ую")
		tvorit = w.Concat("ой")
		predloj = w.Concat("ой")
	}

	cCases := cases.Cases{
		cases.Imenit:  w.Concat(postfix),
		cases.Rodit:   rodit,
		cases.Dat:     dat,
		cases.Vinit:   vinit,
		cases.Tvorit:  tvorit,
		cases.Predloj: predloj,
	}

	if gendr != gender.Female {
		cCases[cases.Vinit] = russian.GetVinitCaseByAnimateness(cCases, animateness)
	}

	return cCases, nil
}

/**
* @param string $adjective
* @param bool   $animateness
* @param string $gender
* @param string $afterConsonantVowel
*
* @return string[]
* @phpstan-return array<string, string>
 */
func declinateSoftAdjective(w str.Word, animateness bool, gendr gender.Gender, afterConsonantVowel str.Word) (cases.Cases, error) {
	var postfix string
	switch gendr {
	case gender.Male:
		postfix = afterConsonantVowel.Concat("й")

	case gender.Female:
		postfix = afterConsonantVowel.Concat("я")

	case gender.Neuter:
		postfix = afterConsonantVowel.Concat("е")

	default:
		return cases.NewCasesWord(w), errors.New("invalid gender in soft adjective")
	}

	var rodit, dat, vinit, tvorit, predloj string
	if gendr != gender.Female {
		rodit = w.Concat("е", "го")
		dat = w.Concat("е", "му")
		tvorit = w.Concat("им")
		predloj = w.Concat("ем")
	} else {
		rodit = w.Concat("е", "й")
		dat = w.Concat("е", "й")
		vinit = w.Concat("юю")
		tvorit = w.Concat("ей")
		predloj = w.Concat("ей")
	}

	cCases := cases.Cases{
		cases.Imenit:  w.Concat(postfix),
		cases.Rodit:   rodit,
		cases.Dat:     dat,
		cases.Vinit:   vinit,
		cases.Tvorit:  tvorit,
		cases.Predloj: predloj,
	}

	if gendr != gender.Female {
		cCases[cases.Vinit] = russian.GetVinitCaseByAnimateness(cCases, animateness)
	}

	return cCases, nil
}

/**
* @param string $adjective
* @param bool   $animateness
* @param string $gender
* @param string $afterConsonantVowel
*
* @return string[]
* @phpstan-return array<string, string>
 */
func declinateMixedAdjective(w str.Word, animateness bool, gendr gender.Gender, afterConsonantVowel str.Word) (cases.Cases, error) {
	var postfix string
	switch gendr {
	case gender.Male:
		postfix = afterConsonantVowel.Concat("й")

	case gender.Female:
		postfix = afterConsonantVowel.Concat("я")

	case gender.Neuter:
		postfix = afterConsonantVowel.Concat("е")

	default:
		return cases.NewCasesWord(w), errors.New("invalid gender in mixed adjective")
	}

	var rodit, dat, tvorit, predloj string
	if gendr != gender.Female {
		rodit = w.Concat("о", "го")
		dat = w.Concat("о", "му")
		tvorit = w.Concat("им")
		predloj = w.Concat("ом")
	} else {
		rodit = w.Concat("о", "й")
		dat = w.Concat("о", "й")
		tvorit = w.Concat("ой")
		predloj = w.Concat("ой")
	}

	cCases := cases.Cases{
		cases.Imenit:  w.Concat(postfix),
		cases.Rodit:   rodit,
		cases.Dat:     dat,
		cases.Tvorit:  tvorit,
		cases.Predloj: predloj,
	}

	switch gendr {
	case gender.Male:
		cCases[cases.Vinit] = russian.GetVinitCaseByAnimateness(cCases, animateness)
	case gender.Neuter:
		cCases[cases.Vinit] = w.Concat("ое")
	default:
		cCases[cases.Vinit] = w.Concat("ую")
	}

	return cCases, nil
}
