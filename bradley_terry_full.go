package openskill

import (
	"math"

	"github.com/samber/lo"
)

// BradleyTerryFull is a implementation of the Bradley-Terry ranking model
// that uses full pairing. The Bradley-Terry model uses logistic distribution
// to properly rank the teams. It accepts the a slice with the team that are
// competing, plus an options parameter, with things such as scores and
// previous rankings. The function return a slice of teams that are properly ranked.
func BradleyTerryFull(game []Team, options *Options) []Team {
	epsilon := epsilon(options)
	tbs := betaSq(options) * 2
	_gamma := gamma(options)

	teamRatings := teamRatings(options)(game)

	return lo.Map(teamRatings, func(item teamRating, index int) Team {
		var iMu, iSigmaSq, iTeam, iRank = item.TeamMu, item.TeamSigmaSq, item.Team, item.Rank

		type _sums struct {
			omegaSum float64
			deltaSum float64
		}

		filteredRatings := lo.Filter(teamRatings, func(localItem teamRating, localIndex int) bool {
			return localIndex != index
		})

		sums := lo.Reduce(filteredRatings, func(agg _sums, localItem teamRating, index int) _sums {
			var qMu, qSigmaSq, qRank = localItem.TeamMu, localItem.TeamSigmaSq, localItem.Rank

			ciq := math.Sqrt(iSigmaSq + qSigmaSq + tbs)
			piq := 1 / (1 + math.Exp((qMu-iMu)/ciq))

			sigSqToCiq := iSigmaSq / ciq

			iGamma := _gamma(ciq, int64(len(teamRatings)), item.TeamMu, item.TeamSigmaSq, item.Team, item.Rank)

			agg.omegaSum += sigSqToCiq * (score(qRank, iRank) - piq)
			agg.deltaSum += ((iGamma * sigSqToCiq) / ciq) * piq * (1 - piq)

			return agg
		}, _sums{omegaSum: 0, deltaSum: 0})

		result := lo.Map([]*Rating(*iTeam), func(finalItem *Rating, index int) *Rating {
			sigmaSq := math.Pow(finalItem.SkillUncertaintyDegree, 2)
			mu := finalItem.AveragePlayerSkill + (sigmaSq/iSigmaSq)*sums.omegaSum
			sigma := finalItem.SkillUncertaintyDegree * math.Sqrt(math.Max(1-(sigmaSq/iSigmaSq)*sums.deltaSum, epsilon))

			return &Rating{
				AveragePlayerSkill:     mu,
				SkillUncertaintyDegree: sigma,
			}
		})

		return Team(result)
	})
}
