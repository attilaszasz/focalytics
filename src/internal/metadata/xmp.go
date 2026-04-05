package metadata

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
	"time"
)

type xmpValues struct {
	DateCreated             *time.Time
	CameraModel             string
	LensModel               string
	FocalLengthMM           *float64
	NormalizedFocalLengthMM *float64
	ApertureF               *float64
	ShutterSeconds          *float64
	ISO                     *int
}

func parseXMP(path string) (xmpValues, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return xmpValues{}, err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	values := map[string]string{}
	stack := make([]string, 0, 8)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return xmpValues{}, err
		}

		switch typed := token.(type) {
		case xml.StartElement:
			stack = append(stack, typed.Name.Local)
			for _, attribute := range typed.Attr {
				key := attribute.Name.Local
				value := strings.TrimSpace(attribute.Value)
				if key == "" || value == "" {
					continue
				}
				if _, exists := values[key]; !exists {
					values[key] = value
				}
			}
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			value := strings.TrimSpace(string(typed))
			if value == "" {
				continue
			}
			key := stack[len(stack)-1]
			if _, exists := values[key]; !exists {
				values[key] = value
			}
		}
	}

	result := xmpValues{}
	result.DateCreated = firstParsedTime(
		values["DateCreated"],
		values["CreateDate"],
		values["ModifyDate"],
	)
	result.CameraModel = firstNonEmpty(values["Model"], values["CameraModel"])
	result.LensModel = firstNonEmpty(values["Lens"], values["LensModel"])
	result.FocalLengthMM = parseFloatPointer(firstNonEmpty(values["FocalLength"], values["AuxFocalLength"]))
	result.NormalizedFocalLengthMM = parseFloatPointer(firstNonEmpty(values["FocalLengthIn35mmFormat"], values["FocalLengthIn35mmFilm"]))
	result.ApertureF = parseFloatPointer(firstNonEmpty(values["FNumber"], values["ApertureValue"]))
	result.ShutterSeconds = parseExposurePointer(firstNonEmpty(values["ExposureTime"], values["ShutterSpeedValue"]))
	result.ISO = parseIntPointer(firstNonEmpty(values["ISO"], values["ISOSpeedRatings"]))

	return result, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
