package fixture

import (
	"context"
	"io/ioutil"
	"path/filepath"

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
	// TODO: impl
	return nil
}

func (f *Fixture) exec(model QueryModelWithYaml) error {
	return nil
}

func getFileData(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(ErrFailReadFile, err.Error())
	}
	return data, nil
}
