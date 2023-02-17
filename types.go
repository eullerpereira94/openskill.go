package openskill

// Rating represents a player's skill in the particular set rankings
type Rating struct {
	// AveragePlayerSkill represents the average
	// value for a player skill. It normally is
	// the value in the middle of normal or
	// logistic distribution.
	AveragePlayerSkill float64

	// SkillUncertaintyDegree represents the amount
	// of uncertainty of the skill of a player.
	// This value is used to set the bounds of the
	// probalistic distribution.
	SkillUncertaintyDegree float64
}

// Team is nothing more than a collection of Ratings
type Team []*Rating

// Gamma represents a function to help to reduce the variance of how much the skill uncertainty degree can change.
// It is defined by user and helps to correct cases in which the player skill can assume huge negative or positive
// values with an uncertainty value near zero. Said function returns a number to adjust the skill uncertainty.
type Gamma func(adjustedTeamUncertainty float64, amountOfTeams int64, averageTeamSkill float64, teamUncertaintySquared float64, team *Team, teamRanking int64) float64

// Model represents the kind of ranking model is chosen to be used
type Model func(teams []Team, options *Options) []Team

// Options contains the values provided of the constants used by the rating system,
// plus optional rankings and scores of the teams to be rated, plus a optional paramenter
// that stops a player true skill rating from going after a victory, which can feel unfair.
type Options struct {
	// StandardizedPlayerSkill is a constant that is set in way that, the following assertion
	// is always true:
	//
	// StandardizedPlayerSkill == (Ordinal(Rating) - Rating.AveragePlayerSkill) / Rating.SkillUncertaintyDegree
	//
	// The bigger this constant is, the smaller the degree of uncertainty relative to average skill.
	// When not set,it defaults to 3.
	StandardizedPlayerSkill *float64

	// AveragePlayerSkill represents the default value of a player average skill level.
	// When not set, it defaults to 25.
	AveragePlayerSkill *float64

	// SkillUncertaintyDegree represents the default value of uncertainty for a player skill.
	// When not set, it defaults to Options.AveragePlayerSkill / Options.NormalizedPlayerSkill
	SkillUncertaintyDegree *float64

	// SmallPositive is a value to use when a ranking model tries to update a player uncertainty
	// with a negative value. It is rarely needed to do such substitutions.
	// When not set, it defaults to 0.001
	SmallPositive *float64

	// GammaFunction is a pointer to a function that is used to adjust
	// how much SkillUncertaintyDegree can vary. When not set, it defaults
	// to an internal implementation that mirrors the one found on item 6.1 of the
	// Weng-Lin paper for the Packett-Luce model.
	GammaFunction *Gamma

	// VarianceForTeamPerformance represents a constant to adjust the value of a team
	// performance. When not set, it defaults to (Options.SkillUncertaintyDegree / 2) ^ 2.
	// The default value for this constant takes into consideration if Options.SkillUncertaintyDegree is set or if
	// either Options.AveragePlayerSkill or Options.NormalizedPlayerSkill are set.
	VarianceForTeamPerformance *float64

	// Model represents the current model of ranking used. When not set, it defaults to Plackett-Luce.
	Model *Model

	// Rankings is a optional slice of rankings that is used when provided order of the teams for the
	// Rate function differs from the actual order of rankings. Other use for this field is to indicate
	// when ties happened after a competition.
	Rankings []int64

	// Scores is slice of the scores of the teams after competing. In this implementation
	// is a slice of integer values because otherwise, a lot of the internal operations
	// would use maps. This field doesn't need to be initialized with values.
	Scores []int64

	weight [][]float64 // There is currently no provided model that takes into account of partial play or contribution, keeping this field because the original project also do so

	// Tau is a value that prevents the uncertainty to drop to a value that is too low.
	// Setting this constant, allows the rating to stay pliable even after many games.
	// A suggested value for this constant is Options.AveragePlayerSkill / 300.
	Tau *float64

	// PreventUncertaintyIncrease is an optional boolean value that, if it is set, and if Options.Tau is set,
	// prevents the uncertainty value to increase, thus stopping the fringe case when the Ordinal of player
	// rating decrease after a victory, which can feel unfair.
	PreventUncertaintyIncrease *bool
}
