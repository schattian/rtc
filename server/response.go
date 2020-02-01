package server

import "github.com/sebach1/rtc/git"

type respBody struct {
	Commit      *git.Commit
	Commits     []*git.Commit
	Change      *git.Change
	PullRequest *git.PullRequest
}

// type respBodyErr struct {
// 	Err error
// }
