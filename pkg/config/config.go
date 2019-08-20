package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

var (
	ErrNotEmpty = errors.New("Can not be empty")
)

type Config struct {
	Databases []Database `yaml:"databases"`
}

func (c *Config) IsDatabase(name string) bool {
	for _, database := range c.Databases {
		if database.Name != name {
			return false
		}
	}

	return true
}

func (c *Config) GetDatabase(name string) Database {
	for _, database := range c.Databases {
		if database.Name == name {
			return database
		}
	}

	return Database{}
}

type Database struct {
	Name      string     `yaml:"name"`
	Tables    []Table    `yaml:"tables"`
	Truncates []Truncate `yaml:"truncates"`
}

func (d *Database) IsTable(tname string) bool {
	for _, table := range d.Tables {
		if table.Name == tname {
			return true
		}
	}

	return false
}

func (d *Database) GetTable(tname string) Table {
	for _, table := range d.Tables {
		if table.Name == tname {
			return table
		}
	}

	return Table{}
}

func (d *Database) GetTruncates() []Truncate {
	var truncates []Truncate

	if len(d.Truncates) == 0 {
		return truncates
	}

	for _, truncate := range d.Truncates {
		if truncate.Name != "" {
			truncates = append(truncates, truncate)
		}
	}

	return truncates
}

type Truncate struct {
	Name string `yaml:"name"`
}

type Table struct {
	Name     string  `yaml:"name"`
	Lang     string  `yaml:"lang,omitempty"`
	TmpName  string  `yaml:"tmp_name"`
	PkeyName string  `yaml:"pkey_name"`
	Fields   []Field `yaml:"fields"`
}

func (t *Table) Validate() error {
	if t.Name == "" || t.TmpName == "" {
		return ErrNotEmpty
	}

	return nil
}

func (t *Table) CheckFields() []error {
	var listErrors []error
	for _, field := range t.Fields {
		if err := field.Validate(); err != nil {
			listErrors = append(listErrors, err)
		}
	}

	return listErrors
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Tag  string `yaml:"tag"`
}

func (f *Field) Validate() error {
	if f.Name == "" || f.Type == "" || f.Tag == "" {
		return ErrNotEmpty
	}

	return nil
}

func InitYaml() (Config, error) {
	var config Config
	filename, _ := filepath.Abs("./config/config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
