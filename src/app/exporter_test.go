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

func Test_countOccurences(t *testing.T) {
	type args struct {
		logfile string
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"count errors in big file", args{logfile: "../test_resources/shopify_webhook.log", pattern: "ERROR"}, 33, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countOccurences(tt.args.logfile, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("countOccurences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countOccurences() got = %v, want %v", got, tt.want)
			}
		})
	}
}
