package dogstatsd

type stringInterner struct {
	strings map[string]string
}

func newStringInterner() *stringInterner {
	return &stringInterner{
		strings: make(map[string]string),
	}
}

func (i *stringInterner) LoadOrStore(key []byte) string {
	if s, found := i.strings[string(key)]; found {
		return s
	}
	s := string(key)
	i.strings[s] = s
	return s
}
