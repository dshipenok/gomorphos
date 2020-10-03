package declension

import (
	"fmt"

	"github.com/dshipenok/gomorphos/russian"
	"github.com/dshipenok/gomorphos/russian/cases"
	"github.com/dshipenok/gomorphos/russian/gender"
	"github.com/dshipenok/gomorphos/str"
)

/**
 * Проверка, изменяемое ли слово.
 * @param string $word Слово для проверки
 * @param bool $animateness Признак одушевленности
 * @return bool
 */
func IsMutable(w str.Word, animateness bool) bool {
	w = w.Lower()
	if w.EndsWith(-1, "у", "и", "е", "о", "ю") || immutableWords.Has(w) {
		return false
	}
	return true
}

/**
 * Определение рода существительного.
 * @param string $word
 * @return string
 */
func DetectGender(w str.Word) gender.Gender {
	w = w.Lower()
	last := w.LastChars(1)
	// пытаемся угадать род объекта, хотя бы примерно, чтобы правильно склонять
	if w.LastChars(2) == "мя" || w.EndsWith(1, "о", "е", "и", "у") {
		return gender.Neuter
	}

	if w.EndsWith(1, "а", "я") ||
		(last == "ь" &&
			!masculineWithSoft.Has(w) &&
			!masculineWithSoftAndRunAwayVowels.Has(w)) {
		return gender.Female
	}

	return gender.Male
}

/**
 * Определение склонения (по школьной программе) существительного.
 * @param string $word
 * @param bool $animateness
 * @return int
 */
func GetDeclension(w str.Word, animateness bool) int {
	w = w.Lower()
	last := w.LastChars(1)
	if abnormalExceptions.HasAny(w) {
		return SecondDeclension
	}

	if w.EndsWith(1, "а", "я") && !w.EndsWith(2, "мя") {
		return FirstDeclension
	} else if russian.IsConsonant(last) || w.EndsWith(1, "о", "е", "ё") ||
		(last == "ь" && russian.IsConsonant(w.Chars(-2, -1)) &&
			!russian.IsHissingConsonant(w.Chars(-2, -1)) &&
			(masculineWithSoft.Has(w) || masculineWithSoftAndRunAwayVowels.Has(w))) {
		return SecondDeclension
	}
	return ThirdDeclension
}

/**
 * Получение слова во всех 6 падежах.
 * @param string $word
 * @param bool $animateness Признак одушевлённости
 * @return string[]
 * @phpstan-return array<string, string>
 */
func GetCases(w str.Word, animateness bool) map[cases.Case]string {
	w = w.Lower()

	// Адъективное склонение (Сущ, образованные от прилагательных и причастий) - прохожий, существительное
	if russian.IsAdjectiveNoun(w) {
		result, _ := DeclinateAdjective(w, animateness)
		return result
	}

	// Субстантивное склонение (существительные)
	if immutableWords.Has(w) {
		return map[cases.Case]string{
			cases.Imenit:  w.String(),
			cases.Rodit:   w.String(),
			cases.Dat:     w.String(),
			cases.Vinit:   w.String(),
			cases.Tvorit:  w.String(),
			cases.Predloj: w.String(),
		}
	}

	abnormalList := abnormalExceptions.SliceOf(w)
	if abnormalExceptions.Has(w) {
		if len(abnormalList) > 0 {
			result := map[cases.Case]string{}
			values := abnormalExceptions.SliceOf(w)
			for ind, pad := range []cases.Case{cases.Imenit, cases.Rodit, cases.Dat, cases.Vinit, cases.Tvorit, cases.Predloj} {
				result[pad] = values[ind]
			}
			return result
		} else {
			prefix := w.Chars(0, -1)
			return map[cases.Case]string{
				cases.Imenit:  string(w),
				cases.Rodit:   prefix + "ени",
				cases.Dat:     prefix + "ени",
				cases.Vinit:   string(w),
				cases.Tvorit:  prefix + "енем",
				cases.Predloj: prefix + "ени",
			}
		}
	}

	switch GetDeclension(w, false) {
	case FirstDeclension:
		return DeclinateFirstDeclension(w)
	case SecondDeclension:
		return DeclinateSecondDeclension(w, animateness)
	case ThirdDeclension:
		return DeclinateThirdDeclension(w)
	}

	return nil // should never reach it
}

