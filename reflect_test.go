package z

import "testing"

func TestCopyNoneZeroField(t *testing.T) {
	type some struct {
		ID           int64
		Name         string
		privateField string
	}
	type args struct {
		from interface{}
		to   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test all zero", args: args{
			from: &some{
				ID:           0,
				Name:         "",
				privateField: "",
			},
			to: &some{
				ID:           0,
				Name:         "",
				privateField: "",
			},
		}},
		{name: "test from zero", args: args{
			from: &some{
				ID:           0,
				Name:         "",
				privateField: "",
			},
			to: &some{
				ID:           1,
				Name:         "Bob",
				privateField: "pField",
			},
		}},
		{name: "test to zero", args: args{
			from: &some{
				ID:           2,
				Name:         "Alice",
				privateField: "p-Field",
			},
			to: &some{
				ID:           0,
				Name:         "",
				privateField: "",
			},
		}},
		{name: "test none zero", args: args{
			from: &some{
				ID:           3,
				Name:         "Jack",
				privateField: "haha",
			},
			to: &some{
				ID:           4,
				Name:         "Rose",
				privateField: "hey",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Before:\tfrom=%#v\tto=%#v", tt.args.from, tt.args.to)
			copyNoneZeroField(tt.args.from, tt.args.to)
			t.Logf("After:\tfrom=%#v\tto=%#v", tt.args.from, tt.args.to)
		})
	}
}
