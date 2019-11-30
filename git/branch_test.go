package git

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/internal/test/assist"
	"github.com/sebach1/rtc/internal/test/thelper"
)

var (
	errFoo = errors.New("foo")
)

func TestNewBranchWithIndex(t *testing.T) {
	type args struct {
		ctx  context.Context
		name integrity.BranchName
	}
	tests := []struct {
		name    string
		args    args
		stubs   []*assist.ExecStubber
		want    *Branch
		wantErr error
	}{
		{
			name:    "INDEX creation returns ERR on db CONNECTion",
			args:    args{name: "foo"},
			wantErr: errFoo,
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Err: errFoo},
			},
		},
		{
			name:    "INDEX creation returns ERR on TX",
			args:    args{name: "foo"},
			wantErr: errFoo,
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Result: sqlmock.NewErrorResult(errFoo)},
			},
		},
		{
			name:    "BRANCH creation returns ERR on TX",
			args:    args{name: "foo"},
			wantErr: errFoo,
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Result: sqlmock.NewResult(1, 1)},
				{Expect: "INSERT INTO branches", Result: sqlmock.NewErrorResult(errFoo)},
			},
		},
		{
			name:    "BRANCH creation returns ERR on CONNECTion",
			args:    args{name: "foo"},
			wantErr: errFoo,
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Result: sqlmock.NewResult(1, 1)},
				{Expect: "INSERT INTO branches", Err: errFoo},
			},
		},
		{
			name:    "BRANCH creation performs SUCCESSfully",
			args:    args{name: "foo"},
			wantErr: nil,
			want:    &Branch{Id: 10, IndexId: 1, Name: "foo"},
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Result: sqlmock.NewResult(1, 1)},
				{Expect: "INSERT INTO branches", Result: sqlmock.NewResult(10, 1)},
			},
		},
		{
			name:    "CTX CANCEL",
			args:    args{name: "foo"},
			wantErr: nil,
			want:    &Branch{Id: 10, IndexId: 1, Name: "foo"},
			stubs: []*assist.ExecStubber{
				{Expect: "INSERT INTO indices DEFAULT VALUES", Result: sqlmock.NewResult(1, 1)},
				{Expect: "INSERT INTO branches", Result: sqlmock.NewResult(10, 1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := thelper.MockDB(t)
			for _, stub := range tt.stubs {
				stub.Stub(mock)
			}
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}
			got, err := NewBranchWithIndex(tt.args.ctx, db, tt.args.name)
			if err != tt.wantErr {
				t.Errorf("NewBranchWithIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBranchWithIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranch_FetchIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		branch    *Branch
		newBranch *Branch
		stubs     []*assist.QueryStubber
		args      args
		wantErr   error
	}{
		{
			name:    "INDEX query returns ERR on db CONNECTion",
			wantErr: errFoo,
			branch:  gBranches.Foo.Copy(),
			stubs: []*assist.QueryStubber{
				{Expect: "SELECT * FROM indices", Err: errFoo},
			},
		},
		{
			name:      "fetches index through index id SUCCESSfully",
			branch:    gBranches.Foo.Copy().rmIndexAndReturn(),
			newBranch: gBranches.Foo.Copy().rmIndexChangesAndReturn(),
			stubs: []*assist.QueryStubber{
				{Expect: "SELECT * FROM indices WHERE id=?", Rows: thelper.RowsFor(&Index{}).AddRow(gIndices.Foo.Id)},
			},
		},
		{
			name:    "the given branch has NIL INDEX ID",
			branch:  gBranches.Foo.Copy().rmIndexIdAndReturn(),
			wantErr: errNilIndexId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := thelper.MockDB(t)
			for _, stub := range tt.stubs {
				stub.Stub(mock)
			}
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}
			originalBranch := tt.branch.Copy()
			err := tt.branch.FetchIndex(tt.args.ctx, db)
			if err != tt.wantErr {
				t.Errorf("Branch.FetchIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			thelper.CmpIfErr(t, err, originalBranch, tt.branch, tt.newBranch, "Branch.FetchIndex()")
		})
	}
}

func (b *Branch) rmIndexIdAndReturn() *Branch {
	b.IndexId = 0
	return b
}

func (b *Branch) rmIndexAndReturn() *Branch {
	b.Index = nil
	return b
}

func (b *Branch) rmIndexChangesAndReturn() *Branch {
	b.Index.Changes = nil
	return b
}
