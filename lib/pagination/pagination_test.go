package pagination

import (
	"fmt"
	"reflect"
	"testing"
)

type TestDBConnector struct {
	testID string
}

func (c *TestDBConnector) ReadPaginationIds(q string) ([]DBResult, error) {
	// id|first_name
	// 1|AA
	// 2|AA
	// 3|AV
	// 4|DE
	// 5|ZZ
	// 6|BC
	switch c.testID {
	case "nocursor_get_primary_asc":
		return []DBResult{
			{PID: "1", SortID: "AA"},
			{PID: "2", SortID: "AA"},
			{PID: "3", SortID: "AV"},
			{PID: "4", SortID: "DE"},
		}, nil

	case "nocursor_no_sortinfo_provided":
		return []DBResult{
			{PID: "1", SortID: "AA"},
			{PID: "2", SortID: "AA"},
			{PID: "3", SortID: "AV"},
			{PID: "4", SortID: "DE"},
		}, nil
	case "nocursor_get_primary_desc":
		return []DBResult{
			{PID: "6", SortID: "BC"},
			{PID: "5", SortID: "ZZ"},
			{PID: "4", SortID: "DE"},
			{PID: "3", SortID: "AV"},
		}, nil
	case "nocursor_get_non_primary_asc":
		return []DBResult{
			{PID: "1", SortID: "AA"},
			{PID: "2", SortID: "AA"},
			{PID: "3", SortID: "AV"},
			{PID: "6", SortID: "BE"},
		}, nil
	case "nocursor_get_non_primary_desc":
		return []DBResult{
			{PID: "5", SortID: "ZZ"},
			{PID: "4", SortID: "DE"},
			{PID: "6", SortID: "BC"},
			{PID: "3", SortID: "AV"},
		}, nil
	case "cursor_no_sortinfo_provided":
		return []DBResult{
			{PID: "3", SortID: "AV"},
			{PID: "4", SortID: "DE"},
			{PID: "5", SortID: "ZZ"},
		}, nil
	case "cursor_get_primary_asc":
		return []DBResult{
			{PID: "3", SortID: "AV"},
			{PID: "4", SortID: "DE"},
			{PID: "5", SortID: "ZZ"},
		}, nil
	case "cursor_get_primary_desc":
		return []DBResult{
			{PID: "2", SortID: "AA"},
			{PID: "1", SortID: "AA"},
		}, nil
	case "cursor_get_non_primary_asc":
		return []DBResult{
			{PID: "6", SortID: "BE"},
			{PID: "4", SortID: "DE"},
			{PID: "5", SortID: "ZZ"},
		}, nil
	case "cursor_get_non_primary_desc":
		return []DBResult{
			{PID: "6", SortID: "BE"},
			{PID: "3", SortID: "AV"},
			{PID: "2", SortID: "AA"},
		}, nil

	default:
		return nil, fmt.Errorf("test %s not found", c.testID)
	}
}

// no cursor no sort info provided
// no cursor get primary asc
// no cursor get primary desc
// no cursor get non primary field asc
// no cursor get non primary field desc

//Cursor -> after
// cursor no sort info provided
// cursor get primary asc
// cursor get primary desc
// cursor get non primary field asc
// cursor get non primary field desc

// TODO
//Cursor -> before
// cursor no sort info provided
// cursor get primary asc
// cursor get primary desc
// cursor get non primary field asc
// cursor get non primary field desc

