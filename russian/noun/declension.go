package declension

import (
	"fmt"
	"strings"

	"github.com/dshipenok/gomorphos/russian"
	"github.com/dshipenok/gomorphos/str"
)

func newCases() map[Case]string {
	return make(map[Case]string, int(Predloj))
}

func GetVinitCaseByAnimateness(forms map[Case]string, animate bool) string {
	if animate {
		return forms[Rodit]
	}

	return forms[Imenit]
}

/**
 * @param string $case
 * @return string
 * @throws \Exception
 */
func canonizeCase(wCase string) Case {
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
func DetectGender(w str.Word) Gender {
	w = w.Lower()
	last := w.LastChars(1)
	// пытаемся угадать род объекта, хотя бы примерно, чтобы правильно склонять
	if w.LastChars(2) == "мя" || w.EndsWith(1, "о", "е", "и", "у") {
		return Neuter
	}

	if w.EndsWith(1, "а", "я") ||
		(last == "ь" &&
			!masculineWithSoft.Has(w) &&
			!masculineWithSoftAndRunAwayVowels.Has(w)) {
		return Female
	}

	return Male
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
func GetCases(w str.Word, animateness bool) map[Case]string {
	w = w.Lower()

	// Адъективное склонение (Сущ, образованные от прилагательных и причастий) - прохожий, существительное
	if russian.IsAdjectiveNoun(w) {
		result, _ := DeclinateAdjective(w, animateness)
		return result
	}

	// Субстантивное склонение (существительные)
	if immutableWords.Has(w) {
		return map[Case]string{
			Imenit:  w.String(),
			Rodit:   w.String(),
			Dat:     w.String(),
			Vinit:   w.String(),
			Tvorit:  w.String(),
			Predloj: w.String(),
		}
	}

	if abnormalExceptions.Has(w) {
		result := map[Case]string{}
		values := abnormalExceptions.SliceOf(w)
		for ind, pad := range []Case{Imenit, Rodit, Dat, Vinit, Tvorit, Predloj} {
			result[pad] = values[ind]
		}
		return result
	}

	if abnormalExceptions.Has(w) {
		prefix := w.Chars(0, -1)
		return map[Case]string{
			Imenit:  string(w),
			Rodit:   prefix + "ени",
			Dat:     prefix + "ени",
			Vinit:   string(w),
			Tvorit:  prefix + "енем",
			Predloj: prefix + "ени",
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
func DeclinateFirstDeclension(w str.Word) (forms map[Case]string) {
	w = w.Lower()
	prefix := w.Lower().Chars(0, -1)
	last := w.LastChars(1)
	softLast := russian.CheckLastConsonantSoftness(w)
	forms = map[Case]string{
		Imenit: string(w),
	}

	// RODIT
	tmpSoftLast := softLast || w.EndsWith(1, "г", "к", "х")
	forms[Rodit] = russian.ChooseVowelAfterConsonant(last, tmpSoftLast, prefix+"и", prefix+"ы")

	// DAT
	forms[Dat] = GetPredCaseOf12Declensions(w, last, prefix)

	// VINIT
	tmpSoftLast = softLast && w.LastChars(1) != "ч"
	forms[Vinit] = russian.ChooseVowelAfterConsonant(last, tmpSoftLast, prefix+"ю", prefix+"у")

	// TVORIT
	if last == "ь" {
		forms[Tvorit] = prefix + "ой"
	} else {
		forms[Tvorit] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"ей", prefix+"ой")
	}

	// 	if ($last == 'й' || (static::isConsonant($last) && !static::isHissingConsonant($last)) || static::checkLastConsonantSoftness($word))
	// 	$forms[Cases::TVORIT] = $prefix.'ей';
	// else
	// 	$forms[Cases::TVORIT] = $prefix.'ой'; # http://morpher.ru/Russian/Spelling.aspx#sibilant

	// PREDLOJ the same as DAT
	forms[Predloj] = forms[Dat]
	return forms
}

/**
 * Получение всех форм слова третьего склонения.
 * @param string $word
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateThirdDeclension(w str.Word) map[Case]string {
	w = w.Lower()
	prefix := w.Chars(0, -1)
	return map[Case]string{
		Imenit:  w.String(),
		Rodit:   prefix + "и",
		Dat:     prefix + "и",
		Vinit:   w.String(),
		Tvorit:  prefix + "ью",
		Predloj: prefix + "и",
	}
}

/**
 * Получение всех форм слова второго склонения.
 * @param string $word
 * @param bool $animateness
 * @return string[]
 * @phpstan-return array<string, string>
 */
func DeclinateSecondDeclension(w str.Word, animateness bool) map[Case]string {
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
	forms := newCases()
	forms[Imenit] = w.String()

	// RODIT
	forms[Rodit] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"я", prefix+"а")

	// DAT
	forms[Dat] = russian.ChooseVowelAfterConsonant(last, softLast, prefix+"ю", prefix+"у")

	// VINIT
	if lastWord.OneOf("о", "е", "ё") {
		forms[Vinit] = w.String()
	} else {
		forms[Vinit] = GetVinitCaseByAnimateness(forms, animateness)
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
		forms[Tvorit] = prefix + "ем"
	} else if lastWord.OneOf("й") || softLast {
		forms[Tvorit] = prefix + "ем"
	} else {
		forms[Tvorit] = prefix + "ом"
	}

	// PREDLOJ
	forms[Predloj] = GetPredCaseOf12Declensions(w, last, prefix)

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
func DeclinateAdjective(w str.Word, animateness bool) (map[Case]string, error) {
	prefix := w.Chars(0, -2)

	switch w.LastChars(2) {
	// Male adjectives
	case "ой", "ый":
		return map[Case]string{
			Imenit:  w.String(),
			Rodit:   prefix + "ого",
			Dat:     prefix + "ому",
			Vinit:   w.String(),
			Tvorit:  prefix + "ым",
			Predloj: prefix + "ом",
		}, nil

	case "ий":
		return map[Case]string{
			Imenit:  w.String(),
			Rodit:   prefix + "его",
			Dat:     prefix + "ему",
			Vinit:   prefix + "его",
			Tvorit:  prefix + "им",
			Predloj: prefix + "ем",
		}, nil

	// Neuter adjectives
	case "ое", "ее":
		prefix = w.Chars(0, -1)
		middle := "и"
		if w.Chars(-2, -1) == "о" {
			middle = "ы"
		}
		return map[Case]string{
			Imenit:  w.String(),
			Rodit:   prefix + "го",
			Dat:     prefix + "му",
			Vinit:   w.String(),
			Tvorit:  w.Chars(0, -2) + middle + "м",
			Predloj: prefix + "м",
		}, nil

	// Female adjectives
	case "ая":
		ending := "ой"
		if russian.IsHissingConsonant(str.Word(prefix).LastChars(1)) {
			ending = "ей"
		}
		return map[Case]string{
			Imenit:  w.String(),
			Rodit:   prefix + ending,
			Dat:     prefix + ending,
			Vinit:   prefix + "ую",
			Tvorit:  prefix + ending,
			Predloj: prefix + ending,
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
	cCase := canonizeCase(wCase)
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
