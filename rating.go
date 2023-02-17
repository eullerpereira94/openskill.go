package openskill

// NewRatingParams represents the values that help initilize a rating.
type NewRatingParams struct {
	AveragePlayerSkill     float64
	SkillUncertaintyDegree float64
}

// NewRating creates a new Rating, with optional initializing rating values and a optional set of default constants.
// Useful when not creating ratings in a vacuum.
func NewRating(init *NewRatingParams, options *Options) *Rating {
	var (
		_mu, _sigma float64
	)
	if init != nil {
		_mu = init.AveragePlayerSkill
		_sigma = init.SkillUncertaintyDegree
	} else {
		_mu = mu(options)
		_sigma = sigma(options)
	}
	return &Rating{
		AveragePlayerSkill:     _mu,
		SkillUncertaintyDegree: _sigma,
	}
}
