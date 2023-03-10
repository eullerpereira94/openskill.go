package openskill

func z(options *Options) float64 {
	if options != nil && options.StandardizedPlayerSkill != nil {
		return *options.StandardizedPlayerSkill
	}
	return 3
}

func mu(options *Options) float64 {
	if options != nil && options.AveragePlayerSkill != nil {
		return *options.AveragePlayerSkill
	}
	return 25
}

// this tau function is kept here for the sake of keeping a similar structure to the original code, even though is unused
func tau(options *Options) float64 {
	if options != nil && options.Tau != nil {
		return *options.Tau
	}
	return mu(options) / 300
}

func sigma(options *Options) float64 {
	if options != nil && options.SkillUncertaintyDegree != nil {
		return *options.SkillUncertaintyDegree
	}
	return mu(options) / z(options)
}

func epsilon(options *Options) float64 {
	if options != nil && options.SmallPositive != nil {
		return *options.SmallPositive
	}
	return 0.0001
}

func beta(options *Options) float64 {
	return sigma(options) / 2
}

func betaSq(options *Options) float64 {
	if options != nil && options.VarianceForTeamPerformance != nil {
		return *options.VarianceForTeamPerformance
	}

	beta := beta(options)

	return beta * beta
}
