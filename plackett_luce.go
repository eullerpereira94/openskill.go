package openskill

import (
	"math"

	"github.com/samber/lo"
)

// PlackettLuce represents the Plackett-Luce ranking model, which is the generalized version of the Bradley-Terry model
func PlackettLuce(game []Team, options *Options) []Team {
	epsilon := epsilon(options)
	teamRatings := teamRatings(options)(game)
	c := utilC(options)(teamRatings)
	sumQ := utilSumQ(teamRatings, c)
	a := utilA(teamRatings)
	gamma := gamma(options)

	return lo.Map(teamRatings, func(item teamRating, index int) Team {
		iMuOverCe := math.Exp(item.TeamMu / c)

		type _sums struct {
			omegaSum float64
			deltaSum float64
		}

		filteredRatings := lo.Filter(teamRatings, func(localItem teamRating, localIndex int) bool {
			return localItem.Rank <= item.Rank
		})

		sums := lo.Reduce(filteredRatings, func(agg _sums, item teamRating, localIndex int) _sums {
			quotient := iMuOverCe / sumQ[localIndex]

			agg.omegaSum = agg.omegaSum + lo.Ternary(index == localIndex, 1-quotient, -quotient)/float64(a[localIndex])
			agg.deltaSum = agg.deltaSum + (quotient*(1-quotient))/float64(a[localIndex])

			return agg
		}, _sums{omegaSum: 0, deltaSum: 0})

		iGamma := gamma(c, int64(len(teamRatings)), item.TeamMu, item.TeamSigmaSq, item.Team, item.Rank)
		iOmega := sums.omegaSum * (item.TeamSigmaSq / c)
		iDelta := iGamma * sums.deltaSum * (item.TeamSigmaSq / math.Pow(c, 2))

		result := lo.Map([]*Rating(*item.Team), func(finalItem *Rating, index int) *Rating {
			return &Rating{
				AveragePlayerSkill:     finalItem.AveragePlayerSkill + (math.Pow(finalItem.SkillUncertaintyDegree, 2)/item.TeamSigmaSq)*iOmega,
				SkillUncertaintyDegree: finalItem.SkillUncertaintyDegree * math.Sqrt(math.Max(1-(math.Pow(finalItem.SkillUncertaintyDegree, 2)/item.TeamSigmaSq)*iDelta, epsilon)),
			}
		})

		return Team(result)
	})
}
