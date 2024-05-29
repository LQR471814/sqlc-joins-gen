package outputs

import "testing"

func TestJoinLine(t *testing.T) {
	cases := []struct {
		input    SqlJoinLine
		expected string
	}{
		{
			input: SqlJoinLine{
				Table: "Table",
				On: []SqlJoinOn{
					{
						SourceTable: "Source",
						SourceAttr:  "attr",
						TargetTable: "Table",
						TargetAttr:  "sourceAttr",
					},
				},
			},
			expected: `inner join "Table" on "Source"."attr" = "Table"."sourceAttr"`,
		},
		{
			input: SqlJoinLine{
				Table: "Table",
				On: []SqlJoinOn{
					{
						SourceTable: "Source",
						SourceAttr:  "attr",
						TargetTable: "Table",
						TargetAttr:  "sourceAttr",
					},
				},
				Opts: SqlSelectOpts{
					Limit:  20,
					Offset: 23,
					Where:  "Table.field1 = ?",
					OrderBy: []SqlOrderBy{
						{
							Table:     "Table",
							Attr:      "field2",
							Ascending: true,
						},
						{
							Table:     "Table",
							Attr:      "field3",
							Ascending: false,
						},
					},
				},
			},
			expected: `inner join (select * from "Table" where Table.field1 = ? order by "Table"."field2" asc, "Table"."field3" dsc limit 20 offset 23) as "Table" on "Source"."attr" = "Table"."sourceAttr"`,
		},
	}

	gen := SqliteGenerator{}
	for _, test := range cases {
		res := gen.joinLine(test.input)
		if res != test.expected {
			t.Fatalf(
				"unexpected result:\nexpected: '%s'\ngot: '%s'",
				test.expected, res,
			)
		}
	}
}


