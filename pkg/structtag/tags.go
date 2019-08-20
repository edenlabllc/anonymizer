package structtag

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errOptionKeyNotSet = errors.New("option key does not exist")
)

type Tag struct {
	Key     string
	Nmae    string
	Options []Option
}

type Size struct {
	Min int
	Max int
}

type Option struct {
	Key        string
	CustomText string
	SizeInt    int
	Sizes      []Size
}

func Parse(tag string) *Tag {
	var Tag Tag

	if tag != "" {
		sTags := strings.Split(tag, ";")

		if len(sTags) > 0 {
			for key, item := range sTags {
				sTag := strings.SplitN(item, ":", 2)
				if key == 0 {
					Tag.Key = sTag[0]
					Tag.Nmae = sTag[1]
				} else {
					var Option Option
					Option.Key = sTag[0]

					switch sTag[0] {
					case "value":
						Option.CustomText = sTag[1]
						break
					case "size_int":
						sizeInt, _ := strconv.Atoi(strings.TrimSpace(sTag[1]))
						Option.SizeInt = sizeInt
						break
					case "size":
						var sSize []Size
						pSize := strings.TrimSpace(sTag[1])
						pSize = strings.Replace(sTag[1], "),", "|", -1)
						pSize = strings.Replace(pSize, "(", "", -1)
						pSize = strings.Replace(pSize, ")", "", -1)

						splitSize := strings.Split(pSize, "|")
						for _, iSize := range splitSize {
							minMax := strings.Split(iSize, ",")
							if len(minMax) > 0 {
								min, _ := strconv.Atoi(strings.TrimSpace(minMax[0]))
								max, _ := strconv.Atoi(strings.TrimSpace(minMax[1]))

								sSize = append(sSize, Size{
									Min: min,
									Max: max,
								})
							}
						}

						Option.Sizes = sSize
						break
					}

					Tag.Options = append(Tag.Options, Option)
				}
			}
		}
	}

	return &Tag
}

func (t *Tag) GetOption(key string) (Option, error) {
	var option Option
	for _, option := range t.Options {
		if option.Key == key {
			return option, nil
		}
	}

	return option, errOptionKeyNotSet
}
