package openskill

// Ordinal takes a Rating and returns a value that has 99.7% of representing the player true skill rating
func Ordinal(rating Rating, options *Options) float64 {
	return rating.AveragePlayerSkill - (z(options) * rating.SkillUncertaintyDegree)
}
