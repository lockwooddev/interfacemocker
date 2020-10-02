package interfacemocker

import (
	"fmt"
	"reflect"
	"strings"
)

type Foo struct{}
type Bar struct{}
type Baz struct{}
type Waldo struct{}

type Mocker interface {
	GetFoo() (Foo, error)
	GetBar() Bar
	GetBaz() *Baz
	GetWaldo() (*Waldo, error)
}

type mocker struct {
	iter  int
	Mocks []interface{}
}

func (m *mocker) GetFoo() (Foo, error) {
	value := m.getMultiMock("interfacemocker.Foo", "errors.errorString")
	if value[1] == nil {
		return value[0].(Foo), nil
	}
	return value[0].(Foo), value[1].(error)
}

func (m *mocker) GetBar() Bar {
	value := m.getSingleMock("interfacemocker.Bar")
	return value.(Bar)
}

func (m *mocker) GetBaz() *Baz {
	value := m.getSingleMock("interfacemocker.Baz")
	return value.(*Baz)
}

func (m *mocker) GetWaldo() (*Waldo, error) {
	value := m.getMultiMock("interfacemocker.Waldo", "errors.errorString")
	if value[1] == nil {
		return value[0].(*Waldo), nil
	}
	return value[0].(*Waldo), value[1].(error)
}

func (m mocker) checkIteration() {
	if m.iter == len(m.Mocks) {
		panic("no more mocks available")
	}
}

func (m *mocker) getMultiMock(types ...string) []interface{} {
	m.checkIteration()

	mock := m.Mocks[m.iter]
	switch reflect.TypeOf(mock).Kind() {
	case reflect.Slice:
		mockValues := mock.([]interface{})
		if len(types) != len(mockValues) {
			msg := fmt.Sprintf(
				"mocker expected %d values, but got %d on call %d",
				len(types), len(mockValues), m.iter,
			)
			panic(msg)
		}

		for idx, mockValue := range mockValues {
			expectedType := types[idx]

			// nill values should be skipped
			if mockValue == nil {
				continue
			}

			mockType := reflect.TypeOf(mockValue).String()
			if !strings.Contains(mockType, expectedType) {
				msg := fmt.Sprintf(
					"mocker expected '%s', but got '%s' on call %d and position %d",
					expectedType, mockType, m.iter, idx,
				)
				panic(msg)
			}
		}
		m.iter++
		return mock.([]interface{})
	default:
		panic(fmt.Sprintf("mocker expected slice, but got '%s' on call %d", reflect.TypeOf(mock).String(), m.iter))
	}
}

func (m *mocker) getSingleMock(t string) interface{} {
	m.checkIteration()

	mock := m.Mocks[m.iter]
	if strings.Contains(reflect.TypeOf(mock).String(), t) {
		m.iter++
		return mock
	}
	panic(fmt.Sprintf("Mocker expected '%s' return value, but got '%s' on call %d", reflect.TypeOf(mock).String(), t, m.iter))
}

func NewMocker(mocks []interface{}) Mocker {
	return &mocker{Mocks: mocks}
}
