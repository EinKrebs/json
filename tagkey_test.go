// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wal-g/json/mocks"
)

type basicLatin2xTag struct {
	V string `json:"$%-/"`
}

type basicLatin3xTag struct {
	V string `json:"0123456789"`
}

type basicLatin4xTag struct {
	V string `json:"ABCDEFGHIJKLMO"`
}

type basicLatin5xTag struct {
	V string `json:"PQRSTUVWXYZ_"`
}

type basicLatin6xTag struct {
	V string `json:"abcdefghijklmno"`
}

type basicLatin7xTag struct {
	V string `json:"pqrstuvwxyz"`
}

type miscPlaneTag struct {
	V string `json:"色は匂へど"`
}

type percentSlashTag struct {
	V string `json:"text/html%"` // https://golang.org/issue/2718
}

type punctuationTag struct {
	V string `json:"!#$%&()*+-./:;<=>?@[]^_{|}~ "` //nolint:staticcheck // https://golang.org/issue/3546
}

type dashTag struct {
	V string `json:"-,"`
}

type emptyTag struct {
	W string
}

type misnamedTag struct {
	X string `jsom:"Misnamed"`
}

type badFormatTag struct {
	Y string `:"BadFormat"` //nolint:govet
}

type badCodeTag struct {
	Z string `json:" !\"#&'()*+,."` //nolint: staticcheck
}

type spaceTag struct {
	Q string `json:"With space"`
}

type unicodeTag struct {
	W string `json:"Ελλάδα"`
}

var structTagObjectKeyTests = []struct {
	raw   interface{}
	value string
	key   string
}{
	{basicLatin2xTag{"2x"}, "2x", "$%-/"},
	{basicLatin3xTag{"3x"}, "3x", "0123456789"},
	{basicLatin4xTag{"4x"}, "4x", "ABCDEFGHIJKLMO"},
	{basicLatin5xTag{"5x"}, "5x", "PQRSTUVWXYZ_"},
	{basicLatin6xTag{"6x"}, "6x", "abcdefghijklmno"},
	{basicLatin7xTag{"7x"}, "7x", "pqrstuvwxyz"},
	{miscPlaneTag{"いろはにほへと"}, "いろはにほへと", "色は匂へど"},
	{dashTag{"foo"}, "foo", "-"},
	{emptyTag{"Pour Moi"}, "Pour Moi", "W"},
	{misnamedTag{"Animal Kingdom"}, "Animal Kingdom", "X"},
	{badFormatTag{"Orfevre"}, "Orfevre", "Y"},
	{badCodeTag{"Reliable Man"}, "Reliable Man", "Z"},
	{percentSlashTag{"brut"}, "brut", "text/html%"},
	{punctuationTag{"Union Rags"}, "Union Rags", "!#$%&()*+-./:;<=>?@[]^_{|}~ "},
	{spaceTag{"Perreddu"}, "Perreddu", "With space"},
	{unicodeTag{"Loukanikos"}, "Loukanikos", "Ελλάδα"},
}

func TestStructTagObjectKey(t *testing.T) {
	buf := mocks.NewBuildCloser()
	for _, tt := range structTagObjectKeyTests {
		require.NoError(t, Marshal(tt.raw, buf))
		var f interface{}
		require.NoError(t, Unmarshal(strings.NewReader(buf.String()), &f))
		for i, v := range f.(map[string]interface{}) {
			switch i {
			case tt.key:
				s, ok := v.(string)
				require.True(t, ok)
				require.Equal(t, tt.value, s)
			default:
				t.Fatalf("Unexpected key: %#q, from %#q", i, buf.String())
			}
		}
		buf.Reset()
	}
}
