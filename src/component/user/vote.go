package user

/*
  Proposed voting rules - for now these are relaxed
  1 points - submit, comment
  1 points - upvote (they start with 10 points)
  10 points - downvote
  100 points - flag

	karma is collected for comment upvotes *only* not for topic upvotes
	karma is sacrificed in negative actions - flagging and downvoting
*/

// CanUpvote returns true if this user can upvote
func (m *User) CanUpvote() bool {
	return m.Score > -5.0 && !m.Anon() && !m.Banned()
}

// CanDownvote returns true if this user can downvote
func (m *User) CanDownvote() bool {
	return m.Score >= 0 && !m.Anon() && !m.Banned()
}

// CanFlag returns true if this user can flag
func (m *User) CanFlag() bool {
	return m.Admin()
}

// CanSubmit returns true if this user can submit
func (m *User) CanSubmit() bool {
	return m.Score > -5.0 && !m.Anon() && !m.Banned() && !m.IsUnverified()
}

func (m *User) CanPublish() bool {
	return m.Score >= 3.0 && !m.Anon() && !m.Banned() && !m.IsUnverified()
}

// CanComment returns true if this user can comment
func (m *User) CanComment() bool {
	return m.Score > -5.0 && !m.Anon() && !m.Banned() && !m.IsUnverified()
}
