package tictactoe

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"appengine/user"

	"github.com/crhym3/go-endpoints/endpoints"
)

const clientId = "YOUR-CLIENT-ID"

var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{clientId, endpoints.ApiExplorerClientId}
	// in case we'll want to use TicTacToe API from an Android app
	audiences = []string{clientId}
)

type BoardMsg struct {
	State string `json:"state" endpoints:"required"`
}

type ScoreReqMsg struct {
	Outcome string `json:"outcome" endpoints:"required"`
}

type ScoreRespMsg struct {
	Id      int64  `json:"id"`
	Outcome string `json:"outcome"`
	Played  string `json:"played"`
}

type ScoresListReq struct {
	Limit int `json:"limit"`
}

type ScoresListResp struct {
	Items []*ScoreRespMsg `json:"items"`
}

// TicTacToe API service
type TicTacToeApi struct {
}

// BoardGetMove simulates a computer move in tictactoe.
// Exposed as API endpoint
func (ttt *TicTacToeApi) BoardGetMove(r *http.Request,
	req *BoardMsg, resp *BoardMsg) error {

	const boardLen = 9
	if len(req.State) != boardLen {
		return fmt.Errorf("Bad Request: Invalid board: %q", req.State)
	}
	runes := []rune(req.State)
	freeIndices := make([]int, 0)
	for pos, r := range runes {
		if r != 'O' && r != 'X' && r != '-' {
			return fmt.Errorf("Bad Request: Invalid rune: %q", r)
		}
		if r == '-' {
			freeIndices = append(freeIndices, pos)
		}
	}
	freeIdxLen := len(freeIndices)
	if freeIdxLen > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIdx := r.Intn(freeIdxLen)
		runes[freeIndices[randomIdx]] = 'O'
		resp.State = string(runes)
	} else {
		return fmt.Errorf("Bad Request: This board is full: %q", req.State)
	}
	return nil
}

// ScoresList queries scores for the current user.
// Exposed as API endpoint
func (ttt *TicTacToeApi) ScoresList(r *http.Request,
	req *ScoresListReq, resp *ScoresListResp) error {

	c := endpoints.NewContext(r)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}
	q := newUserScoreQuery(u)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	scores, err := fetchScores(c, q, req.Limit)
	if err != nil {
		return err
	}
	resp.Items = make([]*ScoreRespMsg, len(scores))
	for i, score := range scores {
		resp.Items[i] = score.toMessage(nil)
	}
	return nil
}

// ScoresInsert inserts a new score for the current user.
func (ttt *TicTacToeApi) ScoresInsert(r *http.Request,
	req *ScoreReqMsg, resp *ScoreRespMsg) error {

	c := endpoints.NewContext(r)
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}
	score := newScore(req.Outcome, u)
	if err := score.put(c); err != nil {
		return err
	}
	score.toMessage(resp)
	return nil
}

// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c endpoints.Context) (*user.User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	c.Debugf("Current user: %#v", u)
	return u, nil
}

// RegisterService exposes TicTacToeApi methods as API endpoints.
// 
// The registration/initialization during startup is not performed here but
// in app package. It is separated from this package (tictactoe) so that the
// service and its methods defined here can be used in another app,
// e.g. http://github.com/crhym3/go-endpoints.appspot.com.
func RegisterService() (*endpoints.RpcService, error) {
	api := &TicTacToeApi{}
	rpcService, err := endpoints.RegisterService(api,
		"tictactoe", "v1", "Tic Tac Toe API", true)
	if err != nil {
		return nil, err
	}

	info := rpcService.MethodByName("BoardGetMove").Info()
	info.Path, info.HttpMethod, info.Name = "board", "POST", "board.getmove"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ScoresList").Info()
	info.Path, info.HttpMethod, info.Name = "scores", "GET", "scores.list"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	info = rpcService.MethodByName("ScoresInsert").Info()
	info.Path, info.HttpMethod, info.Name = "scores", "POST", "scores.insert"
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIds, audiences

	return rpcService, nil
}
