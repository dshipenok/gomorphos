package str

type WordSet map[string]struct{}

// NewWordSet constructor
func NewWordSet(strs []string) WordSet {
	ws := WordSet{}
	for _, s := range strs {
		ws[s] = struct{}{}
	}
	return ws
}

func (ws WordSet) Has(w Word) bool {
	_, has := ws[string(w)]
	return has
}

func (ws WordSet) HasStr(s string) bool {
	_, has := ws[s]
	return has
}
