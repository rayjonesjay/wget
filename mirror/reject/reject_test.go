package reject

import (
	"testing"
	"wget/ctx"
)

func TestReject(t *testing.T) {
	type args struct {
		ctx        ctx.Context
		mirrorPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				ctx:        ctx.Context{Rejects: []string{}},
				mirrorPath: "",
			},
			want: false,
		},

		{
			name: "Mirror root URL",
			args: args{
				ctx:        ctx.Context{Rejects: []string{"/"}},
				mirrorPath: "/",
			},
			want: true,
		},

		{
			name: "Root /img",
			args: args{
				ctx:        ctx.Context{Rejects: []string{"/img"}},
				mirrorPath: "/img",
			},
			want: true,
		},

		{
			name: "Subfolder /img",
			args: args{
				ctx:        ctx.Context{Rejects: []string{"/img"}},
				mirrorPath: "/img",
			},
			want: true,
		},

		{
			name: "Nested Subfolder /img",
			args: args{
				ctx:        ctx.Context{Rejects: []string{"/img"}},
				mirrorPath: "/home/path/img/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Reject(tt.args.ctx, tt.args.mirrorPath); got != tt.want {
					t.Errorf("Reject() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
