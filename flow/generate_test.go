package flow

import (
	"reflect"
	"slices"
	"testing"
)

func TestForward(t *testing.T) {
	type args struct {
		begin int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Simple",
			args: args{begin: 0, end: 3},
			want: []int{0, 1, 2},
		},
		{
			name: "Empty",
			args: args{begin: 2, end: 0},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Forward(tt.args.begin, tt.args.end))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackward(t *testing.T) {
	type args struct {
		begin int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Simple",
			args: args{begin: 0, end: 3},
			want: []int{2, 1, 0},
		},
		{
			name: "Empty",
			args: args{begin: 2, end: 0},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Backward(tt.args.begin, tt.args.end))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Backward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSince(t *testing.T) {
	type args struct {
		start int
	}
	tests := []struct {
		name string
		args args
		take int
		want []int
	}{
		{
			name: "Empty",
			args: args{start: 5},
			take: 0,
			want: []int{},
		},
		{
			name: "Simple",
			args: args{start: 0},
			take: 3,
			want: []int{0, 1, 2},
		},
		{
			name: "FromNegative",
			args: args{start: -2},
			take: 5,
			want: []int{-2, -1, 0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := []int{}
			count := 0
			for val := range Since(tt.args.start) {
				if count == tt.take {
					break
				}
				got = append(got, val)
				count += 1
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Since() = %v, want %v", got, tt.want)
			}
		})
	}
}
