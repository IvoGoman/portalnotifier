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
type ByName []Module

func (md ByName) Len() int {
	return len(md)
}
func (md ByName) Less(i, j int) bool {
	return md[i].Name < md[j].Name
}
func (md ByName) Swap(i, j int) {
	md[i], md[j] = md[j], md[i]
}
