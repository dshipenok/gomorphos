package russian

import (
	"strings"

	"github.com/dshipenok/gomorphos/russian/cases"
	"github.com/dshipenok/gomorphos/str"
)

/**
 * Выбор окончания в зависимости от мягкости
 *
 * @param string $last
 * @param bool $softLast
 * @param string $afterSoft
 * @param string $afterHard
 *
 * @return string
 */
func ChooseVowelAfterConsonant(last string, softLast bool, afterSoft, afterHard string) string {
	if last != "щ" && /*static::isVelarConsonant($last) ||*/ softLast {
		return afterSoft
	}

	return afterHard
}

/**
 * @param string $string
 * @param string[] $chars
 * @return string|false
 */
func FindLastPositionForOneOfChars(w str.Word, chars []string) str.Word {
	lastPosition := -1
	for _, ch := range chars {
		pos := w.LastIndex(ch)
		if pos < 0 {
			continue
		}
		if pos > lastPosition {
			lastPosition = pos
		}
	}
	if lastPosition >= 0 {
		return w.SubWord(lastPosition)
	}
	return str.Word{}
}

/**
 * Проверка мягкости последней согласной
 * @param string $word
 * @return bool
 */
func CheckLastConsonantSoftness(w str.Word) bool {
	substring := FindLastPositionForOneOfChars(w.Lower(), Consonants)
	if len(substring) == 0 {
		return false
	}
	if substring.SliceWord(0, 1).EndsWith(1, "й", "ч", "щ", "ш") { // always soft consonants
		return true
	}
	if substring.Len() > 1 && substring.SliceWord(1, 2).OneOf("е", "ё", "и", "ю", "я", "ь") { // consonants are soft if they are trailed with these vowels
		return true
	}
	return false
}

/**
* Проверка гласной
* @param string $char
* @return bool
 */
func IsVowel(ch string) bool {
	return Vowels.HasStr(ch)
}

/**
* Проверка согласной
* @param string $char
* @return bool
 */
func IsConsonant(ch string) bool {
	return ConsonantsSet.HasStr(ch)
}

/**
* Проверка звонкости согласной
* @param string $char
* @return bool
 */
func IsSonorousConsonant(ch string) bool {
	return SonorousConsonants.HasStr(ch)
}

/**
* Проверка глухости согласной
* @param string $char
* @return bool
 */
func IsDeafConsonant(ch string) bool {
	return DeafConsonants.HasStr(ch)
}

/**
* Щипящая ли согласная
* @param string $consonant
* @return bool
 */
func IsHissingConsonant(consonant string) bool {
	consonant = strings.ToLower(consonant)
	for _, c := range []string{"ж", "ш", "ч", "щ"} {
		if consonant == c {
			return true
		}
	}
	return false
}

/**
* Проверка на велярность согласной
* @param string $consonant
* @return bool
 */
func IsVelarConsonant(consonant string) bool {
	consonant = strings.ToLower(consonant)
	for _, c := range []string{"г", "к", "х"} {
		if consonant == c {
			return true
		}
	}
	return false
}

/**
 * Проверяет, является ли существительно адъективным существительным
 * @param string $noun Существительное
 * @return bool
 */
func IsAdjectiveNoun(noun str.Word) bool {
	return noun.EndsWith(2, "ой", "ий", "ый", "ая", "ое", "ее") &&
		!noun.OneOf("гений", "комментарий")
}

func GetVinitCaseByAnimateness(forms map[cases.Case]string, animate bool) string {
	if animate {
		return forms[cases.Rodit]
	}

	return forms[cases.Imenit]
}

/**
* Проверка мягкости последней согласной, за исключением Н
* @param string $word
* @return bool
 */
func CheckBaseLastConsonantSoftness(w str.Word) bool {
	substring := FindLastPositionForOneOfChars(w.Lower(), ConsonantsAdj)
	if !substring.Empty() {
		if substring.SliceWord(0, 1).OneOf("й", "ч", "щ", "ш") { // always soft consonants
			return true
		}
		if substring.Len() > 1 && substring.SliceWord(1, 2).OneOf("е", "ё", "и", "ю", "я", "ь") { // consonants are soft if they are trailed with these vowels
			return true
		}
	}
	return false
}
