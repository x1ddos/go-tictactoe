package tictactoe

import (
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"golang.org/x/net/context"
)

const (
	TIME_LAYOUT = "Jan 2, 2006 15:04:05 AM"
	SCORE_KIND  = "Score"
)

// Score is an entity to store scores that have been inserted by users.
type Score struct {
	key *datastore.Key

	Outcome string    `datastore:"outcome"`
	Played  time.Time `datastore:"played"`
	Player  string    `datastore:"player"`
}

// Turns the Score struct/entity into a ScoreRespMsg which is then used
// as an API response.
func (s *Score) toMessage(msg *ScoreRespMsg) *ScoreRespMsg {
	if msg == nil {
		msg = &ScoreRespMsg{}
	}
	msg.Id = s.key.IntID()
	msg.Outcome = s.Outcome
	msg.Played = s.timestamp()
	return msg
}

// timestamp formats date/time of the score.
func (s *Score) timestamp() string {
	return s.Played.Format(TIME_LAYOUT)
}

// put stores the score in the Datastore.
func (s *Score) put(c context.Context) (err error) {
	key := s.key
	if key == nil {
		key = datastore.NewIncompleteKey(c, SCORE_KIND, nil)
	}
	key, err = datastore.Put(c, key, s)
	if err == nil {
		s.key = key
	}
	return
}

// newScore returns a new Score ready to be stored in the Datastore.
func newScore(outcome string, u *user.User) *Score {
	return &Score{Outcome: outcome, Played: time.Now(), Player: userId(u)}
}

// newUserScoreQuery returns a Query which can be used to list all previous
// games of a user.
func newUserScoreQuery(u *user.User) *datastore.Query {
	return datastore.NewQuery(SCORE_KIND).Filter("player =", userId(u))
}

// fetchScores runs Query q and returns Score entities fetched from the
// Datastore.
func fetchScores(c context.Context, q *datastore.Query, limit int) (
	[]*Score, error) {

	scores := make([]*Score, 0, limit)
	keys, err := q.Limit(limit).GetAll(c, &scores)
	if err != nil {
		return nil, err
	}
	for i, score := range scores {
		score.key = keys[i]
	}
	return scores, nil
}

// userId returns a string ID of the user u to be used as Player of Score.
func userId(u *user.User) string {
	return u.String()
}
