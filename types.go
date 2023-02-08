package openskill

// Rating represents a player's skill in the particular set rankings
type Rating struct {
	AveragePlayerSkill     float64
	SkillUncertaintyDegree float64
}

// Team is nothing more than a collection of Ratings
type Team []*Rating

// Gamma represents a function to help to reduce the variance of how much the skill uncertainty degree can change.
// It is defined by user and helps to correct cases in which the player skill can assume huge negative or positive
// values with an uncertainty value near zero
type Gamma func(adjustedTeamUncertainty float64, amountOfTeams int64, averageTeamSkill float64, teamUncertaintySquared float64, team *Team, teamRanking int64) float64

// Model represents the kind of ranking model is chosen top be used
type Model func(teams []Team, options *Options) []Team

// Options contains the values provided of the constants used by the rating system,
// plus optional rankings and scores of the teams to be rated, plus a optional paramenter
// that stops a player true skill rating from going after a victory, which can feel unfair.
type Options struct {
	NormalizedPlayerSkill      *float64
	AveragePlayerSkill         *float64
	SkillUncertaintyDegree     *float64
	SmallPositive              *float64
	GammaFunction              *Gamma
	VarianceForTeamPerformance *float64
	Model                      *Model
	Rankings                   []int64
	Scores                     []int64     // This was an executive decision to only allow integer values, or else I would be using maps for a lot of the operations
	weight                     [][]float64 // There is currently no provided model that takes into account of partial play or contribution, keeping this field because the original project also do so
	tau                        *float64
	PreventUncertaintyIncrease *bool
}
