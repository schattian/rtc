package maptypes

import (
	"testing"
)

func TestChange_IsUntracked(t *testing.T) {
	type fields struct {
		ID         int
		TableName  TableName
		ColumnName ColumnName
		Value      Value
		EntityID   ID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "entity_id is not set up", fields: fields{}, want: true},
		{name: "entity_id is nil", fields: fields{EntityID: nil}, want: true},

		{name: "entity_id is zero-value", fields: fields{EntityID: 0}, want: false},
		{name: "entity_id is filled", fields: fields{EntityID: 20}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chg := &Change{
				ID:         tt.fields.ID,
				TableName:  tt.fields.TableName,
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				EntityID:   tt.fields.EntityID,
			}
			if got := chg.IsUntracked(); got != tt.want {
				t.Errorf("Change.IsUntracked() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIsCompatibleWith(t *testing.T) {
	type args struct {
		chg      *Change
		otherChg *Change
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "different entities but same tableName",
			args: args{chg: &Change{EntityID: 10, TableName: "test"},
				otherChg: &Change{EntityID: 20, TableName: "test"}},
			want: false,
		},
		{
			name: "all different",
			args: args{chg: &Change{EntityID: 10, TableName: "test"},
				otherChg: &Change{EntityID: 20, TableName: "another"}},
			want: false,
		},

		{
			name: "diff tableName and same entity",
			args: args{chg: &Change{EntityID: 10, TableName: "test"},
				otherChg: &Change{EntityID: 10, TableName: "another"}},
			want: false,
		},
		{
			name: "both nil entities and sabe tableName",
			args: args{chg: &Change{TableName: "test"},
				otherChg: &Change{TableName: "test"}},
			want: false,
		},

		{
			name: "is mirrored",
			args: args{chg: &Change{EntityID: 10, TableName: "test"},
				otherChg: &Change{EntityID: 10, TableName: "test"}},
			want: true,
		},
		{
			name: "is mirrored but with diff colName",
			args: args{chg: &Change{EntityID: 10, TableName: "test", ColumnName: "test"},
				otherChg: &Change{EntityID: 10, TableName: "test", ColumnName: "another"}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCompatibleWith(tt.args.chg, tt.args.otherChg); got != tt.want {
				t.Errorf("IsCompatibleWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
