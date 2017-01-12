package util

import (
	"math"
)

// CalculateAverage returns the Average Grade
func CalculateAverage(grades map[string]Module) float64 {
	var teamprojekt Module
	// var thesis Module
	var research []Module
	var csFundamentals []Module
	var ieFundamentals []Module
	var mmmFundamentals []Module
	var specs []Module
	var grade float64
	var credits float64
	for _, v := range grades {
		// dividing the grades into the groups that they belong to
		if v.ExamID == 420500 {
			teamprojekt = v
		} else if v.ExamID == 420000 {
			research = append(research, v)
		} else if 400499 < v.ExamID && v.ExamID < 400600 {
			// fmt.Println("CS Fundamental " + v.Name)
			csFundamentals = append(csFundamentals, v)
		} else if 400599 < v.ExamID && v.ExamID < 400700 {
			// fmt.Println("CS Spec " + v.Name)
			specs = append(specs, v)
		} else if 410499 < v.ExamID && v.ExamID < 410600 {
			// fmt.Println("IE Fundamental " + v.Name)
			ieFundamentals = append(ieFundamentals, v)
		} else if 410599 < v.ExamID && v.ExamID < 410700 {
			// fmt.Println("IE Spec " + v.Name)
			specs = append(specs, v)
		} else if 140000 < v.ExamID && v.ExamID < 140900 {
			// fmt.Println("BWL Fundamental " + v.Name)
			mmmFundamentals = append(mmmFundamentals, v)
		} else if 170000 < v.ExamID && v.ExamID < 170999 {
			specs = append(specs, v)
		}
	}
	grade += teamprojekt.Grade * teamprojekt.Bonus
	credits += teamprojekt.Bonus
	c, a := AverageGroup(research...)
	credits += c
	grade += a * c
	c, a = AverageGroup(csFundamentals...)
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

func truncate(grade float64, dec int) float64 {
	decimals := math.Pow10(dec)
	return math.Floor(grade*decimals) / decimals

}
