package assist

import (
	"database/sql/driver"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

type ExecStubber struct {
	Expect string
	Err    error
	Result driver.Result
}

type QueryStubber struct {
	Expect string
	Err    error
	Rows   *sqlmock.Rows
}

func (exec *ExecStubber) Stub(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	exec.Expect = regexp.QuoteMeta(exec.Expect)
	expect := mock.ExpectExec(exec.Expect)
	if exec.Err != nil {
		return expect.WillReturnError(exec.Err)
	}
	return expect.WillReturnResult(exec.Result)
}

func (query *QueryStubber) Stub(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	query.Expect = regexp.QuoteMeta(query.Expect)
	expect := mock.ExpectQuery(query.Expect)
	if query.Err != nil {
		return expect.WillReturnError(query.Err)
	}
	return expect.WillReturnRows(query.Rows)
}
