package openskill

import (
	"math"

	"github.com/samber/lo"
)

// ThurstoneMostellerPart is a implementation of the Thurstone-Mosteller ranking model
// that uses partial pairing. The Thurstone-Mosteller model uses gaussian distribution
// to properly rank the teams. Partial pairing is less accurate than a
// full pairing, but works better in situations with a high number of teams.
// This function accepts the a slice with the team that are
// competing, plus an options parameter, with things such as scores and
// previous rankings. The function return a slice of teams that are properly ranked.
func ThurstoneMostellerPart(game []Team, options *Options) []Team {
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

			ciq := 2 * math.Sqrt(iSigmaSq+qSigmaSq+tbs)
			deltaMu := (iMu - qMu) / ciq
			sigSqToCiq := iSigmaSq / ciq
			iGamma := _gamma(ciq, int64(len(teamRatings)), iTeamRating.TeamMu, iTeamRating.TeamSigmaSq, iTeamRating.Team, iTeamRating.Rank)

			if qRank == iRank {
				agg.omegaSum += sigSqToCiq * vt(deltaMu, epsilon/ciq)
				agg.deltaSum += ((iGamma * sigSqToCiq) / ciq) * wt(deltaMu, epsilon/ciq)

				return agg
			}

			sign := lo.Ternary(qRank > iRank, 1.0, -1.0)

			agg.omegaSum += sign * sigSqToCiq * v(sign*deltaMu, epsilon/ciq)
			agg.deltaSum += ((iGamma * sigSqToCiq) / ciq) * w(sign*deltaMu, epsilon/ciq)

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
