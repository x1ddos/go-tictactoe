// This package performs initialization during startup.
// 
// It is separated from tictactoe package so that the latter can be used
// in another app, e.g. http://github.com/crhym3/go-endpoints.appspot.com.

package app

import (
	"github.com/crhym3/go-endpoints/endpoints"
	"github.com/crhym3/go-tictactoe/tictactoe"
)

func init() {
	if _, err := tictactoe.RegisterService(); err != nil {
		panic(err.Error())
	}
	endpoints.HandleHttp()
}
