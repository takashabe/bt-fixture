package fixture

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	now := time.Now()
	cases := []struct {
		input  QueryModelWithYaml
		expect error
	}{
		{
			QueryModelWithYaml{
				Table: "test",
				ColumnFamilies: []ColumnFamilies{
					ColumnFamilies{
						Family: "d",
						Columns: []Columns{
							Columns{
								Key: "1",
								Rows: map[string]interface{}{
									"name": "foo",
									"age":  "1",
								},
								Version: Version{now},
							},
						},
					},
				},
			},
			nil,
		},
	}
	for _, c := range cases {
		f, err := NewFixture("test-project", "test-instance")
		assert.NoError(t, err)

		err = f.clearTable(c.input.Table)
		assert.NoError(t, err)

		err = f.exec(c.input)
		assert.Equal(t, c.expect, errors.Cause(err))
	}
}

func TestLoad(t *testing.T) {
	cases := []struct {
		input  string
		expect error
	}{
		{"testdata/test.yml", nil},
		{"testdata/version.yaml", nil},
		{"testdata/invalid.ext", ErrUnknownFileExt},
		{"testdata/invalid.yaml", ErrInvalidFixture},
		{"not_exists.yml", ErrFailReadFile},
	}
	for _, c := range cases {
		f, err := NewFixture("test-project", "test-instance")
		assert.NoError(t, err)

		err = f.Load(c.input)
		assert.Equal(t, c.expect, errors.Cause(err))
	}
}
