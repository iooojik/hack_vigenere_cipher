package main

import (
	"reflect"
	"testing"
)

func Test_countEntryIndex(t *testing.T) {
	type args struct {
		entries map[int][]int
		length  int
	}
	tests := []struct {
		name string
		args args
		want map[int]map[string]float64
	}{
		{
			name: "test#1",
			args: args{
				entries: map[int][]int{
					0: {50, 17, 4, 6},
					1: {56, 16, 24, 36},
				},
				length: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPositions(tt.args.entries, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}
