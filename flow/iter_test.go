package flow

import (
	"iter"
	"reflect"
	"slices"
	"testing"
)

func TestPack(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SingleValue",
			args: args{v: 3},
			want: []int{3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Pack(tt.args.v))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTake(t *testing.T) {
	type args struct {
		sequence iter.Seq[int]
		count    int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "TakeSufficient",
			args: args{sequence: slices.Values([]int{1, 2, 3, 4, 5}), count: 3},
			want: []int{1, 2, 3},
		},
		{
			name: "TakeInsufficient",
			args: args{sequence: slices.Values([]int{1, 2, 3}), count: 5},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Take(tt.args.sequence, tt.args.count))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDrop(t *testing.T) {
	type args struct {
		sequence iter.Seq[int]
		count    int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "DropSufficient",
			args: args{sequence: slices.Values([]int{1, 2, 3, 4, 5}), count: 3},
			want: []int{4, 5},
		},
		{
			name: "DropInsufficient",
			args: args{sequence: slices.Values([]int{1, 2, 3}), count: 5},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Drop(tt.args.sequence, tt.args.count))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}
