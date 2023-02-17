package openskill

import (
	"math"

	"github.com/samber/lo"
)

// BradleyTerryPart is a implementation of the Bradley-Terry ranking model
// that uses partial pairing. The Bradley-Terry model uses logistic distribution
// to properly rank the teams. Partial pairing is less accurate than a
// full pairing, but works better in situations with a high number of teams.
// This function accepts the a slice with the team that are
// competing, plus an options parameter, with things such as scores and
// previous rankings. The function return a slice of teams that are properly ranked.
func BradleyTerryPart(game []Team, options *Options) []Team {
	epsilon := epsilon(options)
	tbs := betaSq(options) * 2
	_gamma := gamma(options)

	teamRatings := teamRatings(options)(game)
	adjacentTeams := ladderPairs(teamRatings)

	zipper := lo.Zip2(teamRatings, adjacentTeams)

	return lo.Map(zipper, func(item lo.Tuple2[*teamRating, []*teamRating], index int) Team {
		iTeamRating, iAdjacents := item.Unpack()

		var iMu, iSigmaSq, iTeam, iRank = iTeamRating.TeamMu, iTeamRating.TeamSigmaSq, iTeamRating.Team, iTeamRating.Rank

		_sums := lo.Reduce(iAdjacents, func(agg sums, localItem *teamRating, index int) sums {
			var qMu, qSigmaSq, qRank = localItem.TeamMu, localItem.TeamSigmaSq, localItem.Rank

			ciq := math.Sqrt(iSigmaSq + qSigmaSq + tbs)
			piq := 1 / (1 + math.Exp((qMu-iMu)/ciq))
			sigSqToCiq := iSigmaSq / ciq
			iGamma := _gamma(ciq, int64(len(teamRatings)), iTeamRating.TeamMu, iTeamRating.TeamSigmaSq, iTeamRating.Team, iTeamRating.Rank)

			agg.omegaSum += sigSqToCiq * (score(qRank, iRank) - piq)
			agg.deltaSum += ((iGamma * sigSqToCiq) / ciq) * piq * (1 - piq)

			return agg
		}, sums{omegaSum: 0, deltaSum: 0})

		result := lo.Map([]*Rating(*iTeam), func(finalItem *Rating, index int) *Rating {
			sigmaSq := math.Pow(finalItem.SkillUncertaintyDegree, 2)
			mu := finalItem.AveragePlayerSkill + (sigmaSq/iSigmaSq)*_sums.omegaSum
			sigma := finalItem.SkillUncertaintyDegree * math.Sqrt(math.Max(1-(sigmaSq/iSigmaSq)*_sums.deltaSum, epsilon))

			return &Rating{
				AveragePlayerSkill:     mu,
				SkillUncertaintyDegree: sigma,
			}
		})

		return Team(result)
	})
}
