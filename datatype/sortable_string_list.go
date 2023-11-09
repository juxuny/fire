package datatype

type SortableStringList []string

func (t SortableStringList) Len() int {
	return len(t)
}

func (t SortableStringList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t SortableStringList) Less(i, j int) bool {
	return t[i] < t[j]
}
