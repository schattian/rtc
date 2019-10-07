package git

// The Agent abstracts every interaction with a crud app
// type Agent interface {

// 	// PushSingle fast-pushes a single change
// 	// Its the representation of a bash alias like:
// 	// `git add . && git commit && git push`
// 	PushSingle(context.Context, *git.Change) error

// 	// Push will push the given commit
// 	Push(context.Context, *git.Commit) (*git.Summary, error)

// 	// Pull retrieves the column(s) given from an entity on a table
// 	// In case of giving empty columns (""), it'll pull in batch

// 	// * Postposed
// 	// // PullCommits is the real representation of Pull, but with much less usability
// 	// // It wraps the Pull functions and converts the response into commits
// 	// // Note that, aswell as the commits need to be
// 	// PullCommits(context.Context, ID, TableName, ...ColumnName) (*git.Commit, error)

// 	// Remove deletes a record given its id and table
// 	Remove(context.Context, ID, TableName) error

// 	// Init initializes every persistant needed action (such as token refreshing)
// 	Init(context.Context) error
// }
