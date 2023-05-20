package query

import "testing"

func TestFindTablesFromQuery(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		query          string
		expectedTables []string
	}{
		"basic query":              {},
		"comma separated query":    {},
		"query with joins":         {},
		"query with table aliases": {},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

		})
	}

}
