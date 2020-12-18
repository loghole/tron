package transport

import (
	"reflect"
)

const maxDepth = 32

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
func merge(dst, src map[string]interface{}) map[string]interface{} {
	return mergeMap(dst, src, 0)
}

func mergeMap(dst, src map[string]interface{}, depth int) map[string]interface{} {
	if depth > maxDepth {
		panic("too deep!")
	}

	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapi(srcVal)
			dstMap, dstMapOk := mapi(dstVal)

			if srcMapOk && dstMapOk {
				srcVal = mergeMap(dstMap, srcMap, depth+1)
			}
		}

		dst[key] = srcVal
	}

	return dst
}

func mapi(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}

		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}

		return m, true
	}

	return map[string]interface{}{}, false
}
