package generated

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"ehealth-migration/pkg/structtag"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strconv"
	"time"
)

// TO DO package generated
//const charset = "АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЮЯ"
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const charsetCyrlUkr = "абвгґдеєжзиіїйклмнопрстуфхцчшщюя" +
	"АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЮЯ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GeneratedString(length int, lang string) string {
	var chars string
	chars = charset
	if lang == "ukr" {
		chars = charsetCyrlUkr
	}
	return StringWithCharset(length, chars)
}

func RandNumbersGenerated(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min) + min
}

func GetHex(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

func RandToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func isNull(value interface{}) interface{} {
	switch value.(type) {
	case sql.NullString:
		return sql.NullString{}
	case sql.NullBool:
		return sql.NullBool{}
	case sql.NullInt64:
		return sql.NullInt64{}
	case sql.NullFloat64:
		return sql.NullFloat64{}
	}

	return nil
}

func GetGenerated(objTag interface{}, value interface{}, lang string) interface{} {
	tag := objTag.(*structtag.Tag)
	switch tag.Nmae {
	case "hex":
		switch value.(type) {
		case sql.NullString:
			result := sql.NullString{}
			val := value.(sql.NullString).String
			if val != "" {
				result.String = GetHex(val)
				result.Valid = true
			}
			return result
		case string:
			return GetHex(value.(string))
		}

		return nil
	case "md5Hash":
		switch value.(type) {
		case sql.NullString:
			result := sql.NullString{}
			val := value.(sql.NullString).String
			if val != "" {
				result.String = GetMD5Hash(RandToken())[0:32]
				result.Valid = true
			}
			return result
		case string:
			return GetMD5Hash(RandToken())[0:32]
		}

		return nil
	case "null":
		return isNull(value)
	case "customJsonb":
		customString, err := tag.GetOption("value")

		if err != nil {
			log.Error().Err(err).Msg("GetGenerated error")
			return nil
		}

		jsonData := json.RawMessage(customString.CustomText)

		return postgres.Jsonb{jsonData}
	case "customString":
		customString, err := tag.GetOption("value")

		if err != nil {
			log.Error().Err(err).Msg("GetGenerated error")
			return nil
		}

		return customString.CustomText
	case "randString":
		var randString string
		size, err := tag.GetOption("size")
		if err == nil {
			for _, item := range size.Sizes {
				randString += strconv.Itoa(RandNumbersGenerated(item.Min, item.Max))
			}
		} else {
			randString += strconv.Itoa(RandNumbersGenerated(1000000000, 9999999999))
		}

		return randString
	case "randStringDouble":
		var randStringDouble string
		size, err := tag.GetOption("size")
		if err == nil {
			for key, item := range size.Sizes {
				if key > 0 {
					randStringDouble += "-"
				}
				randStringDouble += strconv.Itoa(RandNumbersGenerated(item.Min, item.Max))
			}
		} else {
			randStringDouble += strconv.Itoa(RandNumbersGenerated(10000000, 99999999)) + "-" + strconv.Itoa(RandNumbersGenerated(10000, 99999))
		}

		return randStringDouble
	case "randIntAndString":
		var randIntAndString string
		sizeInt := 2

		tagSizeInt, err := tag.GetOption("size_int")
		if err == nil {
			sizeInt = tagSizeInt.SizeInt
		}

		randIntAndString = GeneratedString(sizeInt, lang)
		size, err := tag.GetOption("size")
		if err == nil {
			for _, item := range size.Sizes {
				randIntAndString += strconv.Itoa(RandNumbersGenerated(item.Min, item.Max))
			}
		} else {
			randIntAndString += strconv.Itoa(RandNumbersGenerated(100000, 999999))
		}

		return randIntAndString
	case "phoneRandom":
		phoneRandom := "+"
		size, err := tag.GetOption("size")
		if err == nil {
			for _, item := range size.Sizes {
				phoneRandom += strconv.Itoa(RandNumbersGenerated(item.Min, item.Max))
			}
		} else {
			phoneRandom += strconv.Itoa(RandNumbersGenerated(10000000000, 99999999999))
		}

		return phoneRandom
	default:
		return nil
	}
}
