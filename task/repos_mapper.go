package task

import (
	"sort"

	"github.com/yuanjiecloud/fire/datatype"
)

// ReposMapper is a mapper from package name to local file path
type ReposMapper map[string]string

func NewReposMapper() ReposMapper {
	return make(ReposMapper)
}

func (t ReposMapper) Clone() ReposMapper {
	result := make(ReposMapper)
	for k, v := range t {
		result[k] = v
	}
	return result
}

func (t ReposMapper) OverridePatch(in ReposMapper) ReposMapper {
	result := t.Clone()
	for k, v := range in {
		result[k] = v
	}
	return result
}

func (t ReposMapper) MergeIgnoreDuplicated(in ReposMapper) ReposMapper {
	result := t.Clone()
	for k, v := range in {
		if _, b := result[k]; !b {
			result[k] = v
		}
	}
	return result
}

func (t ReposMapper) GetKeys() []string {
	var list datatype.SortableStringList
	for k := range t {
		list = append(list, k)
	}
	sort.Sort(list)
	return []string(list)
}
