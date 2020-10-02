package interfacemocker

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mocker_success(t *testing.T) {
	mocks := []interface{}{
		[]interface{}{
			Foo{}, errors.New("foo"),
		},
		Bar{},
		[]interface{}{
			Foo{}, nil,
		},
		// test pointer
		&Baz{},
		// test multi value return pointer with error
		[]interface{}{
			&Waldo{}, errors.New("foo"),
		},
		// test multi value return with nil error
		[]interface{}{
			&Waldo{}, nil,
		},
	}
	mocker := NewMocker(mocks)

	r1, err := mocker.GetFoo()
	assert.Error(t, err)
	assert.Equal(t, reflect.TypeOf(r1).String(), "interfacemocker.Foo")

	r2 := mocker.GetBar()
	assert.Equal(t, reflect.TypeOf(r2).String(), "interfacemocker.Bar")

	r3, err := mocker.GetFoo()
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(r3).String(), "interfacemocker.Foo")

	r4 := mocker.GetBaz()
	assert.Equal(t, reflect.TypeOf(r4).String(), "*interfacemocker.Baz")

	r5, err := mocker.GetWaldo()
	assert.Error(t, err)
	assert.Equal(t, reflect.TypeOf(r5).String(), "*interfacemocker.Waldo")

	r6, err := mocker.GetWaldo()
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(r6).String(), "*interfacemocker.Waldo")
}

// test panic when wrong type is returned, not matching the mocks list
func Test_Mocker_getSingleMock_panics(t *testing.T) {
	mocks := []interface{}{Bar{}, &Baz{}}
	mocker := NewMocker(mocks)

	r1 := mocker.GetBar()
	assert.Equal(t, reflect.TypeOf(r1).String(), "interfacemocker.Bar")
	expectedPanicMsg := "Mocker expected '*interfacemocker.Baz' return value, but got 'interfacemocker.Bar' on call 1"
	assert.PanicsWithValue(t, expectedPanicMsg, func() { mocker.GetBar() })
}

func Test_Mocker_checkIteration_panics(t *testing.T) {
	mocks := []interface{}{}
	mocker := NewMocker(mocks)
	expectPanicMsg := "no more mocks available"

	assert.PanicsWithValue(t, expectPanicMsg, func() { mocker.GetFoo() })
}

func Test_Mocker_getMultiMock_not_a_slice_panic(t *testing.T) {
	mocks := []interface{}{Bar{}}
	mocker := NewMocker(mocks)

	assert.PanicsWithValue(t, "mocker expected slice, but got 'interfacemocker.Bar' on call 0", func() { mocker.GetFoo() })
}

func Test_Mocker_getMultiMock_invalid_mock_slice_length_panic(t *testing.T) {
	mocks := []interface{}{
		[]interface{}{
			Foo{},
		},
	}
	mocker := NewMocker(mocks)

	assert.PanicsWithValue(t, "mocker expected 2 values, but got 1 on call 0", func() { mocker.GetFoo() })
}

func Test_Mocker_getMultiMock_invalid_mock_return_value_panic(t *testing.T) {
	mocks := []interface{}{
		[]interface{}{
			Foo{}, nil,
		},
	}
	mocker := NewMocker(mocks)

	assert.PanicsWithValue(t, "mocker expected 'interfacemocker.Waldo', but got 'interfacemocker.Foo' on call 0 and position 0", func() { mocker.GetWaldo() })
}
