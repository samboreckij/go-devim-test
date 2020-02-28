package env

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestValue_String(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		get   string
		def   string
		want  string
	}{
		{
			name: "Exists",
			setup: func() {
				os.Setenv("TEST1_USERNAME", "username")
			},
			get:  "TEST1_USERNAME",
			def:  "defuser",
			want: "username",
		},
		{
			name: "Not exists",
			get:  "TEST2_USERNAME",
			def:  "defuser",
			want: "defuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if got := Get(tt.get).String(tt.def); got != tt.want {
				t.Errorf("Get().String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Int(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		get   string
		def   int
		want  int
	}{
		{
			name: "Exists",
			setup: func() {
				os.Setenv("TEST1_BULK", "3000")
			},
			get:  "TEST1_BULK",
			def:  11,
			want: 3000,
		},
		{
			name: "Not exists",
			get:  "TEST2_BULK",
			def:  11,
			want: 11,
		},
		{
			name: "Zero",
			get:  "TEST3_BULK",
			setup: func() {
				os.Setenv("TEST3_BULK", "0")
			},
			def:  12,
			want: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if got := Get(tt.get).Int(tt.def); got != tt.want {
				t.Errorf("Get().Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Bool(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		get   string
		def   bool
		want  bool
	}{
		{
			name: "Exists",
			setup: func() {
				os.Setenv("TEST1_ENABLED", "false")
			},
			get:  "TEST1_ENABLED",
			def:  true,
			want: false,
		},
		{
			name: "Not exists",
			get:  "TEST2_ENABLED",
			def:  true,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if got := Get(tt.get).Bool(tt.def); got != tt.want {
				t.Errorf("Get().Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Duration(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		get   string
		def   time.Duration
		want  time.Duration
	}{
		{
			name: "Exists",
			setup: func() {
				os.Setenv("TEST1_TIMEOUT", "30s")
			},
			get:  "TEST1_TIMEOUT",
			def:  time.Second,
			want: 30 * time.Second,
		},
		{
			name: "Not exists",
			get:  "TEST2_TIMEOUT",
			def:  time.Second,
			want: time.Second,
		},
		{
			name: "Zero",
			setup: func() {
				os.Setenv("TEST3_TIMEOUT", "0s")
			},
			get:  "TEST3_TIMEOUT",
			def:  time.Second,
			want: time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if got := Get(tt.get).Duration(tt.def); got != tt.want {
				t.Errorf("Get().Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatch(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		timeout time.Duration
	}

	tests := []struct {
		name   string
		args   args
		wantfn func(<-chan *Value) error
	}{
		{
			name: "Changes",
			args: args{
				ctx:     context.Background(),
				key:     "TEST1_WATCHED",
				timeout: time.Nanosecond,
			},
			wantfn: func(ch <-chan *Value) error {
				w1 := "first_val"

				os.Setenv("TEST1_WATCHED", w1)

				if g1, ok := <-ch; !ok || g1.String("") != w1 {
					return fmt.Errorf("%v, want %v", g1, w1)
				}

				w2 := "second_val"

				os.Setenv("TEST1_WATCHED", w2)

				if g2, ok := <-ch; !ok || g2.String("") != w2 {
					return fmt.Errorf("%v, want %v", g2, w2)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(tt.args.ctx, time.Second)

			if err := tt.wantfn(Watch(ctx, tt.args.key, tt.args.timeout)); err != nil {
				t.Errorf("Watch() = %v", err)
			}

			cancel()
		})
	}
}
