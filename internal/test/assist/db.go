package assist

import (
	"database/sql/driver"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sebach1/rtc/internal/store"
)

// ExecStubber is the actor which performs stubs of db.Exec()
type ExecStubber struct {
	Expect string
	Err    error
	Result driver.Result
}

// QueryStubber is the actor which performs stubs of db.Query()
type QueryStubber struct {
	Expect string
	Err    error
	Rows   *sqlmock.Rows
}

// Stub stubs the execution with the given mock
// It uses .Err to provide connection errs, and Result to stub the desired output on caller
func (exec *ExecStubber) Stub(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	exec.Expect = regexp.QuoteMeta(exec.Expect)
	expect := mock.ExpectExec(exec.Expect)
	if exec.Err != nil {
		return expect.WillReturnError(exec.Err)
	}
	return expect.WillReturnResult(exec.Result)
}

// Stub stubs the query with the given mock
// It uses .Err to provide connection errs, and Rows to stub the desired output on caller
func (query *QueryStubber) Stub(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	query.Expect = regexp.QuoteMeta(query.Expect)
	expect := mock.ExpectQuery(query.Expect)
	if query.Err != nil {
		return expect.WillReturnError(query.Err)
	}
	return expect.WillReturnRows(query.Rows)
}

// RowsFor returns sqlmock.Rows with the respective row of columns given a store.Storable entity
func RowsFor(entity store.Storable) *sqlmock.Rows {
	return sqlmock.NewRows(entity.SQLColumns())
}
