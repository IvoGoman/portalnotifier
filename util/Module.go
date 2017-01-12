package util

// Module a struct which corresponds to one module
type Module struct {
	ExamID          int64
	Semester        string
	TryCountExam    int64
	Date            string
	Name            string
	Prof            string
	Form            string
	Grade           float64
	Bonus           float64
	Status          string
	TryCountStudent string
}

// ByName sorts the Modules by their name
type ByName []Module

// ModuleMapToArray returns an Array that can be used for sorting the modules
func ModuleMapToArray(moduleMap map[string]Module) (array []Module) {
	// array := make([]Module, 0)
	for _, m := range moduleMap {
		array = append(array, m)
	}
	return
}

func (md ByName) Len() int {
	return len(md)
}
func (md ByName) Less(i, j int) bool {
	return md[i].Name < md[j].Name
}
func (md ByName) Swap(i, j int) {
	md[i], md[j] = md[j], md[i]
}
