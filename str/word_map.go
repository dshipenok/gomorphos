package str

type WordMap struct {
	Map map[string][]string

	set map[string]struct{}
}

// NewWordMap constructor
func NewWordMap(m map[string][]string) WordMap {
	set := map[string]struct{}{}

	for k, slice := range m {
		set[k] = struct{}{}
		for _, s := range slice {
			set[s] = struct{}{}
		}
	}

	return WordMap{
		Map: m,
		set: set,
	}
}

func (wm WordMap) Has(w Word) bool {
	_, has := wm.Map[string(w)]
	return has
}

func (wm WordMap) HasAny(w Word) bool {
	_, has := wm.set[string(w)]
	return has
}

func (wm WordMap) SliceOf(w Word) []string {
	return wm.Map[string(w)]
}
