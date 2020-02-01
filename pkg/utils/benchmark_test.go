package utils

import (
	"reflect"
	"testing"
)

func TestReuseString(t *testing.T) {
	a := []string{"1", "2", "3", "3", "4"}
	b := []string{"1", "2", "4"}
	type args struct {
		origin []string
		v      string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"1",
			args{origin: a, v: "3"},
			b,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReuseString(tt.args.origin, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reuse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReuseInt64(t *testing.T) {
	a := []int64{1, 2, 3, 4, 1, 5, 6, 7}
	b := []int64{2, 3, 4, 5, 6, 7}
	type args struct {
		origin []int64
		v      int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			"1",
			args{origin: a, v: 1},
			b,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReuseInt64(tt.args.origin, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reuse() = %v, want %v", got, tt.want)
			}
		})
	}
}
