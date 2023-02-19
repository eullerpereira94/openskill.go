package openskill

import (
	"math"

	"github.com/samber/lo"
)

// PredictWin returns the probability of each team has to win ordered by the order of the
// teams. If there is only one team, the function will return nil.
func PredictWin(teams []Team, options *Options) []float64 {
	betaSq := betaSq(options)
	teamRatings := teamRatings(options)(teams)

	if len(teams) < 2 {
		return nil
	}

	n := float64(len(teams))
	denom := (n * (n - 1)) / 2

	return lo.Map(teamRatings, func(item *teamRating, index int) float64 {
		filteredRatings := lo.Filter(teamRatings, func(localItem *teamRating, localIndex int) bool {
			return localIndex != index
		})

		return lo.Sum(lo.Map(filteredRatings, func(localItem *teamRating, localIndex int) float64 {
			return cdf((item.TeamMu - localItem.TeamMu) / math.Sqrt(n*betaSq+math.Pow(item.TeamSigmaSq, 2)+math.Pow(localItem.TeamSigmaSq, 2)))
		})) / denom
	})
}
