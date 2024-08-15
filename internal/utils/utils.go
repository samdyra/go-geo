package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCoordinate converts a slice of two strings to a slice of two float64 values
func ParseCoordinate(coordinateStr []string) ([]float64, error) {
    if len(coordinateStr) != 2 {
        return nil, fmt.Errorf("invalid coordinate format")
    }
    
    lon, err := strconv.ParseFloat(coordinateStr[0], 64)
    if err != nil {
        return nil, fmt.Errorf("invalid longitude: %v", err)
    }
    
    lat, err := strconv.ParseFloat(coordinateStr[1], 64)
    if err != nil {
        return nil, fmt.Errorf("invalid latitude: %v", err)
    }
    
    return []float64{lon, lat}, nil
}

// FormatTableName converts snake_case to Title Case
func FormatTableName(tableName string) string {
    words := strings.Split(tableName, "_")
    for i, word := range words {
        words[i] = strings.Title(word)
    }
    return strings.Join(words, " ")
}

// GetLayerType returns the appropriate layer type based on the data type
func GetLayerType(dataType string) string {
    switch dataType {
    case "LINESTRING":
        return "line"
    case "POLYGON":
        return "fill"
    case "POINT":
        return "circle"
    default:
        return "line" // Default to line if unknown
    }
}

// GetPaint returns the appropriate paint configuration based on the data type and color
func GetPaint(dataType, color string) map[string]interface{} {
    switch dataType {
    case "LINESTRING":
        return map[string]interface{}{
            "line-color":   color,
            "line-width":   5,
            "line-opacity": 0.8,
        }
    case "POLYGON":
        return map[string]interface{}{
            "fill-color":   color,
            "fill-opacity": 0.8,
        }
    case "POINT":
        return map[string]interface{}{
            "circle-radius":        7,
            "circle-color":         color,
            "circle-opacity":       0.8,
            "circle-stroke-width":  1,
        }
    default:
        return map[string]interface{}{
            "line-color":   color,
            "line-width":   5,
            "line-opacity": 0.8,
        }
    }
}