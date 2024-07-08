package gitconfig

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

var (
	gitConfig = New()
)

func TestGitConfig_Set(t *testing.T) {
	type args struct {
		key  string
		vals []interface{}
	}
	tests := []struct {
		name    string
		g       *GitConfig
		args    args
		wantErr bool
	}{
		{
			name: "Empty Value",
			g:    gitConfig,
			args: args{
				key: "foo.foo",
			},
			wantErr: true,
		},
		{
			name: "Valid Key Single Value",
			g:    gitConfig,
			args: args{key: "foo.foo", vals: []any{"foo"}},
		},
		{
			name: "Valid Key Multiple Values",
			g:    gitConfig,
			args: args{key: "foo.bar", vals: []any{"foo", "bar"}},
		},
		{
			name:    "Invalid Key Single Value",
			g:       gitConfig,
			args:    args{key: "foo", vals: []any{"bar"}},
			wantErr: true,
		},
		{
			name:    "Invalid Key Multiple Values",
			g:       gitConfig,
			args:    args{key: "foo", vals: []any{"bar", "baz"}},
			wantErr: true,
		},
		{
			name: "Key With Subsection Single Value",
			g:    gitConfig,
			args: args{
				key:  "foo.bar.baz",
				vals: []any{"foo"},
			},
		},
		{
			name: "Key With Subsection Multiple Values",
			g:    gitConfig,
			args: args{
				key:  "baz.bar.foo",
				vals: []any{"foo", "bar"},
			},
		},
		{
			name: "Key With Nested Sections",
			g:    gitConfig,
			args: args{
				key:  "foo.bar.baz.bla.bla.blu.ble.blo",
				vals: []any{"foo"},
			},
		},
		{
			name: "Variable Name Starts With a Non Alphabetic Character",
			g:    gitConfig,
			args: args{
				key:  "foo.1bar",
				vals: []any{"baz"},
			},
			wantErr: true,
		},
		{
			name: "Boolean Val",
			g:    gitConfig,
			args: args{
				key:  "foo.bool",
				vals: []any{true, false},
			},
		},
		{
			name: "Integer Val",
			g:    gitConfig,
			args: args{
				key:  "foo.int",
				vals: []any{1, 2, 3},
			},
		},
		{
			name: "Invalid Val Type",
			g:    gitConfig,
			args: args{
				key:  "foo.invalidType",
				vals: []any{nil, []int{1, 2, 3}, []any{"aa"}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Set(tt.args.key, tt.args.vals...); (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitConfig_Add(t *testing.T) {
	type args struct {
		key  string
		vals []interface{}
	}
	tests := []struct {
		name    string
		g       *GitConfig
		args    args
		wantErr bool
	}{
		{
			name: "Add To Non-Existing Key",
			g:    gitConfig,
			args: args{
				key:  "foo.uwu",
				vals: []any{"1", "2", 3},
			},
		},
		{
			name: "Add To Existing Key",
			g:    gitConfig,
			args: args{
				key:  "foo.bar",
				vals: []any{"1", 2, "3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Add(tt.args.key, tt.args.vals...); (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitConfig_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		g       *GitConfig
		want    Value
		args    args
		wantErr bool
	}{
		{
			name: "Value Exists",
			g:    gitConfig,
			args: args{
				key: "foo.foo",
			},
			want: Value{"foo"},
		},
		{
			name: "Value Non-Existing Key",
			g:    gitConfig,
			args: args{
				key: "foo.baz",
			},
			wantErr: true,
		},
		{
			name: "Get From Key With Multiple Values",
			g:    gitConfig,
			args: args{
				key: "foo.bar",
			},
			want: Value{"3"},
		},
		{
			name: "Get From Key With Subsection",
			g:    gitConfig,
			args: args{
				key: "foo.bar.baz",
			},
			want: Value{"foo"},
		},
		{
			name: "Get From Key With Nested Sections",
			g:    gitConfig,
			args: args{
				key: "foo.bar.baz.bla.bla.blu.ble.blo",
			},
			want: Value{"foo"},
		},
		{
			name: "Get Non-String Value",
			g:    gitConfig,
			args: args{
				key: "foo.bool",
			},
			want: Value{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitConfig.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitConfig_GetAll(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		g       *GitConfig
		args    args
		want    []Value
		wantErr bool
	}{
		{
			name: "Get Key With Single Value",
			g:    gitConfig,
			args: args{
				key: "foo.foo",
			},
			want: []Value{{"foo"}},
		},
		{
			name: "Get Key With Multiple Values",
			g:    gitConfig,
			args: args{
				key: "foo.bar",
			},
			want: []Value{
				{"foo"},
				{"bar"},
				{"1"},
				{2},
				{"3"},
			},
		},
		{
			name: "Get Non-Existing Key",
			g:    gitConfig,
			args: args{
				key: "foo.baz",
			},
			wantErr: true,
		},
		{
			name: "Get From Key With Subsection",
			g:    gitConfig,
			args: args{
				key: "foo.bar.baz",
			},
			want: []Value{{"foo"}},
		},
		{
			name: "Get From Key With Nested Sections",
			g:    gitConfig,
			args: args{
				key: "foo.bar.baz.bla.bla.blu.ble.blo",
			},
			want: []Value{{"foo"}},
		},
		{
			name: "Get Non-String Value",
			g:    gitConfig,
			args: args{
				key: "foo.int",
			},
			want: []Value{{1}, {2}, {3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetAll(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitConfig.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitConfig_Unset(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		g       *GitConfig
		args    args
		wantErr bool
	}{
		{
			name: "Unset Existing Key",
			g:    gitConfig,
			args: args{
				key: "foo.foo",
			},
		},
		{
			name: "Unset Non-Existing Key",
			g:    gitConfig,
			args: args{
				key: "foo.baz",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Unset(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.Unset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitConfig_Save(t *testing.T) {
	filePath := "./test.gitconfig"
	defer func() {
		os.Remove(filePath)
	}()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		g       *GitConfig
		args    args
		wantErr bool
	}{
		{
			name: "Save to a File",
			g:    gitConfig,
			args: args{
				path: filePath,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Save(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("GitConfig.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitConfig(t *testing.T) {
	gitConfig := New()
	err := gitConfig.Set("foo.bar", "boo")
	if err != nil {
		t.Errorf("gitConfig.Set() = %v, want %v", err, nil)
	}
	err = gitConfig.Set("foo.baz", "blablabla")
	if err != nil {
		t.Errorf("gitConfig.Set() = %v, want %v", err, nil)
	}
	err = gitConfig.Set("bar.foo", "uwu")
	if err != nil {
		t.Errorf("gitConfig.Set() = %v, want %v", err, nil)
	}

	got, err := gitConfig.Get("foo.bar")
	if err != nil || got.String() != "boo" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "boo", nil)
	}
	got, err = gitConfig.Get("foo.baz")
	if err != nil || got.String() != "blablabla" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "blablabla", nil)
	}
	got, err = gitConfig.Get("bar.foo")
	if err != nil || got.String() != "uwu" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "uwu", nil)
	}

	err = gitConfig.Unset("foo.bar")
	if err != nil {
		t.Errorf("gitConfig.Unset() = %v, want %v", err, nil)
	}

	got, err = gitConfig.Get("foo.bar")
	if !errors.Is(err, ErrKeyNotFound) || got.String() == "boo" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "", ErrKeyNotFound)
	}
	got, err = gitConfig.Get("foo.baz")
	if err != nil || got.String() != "blablabla" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "blablabla", nil)
	}

	err = gitConfig.Unset("bar.foo")
	if err != nil {
		t.Errorf("gitConfig.Unset() = %v, want %v", err, nil)
	}
	got, err = gitConfig.Get("bar.foo")
	if !errors.Is(err, ErrKeyNotFound) || got.String() == "uwu" {
		t.Errorf("gitConfig.Get() = (%v, %v), want (%v, %v)", got, err, "", ErrKeyNotFound)
	}
}

func TestNewSection(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Section
		wantErr bool
	}{
		{
			name:    "both invalid",
			args:    args{"w*w.aaaa\n"},
			wantErr: true,
		},
		{
			name:    "invalid section name without subsection",
			args:    args{"foo?"},
			wantErr: true,
		},
		{
			name:    "invalid section name with subsection",
			args:    args{"foo?.bar"},
			wantErr: true,
		},
		{
			name: "valid section without subsection",
			args: args{"foo"},
			want: Section{
				Name: "foo",
			},
		},
		{
			name:    "valid section with invalid subsection",
			args:    args{"foo.bar\n"},
			wantErr: true,
		},
		{
			name: "valid all",
			args: args{"foo.bar"},
			want: Section{
				Name:       "foo",
				Subsection: "bar",
			},
		},
		{
			name: "section name with numbers and dashes",
			args: args{"f-o-o-1-2-3"},
			want: Section{
				Name: "f-o-o-1-2-3",
			},
		},
		{
			name: "weird characters in subsection",
			args: args{"foo.c-l[]o_u^d☁sun☀r-a%d*i+[o]-a$c-t\\i-v^*@&$^*e☢"},
			want: Section{
				Name:       "foo",
				Subsection: "c-l[]o_u^d☁sun☀r-a%d*i+[o]-a$c-t\\i-v^*@&$^*e☢",
			},
		},
		{
			name:    "newline character on subsection",
			args:    args{"foo.bar\n"},
			wantErr: true,
		},
		{
			name: "escaped newline character on subsection",
			args: args{"foo.bar\\n"},
			want: Section{
				Name:       "foo",
				Subsection: "bar\\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSection(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSection() = %v, want %v", got, tt.want)
			}
		})
	}
}
