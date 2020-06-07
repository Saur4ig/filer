package main

import (
	"reflect"
	"testing"
)

func Test_getFirst(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Example 1",
			str:  "abcde",
			want: "a",
		},
		{
			name: "Example 2",
			str:  "",
			want: "",
		},
		{
			name: "Example 3",
			str:  "ягвь",
			want: "я",
		},
		{
			name: "Example 4",
			str:  "世界",
			want: "世",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFirst(tt.str); got != tt.want {
				t.Errorf("getFirstLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDatabase(t *testing.T) {
	tests := []struct {
		name string
		want Database
	}{
		{
			name: "Example",
			want: &clickhouse{
				storage: &storage{
					bucket: map[string]*cell{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clickhouse_add(t *testing.T) {
	tests := []struct {
		name    string
		storage *storage
		str     string
	}{
		{
			name:    "Example 1",
			str:     "abc",
			storage: &storage{bucket: map[string]*cell{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clickhouse{
				storage: tt.storage,
			}
			c.add(tt.str)
			if c.isUnique(tt.str) {
				t.Errorf("not added %s", tt.str)
			}
		})
	}
}

func Test_clickhouse_isUnique(t *testing.T) {
	tests := []struct {
		name    string
		storage *storage
		str     string
		want    bool
	}{
		{
			name: "Example 1",
			storage: &storage{
				bucket: map[string]*cell{
					"a": {data: map[string]struct{}{"abcde": struct{}{}}},
				},
			},
			str:  "abcde",
			want: false,
		},
		{
			name: "Example 2",
			storage: &storage{
				bucket: map[string]*cell{
					"s": {data: map[string]struct{}{"saur": struct{}{}}},
				},
			},
			str:  "sau",
			want: true,
		},
		{
			name: "Example 3",
			storage: &storage{
				bucket: map[string]*cell{
					"s": {data: map[string]struct{}{"saur": struct{}{}}},
				},
			},

			str:  "abcdef",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clickhouse{
				storage: tt.storage,
			}
			if got := c.isUnique(tt.str); got != tt.want {
				t.Errorf("isUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clickhouse_all(t *testing.T) {
	tests := []struct {
		name    string
		storage *storage
		want    []string
	}{
		{
			name: "Example 1",
			storage: &storage{
				bucket: map[string]*cell{
					"a": {data: map[string]struct{}{"abc": struct{}{}}},
					"d": {data: map[string]struct{}{"da": struct{}{}}},
					"t": {data: map[string]struct{}{"test": struct{}{}}},
				},
			},
			want: []string{"abc", "da", "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clickhouse{
				storage: tt.storage,
			}
			if got := c.all(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("all() = %v, want %v", got, tt.want)
			}
		})
	}
}
