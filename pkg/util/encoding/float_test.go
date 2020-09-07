// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package encoding

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/weisslj/cockroach/pkg/util/randutil"
)

func TestEncodeFloatOrdered(t *testing.T) {
	testCases := []struct {
		Value    float64
		Encoding []byte
	}{
		{math.NaN(), []byte{0x02}},
		{math.Inf(-1), []byte{0x03, 0x00, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{-math.MaxFloat64, []byte{0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{-1e308, []byte{0x03, 0x00, 0x1e, 0x33, 0x0c, 0x7a, 0x14, 0x37, 0x5f}},
		{-10000.0, []byte{0x03, 0x3f, 0x3c, 0x77, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{-9999.0, []byte{0x03, 0x3f, 0x3c, 0x78, 0x7f, 0xff, 0xff, 0xff, 0xff}},
		{-100.0, []byte{0x03, 0x3f, 0xa6, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{-99.0, []byte{0x03, 0x3f, 0xa7, 0x3f, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{-1.0, []byte{0x03, 0x40, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{-0.00123, []byte{0x03, 0x40, 0xab, 0xd9, 0x01, 0x8e, 0x75, 0x79, 0x28}},
		{-1e-307, []byte{0x03, 0x7f, 0xce, 0x05, 0xe7, 0xd3, 0xbf, 0x39, 0xf2}},
		{-math.SmallestNonzeroFloat64, []byte{0x03, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}},
		{math.Copysign(0, -1), []byte{0x04}},
		{0, []byte{0x04}},
		{math.SmallestNonzeroFloat64, []byte{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{1e-307, []byte{0x05, 0x00, 0x31, 0xfa, 0x18, 0x2c, 0x40, 0xc6, 0x0d}},
		{0.00123, []byte{0x05, 0x3f, 0x54, 0x26, 0xfe, 0x71, 0x8a, 0x86, 0xd7}},
		{0.0123, []byte{0x05, 0x3f, 0x89, 0x30, 0xbe, 0x0d, 0xed, 0x28, 0x8d}},
		{0.123, []byte{0x05, 0x3f, 0xbf, 0x7c, 0xed, 0x91, 0x68, 0x72, 0xb0}},
		{1.0, []byte{0x05, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{10.0, []byte{0x05, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{12.345, []byte{0x05, 0x40, 0x28, 0xb0, 0xa3, 0xd7, 0x0a, 0x3d, 0x71}},
		{99.0, []byte{0x05, 0x40, 0x58, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{99.0001, []byte{0x05, 0x40, 0x58, 0xc0, 0x01, 0xa3, 0x6e, 0x2e, 0xb2}},
		{99.01, []byte{0x05, 0x40, 0x58, 0xc0, 0xa3, 0xd7, 0x0a, 0x3d, 0x71}},
		{100.0, []byte{0x05, 0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{100.01, []byte{0x05, 0x40, 0x59, 0x00, 0xa3, 0xd7, 0x0a, 0x3d, 0x71}},
		{100.1, []byte{0x05, 0x40, 0x59, 0x06, 0x66, 0x66, 0x66, 0x66, 0x66}},
		{1234, []byte{0x05, 0x40, 0x93, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{1234.5, []byte{0x05, 0x40, 0x93, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{9999, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x00, 0x00, 0x00, 0x00}},
		{9999.000001, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x00, 0x08, 0x63, 0x7c}},
		{9999.000009, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x00, 0x4b, 0x7f, 0x5a}},
		{9999.00001, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x00, 0x53, 0xe2, 0xd6}},
		{9999.00009, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x02, 0xf2, 0xf9, 0x87}},
		{9999.000099, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x03, 0x3e, 0x78, 0xe2}},
		{9999.0001, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x03, 0x46, 0xdc, 0x5d}},
		{9999.001, []byte{0x05, 0x40, 0xc3, 0x87, 0x80, 0x20, 0xc4, 0x9b, 0xa6}},
		{9999.01, []byte{0x05, 0x40, 0xc3, 0x87, 0x81, 0x47, 0xae, 0x14, 0x7b}},
		{9999.1, []byte{0x05, 0x40, 0xc3, 0x87, 0x8c, 0xcc, 0xcc, 0xcc, 0xcd}},
		{10000, []byte{0x05, 0x40, 0xc3, 0x88, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{10001, []byte{0x05, 0x40, 0xc3, 0x88, 0x80, 0x00, 0x00, 0x00, 0x00}},
		{12345, []byte{0x05, 0x40, 0xc8, 0x1c, 0x80, 0x00, 0x00, 0x00, 0x00}},
		{123450, []byte{0x05, 0x40, 0xfe, 0x23, 0xa0, 0x00, 0x00, 0x00, 0x00}},
		{1e308, []byte{0x05, 0x7f, 0xe1, 0xcc, 0xf3, 0x85, 0xeb, 0xc8, 0xa0}},
		{math.MaxFloat64, []byte{0x05, 0x7f, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{math.Inf(1), []byte{0x05, 0x7f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	var lastEncoded []byte
	for _, dir := range []Direction{Ascending, Descending} {
		for i, c := range testCases {
			var enc []byte
			var err error
			var dec float64
			if dir == Ascending {
				enc = EncodeFloatAscending(nil, c.Value)
				_, dec, err = DecodeFloatAscending(enc)
			} else {
				enc = EncodeFloatDescending(nil, c.Value)
				_, dec, err = DecodeFloatDescending(enc)
			}
			if dir == Ascending && !bytes.Equal(enc, c.Encoding) {
				t.Errorf("unexpected mismatch for %v. expected [% x], got [% x]",
					c.Value, c.Encoding, enc)
			}
			if i > 0 {
				if (bytes.Compare(lastEncoded, enc) > 0 && dir == Ascending) ||
					(bytes.Compare(lastEncoded, enc) < 0 && dir == Descending) {
					t.Errorf("%v: expected [% x] to be less than or equal to [% x]",
						c.Value, testCases[i-1].Encoding, enc)
				}
			}
			if err != nil {
				t.Error(err)
				continue
			}
			if math.IsNaN(c.Value) {
				if !math.IsNaN(dec) {
					t.Errorf("unexpected mismatch for %v. got %v", c.Value, dec)
				}
			} else if dec != c.Value {
				t.Errorf("unexpected mismatch for %v. got %v", c.Value, dec)
			}
			testPeekLength(t, enc)
			lastEncoded = enc
		}

		// Test that appending the float to an existing buffer works.
		var enc []byte
		var dec float64
		if dir == Ascending {
			enc = EncodeFloatAscending([]byte("hello"), 1.23)
			_, dec, _ = DecodeFloatAscending(enc[5:])
		} else {
			enc = EncodeFloatDescending([]byte("hello"), 1.23)
			_, dec, _ = DecodeFloatDescending(enc[5:])
		}
		if dec != 1.23 {
			t.Errorf("unexpected mismatch for %v. got %v", 1.23, dec)
		}
	}
}

// TestEncodeFloatValue tests that EncodeFloatValue is bit-for-bit equal
// after decoding. Important for the 0/-0 case because 0 == -0, but we need
// to verify that -0 came out of the decoder.
func TestEncodeFloatValue(t *testing.T) {
	testCases := []float64{
		math.NaN(),
		math.Inf(-1),
		-1,
		math.Copysign(0, -1),
		0,
		1,
		math.Inf(1),
	}

	for _, f := range testCases {
		t.Run(fmt.Sprint(f), func(t *testing.T) {
			enc := EncodeFloatValue(nil, NoColumnID, f)
			_, dec, err := DecodeFloatValue(enc)
			if err != nil {
				t.Fatal(err)
			}
			testBits := math.Float64bits(f)
			encBits := math.Float64bits(dec)
			if testBits != encBits {
				t.Fatalf("unexpected decoded float64 value. expected %f, got %f", f, dec)
			}
		})
	}
}

func BenchmarkEncodeFloat(b *testing.B) {
	rng, _ := randutil.NewPseudoRand()

	vals := make([]float64, 10000)
	for i := range vals {
		vals[i] = rng.Float64()
	}

	buf := make([]byte, 0, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeFloatAscending(buf, vals[i%len(vals)])
	}
}

func BenchmarkDecodeFloat(b *testing.B) {
	rng, _ := randutil.NewPseudoRand()

	vals := make([][]byte, 10000)
	for i := range vals {
		vals[i] = EncodeFloatAscending(nil, rng.Float64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DecodeFloatAscending(vals[i%len(vals)])
	}
}
