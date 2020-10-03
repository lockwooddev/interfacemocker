# interfacemocker

I use a lot of dependency injection in Golang. This repo contains a pattern to make your own mocker for interfaces in case you want to bypass external services like a database in tests. The pattern specializes on mocking the return values for single or multiple object types.

## Usage example

When you adapt the pattern in `mocks.go`, you will be able to easily create a list of mocks that should correspond to the calls you expect in your testcase.

```golang
package example

type DatabaseMocker interface {
	DriverName() string
	GetRow() (Row, error)
}

type mocker struct {
	iter  int
	Mocks []interface{}
}

func NewMocker(mocks []interface{}) Mocker {
	return &mocker{Mocks: mocks}
}

func (m *mocker) DriverName() string {
	value := m.getSingleMock("string")
	return value.(string)
}

func (m *mocker) GetRow() (db.Row, error) {
	value := m.getMultiMock("db.Row", "errors.errorString")
	if value[1] == nil {
		return value[0].(db.Row), nil
	}
	return value[0].(db.Row), value[1].(error)
}

mocks := []interface{}{
	[]interface{}{
		Row{}, errors.New("not found"),
	},
	[]interface{}{
		Row{}, nil,
	},
	"psql"
}
mocker := NewMocker(mocks)

row1, err := mocker.GetRow()
row2, err := mocker.GetRow()
row3, err := mocker.GetRow()
// panic
```

## footnotes

1. Panics are intended for testcases where your mock return value does not match the return value of the actual call. For every testcase you need to tune your mock return values.
