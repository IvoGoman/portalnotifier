package util

import (
	"math"
)

// CalculateAverage returns the Average Grade
// ExamIDs are assigned to groups with hardcoded rules
func CalculateAverage(grades map[string]Module) float64 {
	var teamprojekt Module
	// var thesis Module
	// var research []Module
	var csFundamentals []Module
	var ieFundamentals []Module
	var mmmFundamentals []Module
	var specs []Module
	var grade float64
	var credits float64
	for _, v := range grades {
		// Teamproject
		if v.ExamID == 420500 {
			teamprojekt = v
			// Scientific Research Seminar
		} else if v.ExamID == 420000 {
			specs = append(specs, v)
			// Fundamental Course CS
		} else if 400499 < v.ExamID && v.ExamID < 400600 {
			csFundamentals = append(csFundamentals, v)
			// Specialization Course CS
		} else if 400599 < v.ExamID && v.ExamID < 400700 {
			specs = append(specs, v)
			// Fundamental Course IE
		} else if 410499 < v.ExamID && v.ExamID < 410600 {
			ieFundamentals = append(ieFundamentals, v)
			// Specialization Course IE
		} else if 410599 < v.ExamID && v.ExamID < 410700 {
			specs = append(specs, v)
			// Fundamental Course MMM
		} else if 140000 < v.ExamID && v.ExamID < 140900 {
			mmmFundamentals = append(mmmFundamentals, v)
			// Specialization Course Area IS
		} else if 170000 < v.ExamID && v.ExamID < 170999 {
			specs = append(specs, v)
			// Specialization Course Seminar
		} else if 430000 < v.ExamID && v.ExamID < 430999 {
			specs = append(specs, v)
		}
	}
	grade += teamprojekt.Grade * teamprojekt.Bonus
	credits += teamprojekt.Bonus
	c, a := AverageGroup(csFundamentals...)
	credits += c
	grade += a * c
	c, a = AverageGroup(ieFundamentals...)
	credits += c
	grade += a * c
	c, a = AverageGroup(mmmFundamentals...)
	credits += c
	grade += a * c
	c, a = AverageGroup(specs...)
	credits += c
	grade += a * c
	return truncate(grade/credits, 2)
}

// AverageGroup returns the average grade per module group
func AverageGroup(grades ...Module) (credits, average float64) {
	var grade float64
	for _, v := range grades {
		credits += v.Bonus
		grade += v.Bonus * v.Grade
	}
	return credits, truncate(grade/credits, 1)
}

// truncates float value
func truncate(grade float64, dec int) float64 {
	decimals := math.Pow10(dec)
	return math.Floor(grade*decimals) / decimals

}
