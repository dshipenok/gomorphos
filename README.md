# gomorphos
Порт [Morphos](https://github.com/wapmorgan/Morphos) c PHP на Golang.

На данный момент реализованы только падежи и склонения существительных.
Тестов пока мало, очень вероятно, что код содержит ошибки.

## Примеры использования

```golang
// Получение склонения
func Test_GetDeclension(t *testing.T) {
	decls := GetDeclension(str.Word("лень"), false)

	assert.EqualValues(t, 3, decls)
}

// Получение слова в нужном падеже
func Test_GetCase(t *testing.T) {
	casedStr := GetCase(str.Word("кухня"), "винительный", false)

	assert.EqualValues(t, "кухню", casedStr)
}

// Получение всех падежей слова
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
	}

	for _, tst := range tests {
		t.Run(tst.Word, func(t *testing.T) {
			w := str.Word(tst.Word)

			cases := GetCases(w, false)

			assert.Equal(t, tst.Cases, cases)
		})
	}
}
```
