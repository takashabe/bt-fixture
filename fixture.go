package fixture

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// error variables
var (
	ErrFailRegisterDriver = errors.New("failed to register driver")
	ErrFailReadFile       = errors.New("failed to read file")
	ErrInvalidFixture     = errors.New("invalid fixture file format")
	ErrNotFoundDriver     = errors.New("unknown driver(forgotten import?)")
	ErrUnknownFileExt     = errors.New("unknown file ext")
)

// Fixture provide fixture methods
type Fixture struct {
	client      *bigtable.Client
	adminClient *bigtable.AdminClient
}

// QueryModelWithYaml represent fixture yaml file mapper
type QueryModelWithYaml struct {
	Table          string           `yaml:"table"`
	ColumnFamilies []ColumnFamilies `yaml:"column_families"`
}

type ColumnFamilies struct {
	Family  string    `yaml:"family"`
	Columns []Columns `yaml:"columns"`
}

type Columns struct {
	Key  string            `yaml:"key"`
	Rows map[string]string `yaml:"rows"`
}

// NewFixture returns initialized Fixture
func NewFixture(project, instance string) (*Fixture, error) {
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		return nil, err
	}
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		return nil, err
	}
	return &Fixture{
		client:      client,
		adminClient: adminClient,
	}, nil
}

// Load load .yml script
func (f *Fixture) Load(path string) error {
	data, err := getFileData(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yml", ".yaml":
		return f.loadYaml(data)
	default:
		return errors.Wrapf(ErrUnknownFileExt, "ext:%s, ", ext)
	}
}

func (f *Fixture) loadYaml(file []byte) error {
	model := QueryModelWithYaml{}
	err := yaml.Unmarshal(file, &model)
	if err != nil {
		return errors.Wrapf(ErrInvalidFixture, "%v:, ", err)
	}

	err = f.clearTable(model.Table)
	if err != nil {
		return err
	}
	return f.exec(model)
}

func (f *Fixture) clearTable(table string) error {
	ctx := context.Background()
	tables, err := f.adminClient.Tables(ctx)
	if err != nil {
		return err
	}
	for _, t := range tables {
		if t == table {
			return f.adminClient.DeleteTable(ctx, table)
		}
	}
	return nil
}

func (f *Fixture) exec(model QueryModelWithYaml) error {
	now := time.Now()
	ctx := context.Background()
	if err := f.adminClient.CreateTable(ctx, model.Table); err != nil {
		return err
	}
	table := f.client.Open(model.Table)

	for _, cf := range model.ColumnFamilies {
		fam := cf.Family
		if err := f.adminClient.CreateColumnFamily(ctx, model.Table, fam); err != nil {
			return err
		}

		for _, cs := range cf.Columns {
			var (
				muts = make([]*bigtable.Mutation, 0, len(cs.Rows))
				keys = make([]string, 0, len(cs.Rows))
			)

			for q, v := range cs.Rows {
				mut := bigtable.NewMutation()
				mut.Set(fam, q, bigtable.Time(now), []byte(v))

				muts = append(muts, mut)
				keys = append(keys, cs.Key)
			}

			rowErrs, err := table.ApplyBulk(ctx, keys, muts)
			if err != nil {
				return err
			}
			for _, err := range rowErrs {
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getFileData(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(ErrFailReadFile, err.Error())
	}
	return data, nil
}
