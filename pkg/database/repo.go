package database

import (
	"errors"
	"github.com/satori/go.uuid"
	"time"
)

var (
	errRecordNotFound = errors.New("record not found")
)

type IFailed struct {
	Id        uuid.UUID `gorm:"type:uuid; primary_key"`
	CreatedAt time.Time
}

type Status struct {
	Id     uint
	Result int
}

type Index struct {
	Id        uint `gorm:"primary_key"`
	Tablename string
	Indexname string
	Indexdef  string
}

func (d *Database) CreateTableIndexes(tName string) error {
	return d.DB.Exec("CREATE TABLE IF NOT EXISTS " + tName + "_indexes (id serial not null constraint " + tName + "_indexes_pkey primary key, tablename varchar(255) not null, indexname varchar(255) not null, indexdef text)").Error
}

func (d *Database) CreateStatusEncodingTable(tName string) error {
	return d.DB.Exec("CREATE TABLE IF NOT EXISTS " + tName + "_status_encoding (id serial not null constraint " + tName + "_status_encoding_pkey primary key, result int default 0)").Error
}

func (d *Database) SetStatusEncodingDefaultValues(tName string) error {
	var count int
	var err error

	err = d.DB.Table(tName + "_status_encoding").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		err = d.DB.Exec("INSERT INTO " + tName + "_status_encoding DEFAULT VALUES").Error // "INSERT INTO " + tName + "_status_encoding DEFAULT VALUES")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) CreateFailedTables(tName string) error {
	return d.DB.Exec("CREATE TABLE IF NOT EXISTS " + tName + "_encoding_failed (id uuid not null, created_at timestamp not null)").Error
}

func (d *Database) Truncate(tableName string) error {
	err := d.DB.Exec("TTRUNCATE TABLE " + tableName).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetCurrentTableIndexs(tName string, indexPkeyName string) ([]Index, error) {
	var indexses []Index
	err := d.DB.Table("pg_indexes").Where("tablename = ? AND (indexname != ? AND indexname NOT LIKE '%_pkey%')", tName, indexPkeyName).Find(&indexses).Error
	if err != nil {
		return indexses, err
	}

	return indexses, nil
}

func (d *Database) SetFailed(tName string, model IFailed) error {
	err := d.DB.Table(tName + "_encoding_failed").Create(&model).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) SetIndex(tNmae string, model *Index) error {
	err := d.DB.Table(tNmae + "_indexes").Save(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetIndex(tNmae string, indexname string) (Index, error) {
	var index Index
	err := d.DB.Table(tNmae+"_indexes").Where("indexname = ?", indexname).Find(&index).Error
	if err != nil {
		if err.Error() != errRecordNotFound.Error() {
			return index, err
		}
	}

	return index, nil
}

func (d *Database) GetAllIndexes(tNmae string) ([]Index, error) {
	var indexses []Index
	err := d.DB.Table(tNmae + "_indexes").Find(&indexses).Error
	if err != nil {
		if err.Error() != errRecordNotFound.Error() {
			return indexses, err
		}
	}

	return indexses, nil
}

func (d *Database) DropIndex(indexName string) error {
	err := d.DB.Exec("DROP INDEX IF EXISTS " + indexName).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) CreatedIndex(query string) error {
	err := d.DB.Exec(query).Error
	if err != nil {
		return err
	}

	return nil
}

// offset
func (d *Database) GetOffset(tName string) (Status, error) {
	var status Status
	err := d.DB.Raw("SELECT id, result FROM " + tName + "_status_encoding WHERE id =1").Scan(&status).Error
	if err != nil {
		return status, err
	}

	return status, nil
}

func (d *Database) UpdateOffset(tName string) error {
	return d.DB.Exec("UPDATE " + tName + "_status_encoding SET result = result + 1").Error
}

// query
func (d *Database) GetTmpLasInsert(tNameTmp string) (string, error) {
	lastInsert := struct {
		ID string
	}{}

	err := d.DB.Raw("SELECT id FROM " + tNameTmp + " WHERE true ORDER BY id ASC LIMIT 1 ").Scan(&lastInsert).Error
	if err != nil {
		if err.Error() != errRecordNotFound.Error() {
			return lastInsert.ID, err
		}
	}

	return lastInsert.ID, nil
}

func (d *Database) GetList(tName string, models interface{}, lasId string, limit int) (interface{}, error) {
	query := d.DB.Table(tName)
	if lasId != "" {
		query = query.Where("id < ?", lasId)
	}

	err := query.Limit(limit).Order("id DESC").Find(models).Error
	if err != nil {
		if err.Error() != errRecordNotFound.Error() {
			return models, err
		}
	}

	return models, nil
}

func (d *Database) InsertTmp(tNameTmp string, model interface{}) error {
	err := d.DB.Table(tNameTmp).Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) RenameTable(fromName string, toName string) error {
	err := d.DB.Exec("ALTER TABLE " + fromName + " RENAME TO " + toName).Error
	if err != nil {
		return err
	}

	return nil
}
