package store

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/internal/test/assist"
	"github.com/sebach1/rtc/internal/test/thelper"
)

var errFoo = errors.New("foo")

func TestInsertIntoDB(t *testing.T) {
	type args struct {
		ctx       context.Context
		storables []Storable
	}
	tests := []struct {
		name          string
		args          args
		wantErr       error
		wantStorables []Storable
		stub          *assist.ExecStubber
	}{
		{
			name: "successfully single insert",
			args: args{
				storables: []Storable{&storableStub{Name: "foo"}},
			},
			wantStorables: []Storable{&storableStub{Name: "foo", Id: 10}},
			stub:          &assist.ExecStubber{Expect: "INSERT INTO entities_stub", Result: sqlmock.NewResult(10, 0)},
		},
		{
			name: "exec returns err",
			args: args{
				storables: []Storable{&storableStub{}},
			},
			wantStorables: []Storable{&storableStub{}},
			wantErr:       errFoo,
			stub:          &assist.ExecStubber{Expect: "INSERT INTO entities_stub", Result: sqlmock.NewErrorResult(errFoo)},
		},
		{
			name:    "nil interface given",
			args:    args{},
			wantErr: errNilStorableEntity,
			stub:    &assist.ExecStubber{Expect: "INSERT INTO entities_stub", Result: sqlmock.NewErrorResult(errFoo)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStorables := tt.args.storables
			db, mock := thelper.MockDB(t)
			tt.stub.Stub(mock)
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}
			err := InsertIntoDB(tt.args.ctx, db, tt.args.storables...)
			if err != tt.wantErr {
				t.Errorf("InsertIntoDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if diff := cmp.Diff(oldStorables, tt.args.storables); diff != "" {
					t.Errorf("InsertIntoDB() errored mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.wantStorables, tt.args.storables); diff != "" {
				t.Errorf("InsertIntoDB() mismatch (-want +got): %s", diff)
			}
		})
	}
}
func TestUpdateIntoDB(t *testing.T) {
	type args struct {
		ctx       context.Context
		storables []Storable
	}
	tests := []struct {
		name          string
		args          args
		wantErr       error
		wantStorables []Storable
		stub          *assist.ExecStubber
	}{
		{
			name:    "update returns err",
			wantErr: errFoo,
			args: args{
				storables: []Storable{&storableStub{}},
			},
			stub: &assist.ExecStubber{Expect: "UPDATE entities_stub", Result: sqlmock.NewErrorResult(errFoo)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// oldStorable := tt.args.storable
			db, mock := thelper.MockDB(t)
			tt.stub.Stub(mock)
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}
			err := UpdateIntoDB(tt.args.ctx, db, tt.args.storables...)
			if err != tt.wantErr {
				t.Errorf("UpdateIntoDB() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

type storableStub struct {
	Id   int64
	Name string
}

func (s *storableStub) GetId() int64 {
	return s.Id
}

func (s *storableStub) SetId(id int64) {
	s.Id = id
}

func (s *storableStub) SQLTable() string {
	return "entities_stub"
}

func (s *storableStub) SQLColumns() []string {
	return []string{
		"id",
		"name",
	}
}