func TestNew(t *testing.T) {
	dbInfo := DatabaseInfo{
		Tables: map[string]TableInfo{
			"users": {PrimaryField: "id", PrimaryFieldToString: func(a any) string { return fmt.Sprintf("%d", a) }},
		},
	}
	type args struct {
		sortInfo SortInfo
		dbReader IDReader
		dbInfo   DatabaseInfo
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   *Cursor
		wantErr bool
	}{
		{
			name: "nocursor no sort info provided",
			args: args{
				dbReader: &TestDBConnector{testID: "nocursor_no_sortinfo_provided"},
				sortInfo: SortInfo{
					Field: "",
					Order: "",
					Limit: 3,
					Table: "users"},
				dbInfo: dbInfo,
			},
			want:    []string{"1", "2", "3"},
			want1:   &Cursor{After: ToPointer("AFTER|3|id")},
			wantErr: false,
		},
		{
			name: "nocursor get primary asc",
			args: args{
				dbReader: &TestDBConnector{testID: "nocursor_get_primary_asc"},
				sortInfo: SortInfo{
					Field: "id",
					Order: "asc",
					Limit: 3,
					Table: "users"},
				dbInfo: dbInfo,
			},
			want:    []string{"1", "2", "3"},
			want1:   &Cursor{After: ToPointer("AFTER|3|id")},
			wantErr: false,
		},
		{
			name: "nocursor get primary desc",
			args: args{
				dbReader: &TestDBConnector{testID: "nocursor_get_primary_desc"},
				sortInfo: SortInfo{
					Field: "id",
					Order: "desc",
					Limit: 3,
					Table: "users"},
				dbInfo: dbInfo,
			},
			want:    []string{"6", "5", "4"},
			want1:   &Cursor{After: ToPointer("AFTER|4|id")},
			wantErr: false,
		},
		{
			name: "nocursor get non primary asc",
			args: args{
				dbReader: &TestDBConnector{testID: "nocursor_get_non_primary_asc"},
				sortInfo: SortInfo{
					Field: "first_name",
					Order: "asc",
					Limit: 3,
					Table: "users"},
				dbInfo: dbInfo,
			},
			want:    []string{"1", "2", "3"},
			want1:   &Cursor{After: ToPointer("AFTER|3;AV|first_name")},
			wantErr: false,
		},
		{
			name: "nocursor get non primary desc",
			args: args{
				dbReader: &TestDBConnector{testID: "nocursor_get_non_primary_desc"},
				sortInfo: SortInfo{
					Field: "first_name",
					Order: "desc",
					Limit: 3,
					Table: "users"},
				dbInfo: dbInfo,
			},
			want:    []string{"5", "4", "6"},
			want1:   &Cursor{After: ToPointer("AFTER|6;BC|first_name")},
			wantErr: false,
		},
		{
			name: "cursor no sort info provided",
			args: args{
				dbReader: &TestDBConnector{testID: "cursor_no_sortinfo_provided"},
				sortInfo: SortInfo{
					Field:      "",
					Order:      "",
					Limit:      2,
					Table:      "users",
					CursorInfo: ToPointer(ParseCursor("AFTER|2|id")),
				},
				dbInfo: dbInfo,
			},
			want:    []string{"3", "4"},
			want1:   &Cursor{After: ToPointer("AFTER|4|id"), Before: ToPointer("BEFORE|3|id")},
			wantErr: false,
		},
		{
			name: "cursor get primary asc",
			args: args{
				dbReader: &TestDBConnector{testID: "cursor_get_primary_asc"},
				sortInfo: SortInfo{
					Field:      "id",
					Order:      "asc",
					Limit:      2,
					Table:      "users",
					CursorInfo: ToPointer(ParseCursor("AFTER|2|id"))},
				dbInfo: dbInfo,
			},
			want:    []string{"3", "4"},
			want1:   &Cursor{After: ToPointer("AFTER|4|id"), Before: ToPointer("BEFORE|3|id")},
			wantErr: false,
		},
		{
			name: "cursor get primary desc",
			args: args{
				dbReader: &TestDBConnector{testID: "cursor_get_primary_desc"},
				sortInfo: SortInfo{
					Field:      "id",
					Order:      "desc",
					Limit:      2,
					Table:      "users",
					CursorInfo: ToPointer(ParseCursor("AFTER|3|id"))},
				dbInfo: dbInfo,
			},
			want:    []string{"2", "1"},
			want1:   &Cursor{Before: ToPointer("BEFORE|2|id")},
			wantErr: false,
		},
		{
			name: "cursor get non primary asc",
			args: args{
				dbReader: &TestDBConnector{testID: "cursor_get_non_primary_asc"},
				sortInfo: SortInfo{
					Field:      "first_name",
					Order:      "asc",
					Limit:      2,
					Table:      "users",
					CursorInfo: ToPointer(ParseCursor("AFTER|3;AV|first_name"))},
				dbInfo: dbInfo,
			},
			want:    []string{"6", "4"},
			want1:   &Cursor{After: ToPointer("AFTER|4;DE|first_name"), Before: ToPointer("BEFORE|6;BE|first_name")},
			wantErr: false,
		},
		{
			name: "cursor get non primary desc",
			args: args{
				dbReader: &TestDBConnector{testID: "cursor_get_non_primary_desc"},
				sortInfo: SortInfo{
					Field:      "first_name",
					Order:      "desc",
					Limit:      2,
					Table:      "users",
					CursorInfo: ToPointer(ParseCursor("AFTER|4;DE|first_name"))},
				dbInfo: dbInfo,
			},
			want:    []string{"6", "3"},
			want1:   &Cursor{After: ToPointer("AFTER|3;AV|first_name"), Before: ToPointer("BEFORE|6;BE|first_name")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := New(tt.args.sortInfo, tt.args.dbReader, tt.args.dbInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				if got1.After != nil {
					if tt.want1.After == nil {
						fmt.Printf("wanted no after cursor. got %v\n", *got1.After)
					} else {
						fmt.Printf("got '%v', want '%v'\n", *got1.After, *tt.want1.After)
					}
				}
				if got1.Before != nil {
					if tt.want1.Before == nil {
						fmt.Printf("wanted no before cursor. got %v\n", *got1.Before)
					} else {
						fmt.Printf("got '%v', want '%v'\n", *got1.Before, *tt.want1.Before)
					}
				}
				t.Errorf("New() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
