package main

import "testing"

func Test_prometheusFormat(t *testing.T) {
	type args struct {
		logfile  string
		pattern  string
		nbErrors int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic test", args{"/var/example", "patrn", 10}, "pattern_in_log_count{logfile=\"/var/example\", pattern=\"patrn\"} 10\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prometheusFormat(tt.args.logfile, tt.args.pattern, tt.args.nbErrors); got != tt.want {
				t.Errorf("prometheusFormat() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
