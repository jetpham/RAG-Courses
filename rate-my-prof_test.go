package main

import (
	"testing"
)

func Test_getRateMyProfessorData(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    ProfessorInfo
		wantErr bool
	}{
		{
			name: "philip peterson",
			args: args{name: "philip peterson"},
			want: ProfessorInfo{
				Name:           "Phil Peterson", // name doesn't match exactly
				Department:     "Computer Science",
				School:         "University of San Francisco",
				Rating:         3.8,
				Difficulty:     4.4,
				TotalRatings:   17,
				WouldTakeAgain: 70.6,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRateMyProfessorData(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRateMyProfessorData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Department != tt.want.Department || got.School != tt.want.School || got.Rating != tt.want.Rating || got.Difficulty != tt.want.Difficulty || got.TotalRatings != tt.want.TotalRatings || got.WouldTakeAgain != tt.want.WouldTakeAgain {
				t.Errorf("getRateMyProfessorData() = %v, want %v", got, tt.want)
			}
		})
	}
}
