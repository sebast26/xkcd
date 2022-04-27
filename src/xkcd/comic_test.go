package xkcd

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_matchesTerm(t *testing.T) {
	type args struct {
		comic *Comic
		term  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "matching term",
			args: args{
				comic: &Comic{
					Transcript: "it is a match",
				},
				term: "match",
			},
			want: true,
		},
		{
			name: "not matching term",
			args: args{
				comic: &Comic{
					Transcript: "something else",
				},
				term: "not match",
			},
			want: false,
		},
		{
			name: "empty transcript",
			args: args{
				comic: &Comic{
					Transcript: "",
				},
				term: "not match",
			},
			want: false,
		},
		{
			name: "empty term",
			args: args{
				comic: &Comic{
					Transcript: "empty term",
				},
				term: "",
			},
			want: true,
		},
		{
			name: "whole term",
			args: args{
				comic: &Comic{
					Transcript: "should take into account whole term",
				},
				term: "accountwhole",
			},
			want: false,
		},
		{
			name: "match characters",
			args: args{
				comic: &Comic{
					Transcript: "seafood",
				},
				term: "afo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesTerm(tt.args.comic, tt.args.term); got != tt.want {
				t.Errorf("matchesTerm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		handle io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Comic
		wantErr bool
	}{
		{
			name: "comic read",
			args: args{
				handle: strings.NewReader(`{"month": "4", "num": 571, "link": "", "year": "2009", "news": "", "safe_title": "Can't Sleep", "transcript": "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}", "alt": "If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.", "img": "https://imgs.xkcd.com/comics/cant_sleep.png", "title": "Can't Sleep", "day": "20"}`), // nolint
			},
			want: &Comic{
				URL:        "https://imgs.xkcd.com/comics/cant_sleep.png",
				Transcript: "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}", // nolint
			},
			wantErr: false,
		},
		{
			name: "minimum comic",
			args: args{
				handle: strings.NewReader(`{"transcript": "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}", "img": "https://imgs.xkcd.com/comics/cant_sleep.png"}`), // nolint
			},
			want: &Comic{
				URL:        "https://imgs.xkcd.com/comics/cant_sleep.png",
				Transcript: "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}", // nolint
			},
			wantErr: false,
		},
		{
			name: "empty handle",
			args: args{
				handle: strings.NewReader(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid JSON",
			args: args{
				handle: strings.NewReader(`{"transci:"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
