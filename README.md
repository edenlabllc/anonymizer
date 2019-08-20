# script encoding

1. Used:

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

2. When app runs, it try to read config file from current dir - `./config/config.yaml`

3. Some limitations:

- name of the table in config file must be in `camelcase`
- if column name is `id`, in config file it should be with `- name` in uppercase: `ID`

## Documentation config.yaml

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

4. Example config file can be found in `config` directory
