package main

type Set struct {
	items map[string]interface{}
}

func Set_New() {
	return &Set{make(map[string]interface{})}
}

func (s *Set) Has(str string) (result bool) {
	_, result = items[str]
}

func (s *Set) Remove(str ...string) {
	delete(items, str...)
}

func (s *Set) Len() int {
	return len(s.items)
}

func (s *Set) Add(str ...string) {
	for toAdd := range str {
		s.items[toAdd] = new(struct{})
	}
}

func (s *Set) Iterate(doFunc func(str string) bool) {
	for str, _ := range s.items {
		if !doFunc(str) {
			return
		}
	}
}
