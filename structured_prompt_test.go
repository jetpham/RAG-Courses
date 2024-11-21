package main

import "testing"

func Test_isSimilar(t *testing.T) {
	setup, err := newSetup()
	if err != nil {
		t.Errorf("Error during setup: %v", err)
		return
	}
	type args struct {
		text1 string
		text2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Similar texts",
			args: args{
				text1: "The quick brown fox jumps over the lazy dog.",
				text2: "A fast, dark-colored fox leaps over a sleepy dog.",
			},
			want: true,
		},
		{
			name: "Different texts",
			args: args{
				text1: "The quick brown fox jumps over the lazy dog.",
				text2: "The weather is sunny today.",
			},
			want: false,
		},
		{
			name: "Identical texts",
			args: args{
				text1: "The quick brown fox jumps over the lazy dog.",
				text2: "The quick brown fox jumps over the lazy dog.",
			},
			want: true,
		},
		{
			name: "Similar meaning, different words",
			args: args{
				text1: "He is very happy.",
				text2: "He feels great joy.",
			},
			want: true,
		},
		{
			name: "Different topics",
			args: args{
				text1: "The stock market is volatile.",
				text2: "She loves to paint landscapes.",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similar, reason := isSimilar(setup, tt.args.text1, tt.args.text2)
			if similar != tt.want {
				t.Errorf("text1:\n\n%v\n\ntext2:\n\n%v\n\nreason:\n\n%v", tt.args.text1, tt.args.text2, reason)
			}
		})
	}
}
