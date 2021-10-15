package response

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

/*
	Public Functions
*/

// For error JSON ｒesponse
func SetErrorResponse(ctx *gin.Context, status_code int, err string) {
	ctx.JSON(status_code, gin.H {
		"result" : ResultFail,
		"error" : err,
	})
}

// For multiple JSON structure response
func SetNormalResponse(ctx *gin.Context, status_code int, result_str string, properties ...map[string]interface{}) {
	result := gin.H { "result" : result_str }
	properties = append(properties, result)
	merged_property := mergeMultipleInterface(properties...)
	ctx.JSON(status_code, merged_property)
}

// TODO: For JSON Structure
func SetStructResponse(ctx *gin.Context, status_code int, result_str string, property interface{}){
	result := gin.H { "result" : result_str }
	marshaled := make(map[string]interface{})
	data, _ := json.Marshal(property)
	_ = json.Unmarshal(data, &marshaled)
	merged := mergeMultipleInterface(result, marshaled)
	ctx.JSON(status_code, merged)
}

/*
	Private Functions
*/

// merging multiple interface type
func mergeMultipleInterface(elements ...map[string]interface{}) map[string]interface{} {
	merged := map[string]interface{}{}
	for _, element := range elements {
		merged = mergeInterface(merged, element)
	}
	return merged
}

// merging interface type
func mergeInterface(m1, m2 map[string]interface{}) map[string]interface{} {
    merged := map[string]interface{}{}

    for k, v := range m1 {
        merged[k] = v
    }
    for k, v := range m2 {
        merged[k] = v
    }
    return merged
}