/**
 * @param string $word
 * @param string $last
 * @param string $prefix
 * @return string
 */
func GetPredCaseOf12Declensions(w str.Word, last, prefix string) string {
	if w.EndsWith(2, "ий", "ие") {
		if last == "ё" {
			return prefix + "е"
		} else {
			return prefix + "и"
		}
	} else {
		return prefix + "е"
	}
}

/**
 * Получение всех форм слова первого склонения.
 * @param string $word
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateFirstDeclension(w str.Word) (forms map[cases.Case]string) {
	w = w.Lower()
	prefix := w.Lower().Chars(0, -1)
	last := w.LastChars(1)
	softLast := russian.CheckLastConsonantSoftness(w)
	forms = map[cases.Case]string{
		cases.Imenit: string(w),
	}

	// RODIT
	tmpSoftLast := softLast || w.EndsWith(1, "г", "к", "х")
	forms[cases.Rodit] = russian.ChooseVowelAfterConsonant(last, tmpSoftLast, prefix+"и", prefix+"ы")

	// DAT
	forms[cases.Dat] = GetPredCaseOf12Declensions(w, last, prefix)

	// VINIT
	tmpSoftLast = softLast && w.LastChars(1) != "ч"
	forms[cases.Vinit] = russian.ChooseVowelAfterConsonant(last, tmpSoftLast, prefix+"ю", prefix+"у")

	// TVORIT
	if last == "ь" {
		forms[cases.Tvorit] = prefix + "ой"
	} else {
		forms[cases.Tvorit] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"ей", prefix+"ой")
	}

	// 	if ($last == 'й' || (static::isConsonant($last) && !static::isHissingConsonant($last)) || static::checkLastConsonantSoftness($word))
	// 	$forms[Cases::TVORIT] = $prefix.'ей';
	// else
	// 	$forms[Cases::TVORIT] = $prefix.'ой'; # http://morpher.ru/Russian/Spelling.aspx#sibilant

	// PREDLOJ the same as DAT
	forms[cases.Predloj] = forms[cases.Dat]
	return forms
}

/**
 * Получение всех форм слова третьего склонения.
 * @param string $word
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateThirdDeclension(w str.Word) map[cases.Case]string {
	w = w.Lower()
	prefix := w.Chars(0, -1)
	return map[cases.Case]string{
		cases.Imenit:  w.String(),
		cases.Rodit:   prefix + "и",
		cases.Dat:     prefix + "и",
		cases.Vinit:   w.String(),
		cases.Tvorit:  prefix + "ью",
		cases.Predloj: prefix + "и",
	}
}

/**
 * Получение всех форм слова второго склонения.
 * @param string $word
 * @param bool $animateness
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateSecondDeclension(w str.Word, animateness bool) map[cases.Case]string {
	w = w.Lower()
	lastWord := w.LastCharsWord(1)
	last := lastWord.String()
	prelast := w.Chars(-2, -1)
	softLast := last == "й" ||
		(lastWord.OneOf("ь", "е", "ё", "ю", "я") &&
			((russian.IsConsonant(prelast) &&
				!russian.IsHissingConsonant(prelast)) ||
				prelast == "и"))
	prefix := GetPrefixOfSecondDeclension(w, lastWord)
	forms := cases.NewCases()
	forms[cases.Imenit] = w.String()

	// RODIT
	forms[cases.Rodit] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"я", prefix+"а")

	// DAT
	forms[cases.Dat] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"ю", prefix+"у")

	// VINIT
	if lastWord.OneOf("о", "е", "ё") {
		forms[cases.Vinit] = w.String()
	} else {
		forms[cases.Vinit] = russian.GetVinitCaseByAnimateness(forms, animateness)
	}

	// TVORIT
	// if ($last == 'ь')
	// 	$forms[Cases::TVORIT] = $prefix.'ом';
	// else if ($last == 'й' || (static::isConsonant($last) && !static::isHissingConsonant($last)))
	// 	$forms[Cases::TVORIT] = $prefix.'ем';
	// else
	// 	$forms[Cases::TVORIT] = $prefix.'ом'; # http://morpher.ru/Russian/Spelling.aspx#sibilant
	if (russian.IsHissingConsonant(last) && last != "ш") ||
		(lastWord.OneOf("ь", "е", "ё", "ю", "я") && russian.IsHissingConsonant(w.Chars(-2, -1))) ||
		(last == "ц" && w.LastChars(2) != "ец") {
		forms[cases.Tvorit] = prefix + "ем"
	} else if lastWord.OneOf("й") || softLast {
		forms[cases.Tvorit] = prefix + "ем"
	} else {
		forms[cases.Tvorit] = prefix + "ом"
	}

	// PREDLOJ
	forms[cases.Predloj] = GetPredCaseOf12Declensions(w, last, prefix)

	return forms
}

/**
 * Склонение существительных, образованных от прилагательных и причастий.
 * Rules are from http://rusgram.narod.ru/1216-1231.html
 * @param string $word
 * @param bool $animateness
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateAdjective(w str.Word, animateness bool) (map[cases.Case]string, error) {
	prefix := w.Chars(0, -2)

	switch w.LastChars(2) {
	// Male adjectives
	case "ой", "ый":
		return map[cases.Case]string{
			cases.Imenit:  w.String(),
			cases.Rodit:   prefix + "ого",
			cases.Dat:     prefix + "ому",
			cases.Vinit:   w.String(),
			cases.Tvorit:  prefix + "ым",
			cases.Predloj: prefix + "ом",
		}, nil

	case "ий":
		return map[cases.Case]string{
			cases.Imenit:  w.String(),
			cases.Rodit:   prefix + "его",
			cases.Dat:     prefix + "ему",
			cases.Vinit:   prefix + "его",
			cases.Tvorit:  prefix + "им",
			cases.Predloj: prefix + "ем",
		}, nil

	// Neuter adjectives
	case "ое", "ее":
		prefix = w.Chars(0, -1)
		middle := "и"
		if w.Chars(-2, -1) == "о" {
			middle = "ы"
		}
		return map[cases.Case]string{
			cases.Imenit:  w.String(),
			cases.Rodit:   prefix + "го",
			cases.Dat:     prefix + "му",
			cases.Vinit:   w.String(),
			cases.Tvorit:  w.Chars(0, -2) + middle + "м",
			cases.Predloj: prefix + "м",
		}, nil

	// Female adjectives
	case "ая":
		ending := "ой"
		if russian.IsHissingConsonant(str.Word(prefix).LastChars(1)) {
			ending = "ей"
		}
		return map[cases.Case]string{
			cases.Imenit:  w.String(),
			cases.Rodit:   prefix + ending,
			cases.Dat:     prefix + ending,
			cases.Vinit:   prefix + "ую",
			cases.Tvorit:  prefix + ending,
			cases.Predloj: prefix + ending,
		}, nil
	default:
		return nil, fmt.Errorf("unknown ending %q", w.LastChars(2))
	}
}

/**
 * Получение одной формы слова (падежа).
 * @param string $word Слово
 * @param string $case Падеж
 * @param bool $animateness Признак одушевленности
 * @return string
 * @throws \Exception
 */
func GetCase(w str.Word, wCase string, animateness bool) string {
	cCase := cases.CanonizeCase(wCase)
	forms := GetCases(w, animateness)
	return forms[cCase]
}

/**
 * @param string $word
 * @param string $last
 * @return string
 */
func GetPrefixOfSecondDeclension(w, last str.Word) string {
	var prefix string
	// слова с бегающей гласной в корне
	if masculineWithSoftAndRunAwayVowels.Has(w) {
		prefix = w.Chars(0, -3) + w.Chars(-2, -1)
	} else if last.OneOf("о", "е", "ё", "ь", "й") {
		prefix = w.Chars(0, -1)
	} else if w.LastChars(2) == "ок" && w.Len() > 3 {
		// уменьшительные формы слов (котенок) и слова с суффиксом ок
		prefix = w.Chars(0, -2) + "к"
	} else if w.LastChars(3) == "бец" && w.Len() > 4 {
		// слова с суффиксом бец
		prefix = w.Chars(0, -3) + "бц"
	} else {
		prefix = w.String()
	}
	return prefix
}
