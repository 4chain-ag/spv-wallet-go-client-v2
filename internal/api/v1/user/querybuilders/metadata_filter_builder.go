package querybuilders

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Metadata map[string]any

const DefaultMaxDepth = 100

type MetadataFilterBuilder struct {
	MaxDepth int
	Metadata Metadata
}

func (m *MetadataFilterBuilder) Build() (url.Values, error) {
	var sb strings.Builder
	for k, v := range m.Metadata {
		var path strings.Builder
		pref := fmt.Sprintf("metadata[%s]", k)
		path.WriteString(pref)
		if err := m.generateQueryParams(0, &path, v, &sb, pref); err != nil {
			return nil, err
		}
	}

	params := make(url.Values)
	if sb.Len() == 0 {
		return params, nil
	}

	pairs := strings.Split(TrimLastAmpersand(sb.String()), "&")
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		split := strings.Split(pair, "=")
		params.Add(split[0], split[1])
	}

	return params, nil
}

func (m *MetadataFilterBuilder) generateQueryParams(depth int, path *strings.Builder, val any, ans *strings.Builder, pref string) error {
	if depth > m.MaxDepth {
		return fmt.Errorf("%w - max depth: %d", ErrMetadataFilterMaxDepthExceeded, m.MaxDepth)
	}

	if val == nil {
		return nil
	}

	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Map {
		return m.processMapQueryParams(depth+1, val, path, ans, pref)
	}

	if t.Kind() == reflect.Slice {
		return m.processSlicQueryParams(depth+1, val, path, ans, pref)
	}

	path.WriteString(fmt.Sprintf("=%v&", val))
	ans.WriteString(path.String())
	return nil
}

func (m *MetadataFilterBuilder) processMapQueryParams(depth int, val any, path *strings.Builder, ans *strings.Builder, pref string) error {
	rval := reflect.ValueOf(val)
	for _, k := range rval.MapKeys() {
		mpv := rval.MapIndex(k)
		path.WriteString(fmt.Sprintf("[%v]", k.Interface()))
		if err := m.generateQueryParams(depth+1, path, mpv.Interface(), ans, pref); err != nil {
			return err
		}
		// Reset path after processing each map entry (backtracking).
		str := path.String()
		trim := str[:strings.LastIndex(str, "[")]
		path.Reset()
		path.WriteString(trim)
	}

	return nil
}

func (m *MetadataFilterBuilder) processSlicQueryParams(depth int, val any, path *strings.Builder, ans *strings.Builder, pref string) error {
	slice := reflect.ValueOf(val)
	for i := 0; i < slice.Len(); i++ {
		path.WriteString("[]")
		slv := slice.Index(i)
		if err := m.generateQueryParams(depth+1, path, slv.Interface(), ans, pref); err != nil {
			return err
		}
		// Reset path after processing each slice element (backtracking).
		str := path.String()
		trim := str[:strings.LastIndex(str, "[]")]
		path.Reset()
		path.WriteString(trim)
	}

	return nil
}

var ErrMetadataFilterMaxDepthExceeded = errors.New("maximum depth of nesting in metadata map exceeded")

func TrimLastAmpersand(s string) string {
	if len(s) > 0 && s[len(s)-1] == '&' {
		return s[:len(s)-1]
	}

	return s
}
