package russian

import "github.com/dshipenok/gomorphos/str"

/**
 * @var string[] Все гласные
 */
var Vowels = str.NewWordSet([]string{
	"а",
	"е",
	"ё",
	"и",
	"о",
	"у",
	"ы",
	"э",
	"ю",
	"я",
})

/**
 * @var string[] Все согласные
 */
var Consonants = []string{
	"б",
	"в",
	"г",
	"д",
	"ж",
	"з",
	"й",
	"к",
	"л",
	"м",
	"н",
	"п",
	"р",
	"с",
	"т",
	"ф",
	"х",
	"ц",
	"ч",
	"ш",
	"щ",
}

var ConsonantsSet = str.NewWordSet(Consonants)

/**
* @var string[] Звонкие согласные
 */
var SonorousConsonants = str.NewWordSet([]string{"б", "в", "г", "д", "з", "ж", "л", "м", "н", "р"})

/**
* @var string[] Глухие согласные
 */
var DeafConsonants = str.NewWordSet([]string{"п", "ф", "к", "т", "с", "ш", "х", "ч", "щ"})
