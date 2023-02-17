package openskill

// NewTeam is a small utility function to create a Team from many players.
func NewTeam(teams ...*Rating) Team {
	slc := make([]*Rating, 0)

	slc = append(slc, teams...)

	return Team(slc)
}
