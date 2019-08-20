package generated

import (
	"database/sql"
	"ehealth-migration/pkg/database"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"
	"math/big"
	"reflect"
	"time"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type GeneratedStruct struct {
	Fields []Field
}

func MakeSlices(model interface{}) interface{} {
	v := reflect.ValueOf(model)
	fmt.Println(v.Type())
	sliceType := reflect.SliceOf(v.Type())
	emptySlice := reflect.MakeSlice(sliceType, 10, 10)

	x := reflect.New(emptySlice.Type())
	x.Elem().Set(emptySlice)

	return x.Interface()
}

func (gs *GeneratedStruct) ToStruct() interface{} {
	var structFields []reflect.StructField
	for _, field := range gs.Fields {
		structFields = append(structFields, reflect.StructField{
			Name: field.Name,
			Type: goType(getBaseType(field.Type)),
			Tag:  reflect.StructTag(field.Tag),
		})
	}

	typ := reflect.StructOf(structFields)
	v := reflect.New(typ).Elem()
	s := v.Addr().Interface()

	return s
}

type Type int

const (
	TypeCustom     Type = 0x0000
	TypeAscii      Type = 0x0001
	TypeBigInt     Type = 0x0002
	TypeBlob       Type = 0x0003
	TypeBoolean    Type = 0x0004
	TypeCounter    Type = 0x0005
	TypeDecimal    Type = 0x0006
	TypeDouble     Type = 0x0007
	TypeFloat      Type = 0x0008
	TypeInt        Type = 0x0009
	TypeText       Type = 0x000A
	TypeTimestamp  Type = 0x000B
	TypeUUID       Type = 0x000C
	TypeVarchar    Type = 0x000D
	TypeVarint     Type = 0x000E
	TypeTimeUUID   Type = 0x000F
	TypeInet       Type = 0x0010
	TypeDate       Type = 0x0011
	TypeTime       Type = 0x0012
	TypeSmallInt   Type = 0x0013
	TypeTinyInt    Type = 0x0014
	TypeDuration   Type = 0x0015
	TypeList       Type = 0x0020
	TypeMap        Type = 0x0021
	TypeSet        Type = 0x0022
	TypeUDT        Type = 0x0030
	TypeJsonB      Type = 0x0032
	TypeJsonBMap   Type = 0x0033
	TypeNullString Type = 0x0034
	TypeNullBool   Type = 0x0035
	TypeNullInt    Type = 0x0036
	TypeNullFloat  Type = 0x0037
	TypeDateIsNull Type = 0x0038
)

// String returns the name of the identifier.
func (t Type) String() string {
	switch t {
	case TypeCustom:
		return "custom"
	case TypeAscii:
		return "ascii"
	case TypeBigInt:
		return "bigint"
	case TypeBlob:
		return "blob"
	case TypeBoolean:
		return "boolean"
	case TypeCounter:
		return "counter"
	case TypeDecimal:
		return "decimal"
	case TypeDouble:
		return "double"
	case TypeFloat:
		return "float"
	case TypeInt:
		return "int"
	case TypeText:
		return "text"
	case TypeTimestamp:
		return "timestamp"
	case TypeUUID:
		return "uuid"
	case TypeVarchar:
		return "varchar"
	case TypeTimeUUID:
		return "timeuuid"
	case TypeInet:
		return "inet"
	case TypeDate:
		return "date"
	//case TypeDuration:
	//	return "duration"
	case TypeTime:
		return "time"
	case TypeSmallInt:
		return "smallint"
	case TypeTinyInt:
		return "tinyint"
	//case TypeList:
	//	return "list"
	//case TypeMap:
	//	return "map"
	case TypeSet:
		return "set"
	case TypeVarint:
		return "varint"
	case TypeJsonB:
		return "jsonb"
	case TypeJsonBMap:
		return "jsonbMap"
	case TypeNullString:
		return "nullString"
	case TypeNullBool:
		return "nullBool"
	case TypeNullInt:
		return "nullInt"
	case TypeNullFloat:
		return "nullFloat"
	case TypeDateIsNull:
		return "dateIsNull"
	default:
		return fmt.Sprintf("unknown_type_%d", t)
	}
}

func getBaseType(name string) Type {
	switch name {
	case "ascii":
		return TypeAscii
	case "bigint":
		return TypeBigInt
	case "blob":
		return TypeBlob
	case "boolean":
		return TypeBoolean
	case "counter":
		return TypeCounter
	case "decimal":
		return TypeDecimal
	case "double":
		return TypeDouble
	case "float":
		return TypeFloat
	case "int":
		return TypeInt
	case "tinyint":
		return TypeTinyInt
	case "time":
		return TypeTime
	case "date":
		return TypeDate
	case "dateIsNull":
		return TypeDateIsNull
	case "timestamp":
		return TypeTimestamp
	case "uuid":
		return TypeUUID
	case "varchar":
		return TypeVarchar
	case "text":
		return TypeText
	case "varint":
		return TypeVarint
	case "timeuuid":
		return TypeTimeUUID
	case "inet":
		return TypeInet
	//case "MapType":
	//	return TypeMap
	//case "ListType":
	//	return TypeList
	case "SetType":
		return TypeSet
	case "jsonb":
		return TypeJsonB
	case "jsonbMap":
		return TypeJsonBMap
	case "nullString":
		return TypeNullString
	case "nullBool":
		return TypeNullBool
	case "nullInt":
		return TypeNullInt
	case "nullFloat":
		return TypeNullFloat
	default:
		return TypeCustom
	}
}

func goType(t Type) reflect.Type {
	switch t {
	case TypeVarchar, TypeAscii, TypeInet, TypeText:
		return reflect.TypeOf(*new(string))
	case TypeBigInt, TypeCounter:
		return reflect.TypeOf(*new(int64))
	case TypeTime:
		return reflect.TypeOf(*new(time.Duration))
	case TypeTimestamp:
		return reflect.TypeOf(*new(time.Time))
	case TypeBlob:
		return reflect.TypeOf(*new([]byte))
	case TypeBoolean:
		return reflect.TypeOf(*new(bool))
	case TypeFloat:
		return reflect.TypeOf(*new(float32))
	case TypeDouble:
		return reflect.TypeOf(*new(float64))
	case TypeInt:
		return reflect.TypeOf(*new(int))
	case TypeSmallInt:
		return reflect.TypeOf(*new(int16))
	case TypeTinyInt:
		return reflect.TypeOf(*new(int8))
	//case TypeDecimal:
	//	return reflect.TypeOf(*new(*inf.Dec))
	case TypeUUID, TypeTimeUUID:
		return reflect.TypeOf(*new(uuid.UUID))
	//case TypeList, TypeSet:
	//	return reflect.SliceOf(goType(t.(CollectionType).Elem))
	//case TypeMap:
	//	return reflect.MapOf(goType(t.(CollectionType).Key), goType(t.(CollectionType).Elem))
	case TypeVarint:
		return reflect.TypeOf(*new(*big.Int))
	case TypeUDT:
		return reflect.TypeOf(make(map[string]interface{}))
	case TypeDate:
		return reflect.TypeOf(*new(time.Time))
	case TypeDateIsNull:
		return reflect.TypeOf(*new(database.NullTime))
	case TypeJsonB:
		return reflect.TypeOf(*new(postgres.Jsonb))
	case TypeNullString:
		return reflect.TypeOf(*new(sql.NullString))
	case TypeNullBool:
		return reflect.TypeOf(*new(sql.NullBool))
	case TypeNullInt:
		return reflect.TypeOf(*new(sql.NullInt64))
	case TypeNullFloat:
		return reflect.TypeOf(*new(sql.NullFloat64))
	case TypeJsonBMap:
		return reflect.TypeOf(*new(database.JsonBMap))
	default:
		return nil
	}
}
