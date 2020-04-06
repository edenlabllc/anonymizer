# Anonymizer for postgrsql

![Anonim](https://ic.pics.livejournal.com/hannibal_md/26923335/87805/87805_600.png)

This application created for anonimiztion of the data in the postgres tables.

## features

- Use gorutines for parallelism
- Simple config in yaml format
- Can be used without compilation
- Can be dockerize

## how it works

- Creates new table with suffix `_tmp` (example `original_name_tmp`)
- Create new table with sufix `_indexes` its contains all indexes from original table (example `original_name_indexes`)
- Take original record and transformed it insert in new table with sufix `_tmp`
- After all records from original table are done aplication rename original tabel (add sufix `_timestamp`, example: `original_name_819636251`) after that rename `_tmp` tabel to original name
- Apply all indexes from `_indexes` for renamed `_tmp` table

## script encoding

1. Create config file `./config/config.yaml` (example can be found `config/__config.yaml`)

2. Use without compilation:

```shell
go run ./main.go -user=user -pass=password -name=dbname -host=localhost -port=25432 -limit=100 -tname=table_name
```

- `-user` - databes user name
- `-pass` - database password
- `-name` - database name
- `-tname` - database table name
- `-host` - database host
- `-port` - database port
- `-limit` - sql query limit rows (default: 100)
- `-concurrency` - numbers workers goroutine (default: 10)

3. Some limitations:

- name of the table in config file must be in `camelcase` (exempl: OneField)
- if column name is `id`, in config file it should be with `- name` in uppercase: `ID`

## config.yaml options

### 1. Field Types

- `uuid`
- `varchar`
- `string`
- `int`
- `jsonb`
- `timestamp`
- `time`
- `date`
- `boolean`
- `text`
- `decimal`
- `double`
- `float`
- `timeuuid`
- `smallint`
- `jsonbMap`
- `nullString`
- `nullBool`
- `nullInt`
- `nullFloat`

### 2. Tags 'generated'

- set **null**

 ```yaml
  name: "TaxId"
  type: "nullString"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:null\"`"
```

- set **hex**

```yaml
  name: "FirstName"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:hex\"`"
```

- set **md5Hash**

```yaml
  name: "FirstName"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:md5Hash\"`"
```

- set **customString**

```yaml
  name: "Gender"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:customString;value:GenderTest\"`"
```

- set **randString** (default: size(1000000000, 9999999999))

```yaml
  name: "SecondName"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:randString;size:(1000000000, 9999999999)\"`"
```

- set **randStringDouble** (default: size:(10000000, 99999999),(10000, 99999))

```yaml
  name: "LastName"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:randStringDouble;size:(10000000, 99999999),(10000, 99999)\"`"
```

- set **randIntAndString** (default: size_int: 2 size:(100000, 999999))

```yaml
  name: "FirstName"
  type: "varchar"
  tag: "`gorm:\"type: varchar(255)\" sql:\"type:varchar(255)\" generated:\"type:randIntAndString;size_int: 2;size:(100000, 999999)\"`"
```

- set **phoneRandom** (default: size:(10000000000, 99999999999))

```yaml
    name: "AboutMyself"
    type: "text"
    tag: "`gorm:\"type: text\" generated:\"type:phoneRandom;size:(10000000000, 99999999999)\"`"
```
