package openskill

import (
	"math"

	"github.com/samber/lo"
)

// PredictDraw returns the probability that a set of teams will tie based on their rating.
// If there is only one team, the function will return 1, and if there is no teams, it will
// return -1
func PredictDraw(teams []Team, options *Options) float64 {
	m := len(teams)

	if m <= 0 {
		return -1
	}

	if m == 1 {
		return 1
	}

	n := float64(m)

	teamRatings := teamRatings(options)(teams)
	beta := beta(options)
	betaSq := betaSq(options)

	denom := (n * (n - 1)) / (lo.Ternary(n > 2, 1.0, 2.0))

	var preFlattening [][]*Rating

	for _, v := range teams {
		var team []*Rating

		for _, w := range v {
			u := *w

			team = append(team, &u)
		}

		preFlattening = append(preFlattening, team)
	}

	drawMargin := math.Sqrt(float64(len(lo.Flatten(preFlattening)))) * beta * ppf((1+1/n)/2)

	processedRatings := lo.Map(teamRatings, func(item *teamRating, index int) []float64 {
		filteredRatings := lo.Filter(teamRatings, func(localItem *teamRating, localIndex int) bool {
			return localIndex != index
		})

		return lo.Map(filteredRatings, func(localItem *teamRating, localIndex int) float64 {
			sigmaBar := math.Sqrt(n*betaSq + math.Pow(item.TeamSigmaSq, 2) + math.Pow(localItem.TeamSigmaSq, 2))
			return cdf((drawMargin-item.TeamMu+localItem.TeamMu)/sigmaBar) - cdf((item.TeamMu-localItem.TeamMu-drawMargin)/sigmaBar)
		})
	})

	sum := lo.Sum(lo.Flatten(processedRatings))

	return math.Abs(sum) / denom
}
