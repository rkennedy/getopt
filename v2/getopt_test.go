// Copyright 2020 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package getopt

import (
	"slices"
	"testing"
)

func TestAnyOrder(t *testing.T) {
	for _, tt := range []struct {
		name      string
		in        []string
		anyOrder  bool
		err       string
		remaining []string
	}{
		{
			name:      "allow any order",
			in:        []string{"test", "-a", "arg", "-b"},
			anyOrder:  true,
			err:       "",
			remaining: []string{"arg"},
		},
		{
			name:      "disallow any order",
			in:        []string{"test", "-a", "arg", "-b"},
			anyOrder:  false,
			err:       "test: option -b is mandatory",
			remaining: []string{"arg", "-b"},
		},
		{
			name:      "exclude dash-dash",
			in:        []string{"test", "-a", "--", "arg", "-b"},
			anyOrder:  true,
			err:       "test: option -b is mandatory",
			remaining: []string{"arg", "-b"},
		},
	} {
		reset()
		AllowAnyOrder(tt.anyOrder)
		var vala, valb bool
		Flag(&vala, 'a')
		Flag(&valb, 'b').Mandatory()
		parse(tt.in)
		if s := checkError(tt.err); s != "" {
			t.Errorf("%s: %s", tt.name, s)
		}
		if !slices.Equal(tt.remaining, Args()) {
			t.Errorf("%s: expected %+v but got %+v", tt.name, tt.remaining, Args())
		}
	}
}

func TestMandatory(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		err  string
	}{
		{
			name: "required option present",
			in:   []string{"test", "-r"},
		},
		{
			name: "required option not present",
			in:   []string{"test", "-o"},
			err:  "test: option -r is mandatory",
		},
		{
			name: "no options",
			in:   []string{"test"},
			err:  "test: option -r is mandatory",
		},
	} {
		reset()
		var val bool
		Flag(&val, 'o')
		Flag(&val, 'r').Mandatory()
		parse(tt.in)
		if s := checkError(tt.err); s != "" {
			t.Errorf("%s: %s", tt.name, s)
		}
	}
}

func TestGroup(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		err  string
	}{
		{
			name: "no args",
			in:   []string{"test"},
			err:  "test: exactly one of the following options must be specified: -A, -B",
		},
		{
			name: "one of each",
			in:   []string{"test", "-A", "-C"},
		},
		{
			name: "Two in group One",
			in:   []string{"test", "-A", "-B"},
			err:  "test: options -A and -B are mutually exclusive",
		},
		{
			name: "Two in group Two",
			in:   []string{"test", "-A", "-D", "-C"},
			err:  "test: options -C and -D are mutually exclusive",
		},
	} {
		reset()
		var val bool
		Flag(&val, 'o')
		Flag(&val, 'A').SetGroup("One")
		Flag(&val, 'B').SetGroup("One")
		Flag(&val, 'C').SetGroup("Two")
		Flag(&val, 'D').SetGroup("Two")
		RequiredGroup("One")
		parse(tt.in)
		if s := checkError(tt.err); s != "" {
			t.Errorf("%s: %s", tt.name, s)
		}
	}
}
