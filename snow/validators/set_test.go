// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ava-labs/gecko/ids"
)

func TestSetSet(t *testing.T) {
	vdr0 := NewValidator(ids.ShortEmpty, 1)
	vdr1 := NewValidator(ids.NewShortID([20]byte{0xFF}), math.MaxInt64-1)
	// Should be discarded, because it has a weight of 0
	vdr2 := NewValidator(ids.NewShortID([20]byte{0xAA}), 0)

	s := NewSet()
	err := s.Set([]Validator{vdr0, vdr1, vdr2})
	assert.NoError(t, err)

	length := s.Len()
	assert.Equal(t, 2, length, "should have two validators")

	contains := s.Contains(vdr0.ID())
	assert.True(t, contains, "should have contained vdr0")

	contains = s.Contains(vdr1.ID())
	assert.True(t, contains, "should have contained vdr1")

	sampled, err := s.Sample(1)
	assert.NoError(t, err)
	assert.Len(t, sampled, 1, "should have only sampled one validator")
	assert.Equal(t, vdr1.ID(), sampled[0].ID(), "should have sampled vdr1")
}

func TestSamplerSample(t *testing.T) {
	vdr0 := GenerateRandomValidator(1)
	vdr1 := GenerateRandomValidator(math.MaxInt64 - 1)

	s := NewSet()
	err := s.Add(vdr0)
	assert.NoError(t, err)

	sampled, err := s.Sample(1)
	assert.NoError(t, err)
	assert.Len(t, sampled, 1, "should have only sampled one validator")
	assert.Equal(t, vdr0.ID(), sampled[0].ID(), "should have sampled vdr0")

	_, err = s.Sample(2)
	assert.Error(t, err, "should have errored during sampling")

	err = s.Add(vdr1)
	assert.NoError(t, err)

	sampled, err = s.Sample(1)
	assert.NoError(t, err)
	assert.Len(t, sampled, 1, "should have only sampled one validator")
	assert.Equal(t, vdr1.ID(), sampled[0].ID(), "should have sampled vdr1")

	sampled, err = s.Sample(2)
	assert.NoError(t, err)
	assert.Len(t, sampled, 2, "should have sampled two validators")
	assert.Equal(t, vdr1.ID(), sampled[0].ID(), "should have sampled vdr1")
	assert.Equal(t, vdr1.ID(), sampled[1].ID(), "should have sampled vdr1")

	sampled, err = s.Sample(3)
	assert.NoError(t, err)
	assert.Len(t, sampled, 3, "should have sampled three validators")
	assert.Equal(t, vdr1.ID(), sampled[0].ID(), "should have sampled vdr1")
	assert.Equal(t, vdr1.ID(), sampled[1].ID(), "should have sampled vdr1")
	assert.Equal(t, vdr1.ID(), sampled[2].ID(), "should have sampled vdr1")
}

func TestSamplerDuplicate(t *testing.T) {
	vdr0 := GenerateRandomValidator(1)
	vdr1_0 := GenerateRandomValidator(math.MaxInt64 - 1)
	vdr1_1 := NewValidator(vdr1_0.ID(), 0)

	s := NewSet()
	err := s.Add(vdr0)
	assert.NoError(t, err)

	err = s.Add(vdr1_0)
	assert.NoError(t, err)

	sampled, err := s.Sample(1)
	assert.NoError(t, err)
	assert.Len(t, sampled, 1, "should have only sampled one validator")
	assert.Equal(t, vdr1_0.ID(), sampled[0].ID(), "should have sampled vdr1")

	err = s.Add(vdr1_1)
	assert.NoError(t, err)

	sampled, err = s.Sample(1)
	assert.NoError(t, err)
	assert.Len(t, sampled, 1, "should have only sampled one validator")
	assert.Equal(t, vdr0.ID(), sampled[0].ID(), "should have sampled vdr0")
}

func TestSamplerContains(t *testing.T) {
	vdr := GenerateRandomValidator(1)

	s := NewSet()
	err := s.Add(vdr)
	assert.NoError(t, err)

	contains := s.Contains(vdr.ID())
	assert.True(t, contains, "should have contained validator")

	err = s.Remove(vdr.ID())
	assert.NoError(t, err)

	contains = s.Contains(vdr.ID())
	assert.False(t, contains, "shouldn't have contained validator")
}

func TestSamplerString(t *testing.T) {
	vdr0 := NewValidator(ids.ShortEmpty, 1)
	vdr1 := NewValidator(
		ids.NewShortID([20]byte{
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		}),
		math.MaxInt64-1,
	)

	s := NewSet()
	err := s.Add(vdr0)
	assert.NoError(t, err)

	err = s.Add(vdr1)
	assert.NoError(t, err)

	expected := "Validator Set: (Size = 2)\n" +
		"    Validator[0]:        111111111111111111116DBWJs, 1\n" +
		"    Validator[1]: QLbz7JHiBTspS962RLKV8GndWFwdYhk6V, 9223372036854775806"
	result := s.String()
	assert.Equal(t, expected, result, "wrong string returned")
}

func TestSetWeight(t *testing.T) {
	weight0 := uint64(93)
	vdr0 := NewValidator(ids.NewShortID([20]byte{1}), weight0)
	weight1 := uint64(123)
	vdr1 := NewValidator(ids.NewShortID([20]byte{2}), weight1)

	s := NewSet()
	err := s.Add(vdr0)
	assert.NoError(t, err)

	err = s.Add(vdr1)
	assert.NoError(t, err)

	setWeight := s.Weight()
	expectedWeight := weight0 + weight1
	assert.Equal(t, expectedWeight, setWeight, "wrong set weight")
}
