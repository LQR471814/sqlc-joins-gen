package querycfg

var TESTING_METHODS = []Method{
	{
		Name:   "getUserData",
		Table:  "User",
		Return: MANY,
		Query: Query{
			Columns: map[string]bool{
				"gpa": true,
			},
			With: map[string]Query{
				"PSUserCourse": {
					Columns: map[string]bool{
						"courseName": true,
					},
					With: map[string]Query{
						"PSUserAssignment": {
							With: map[string]Query{
								"PSAssignment": {},
							},
						},
						"PSUserMeeting": {},
					},
				},
				"MoodleUserCourse": {
					Columns: map[string]bool{
						"courseId":  true,
						"userEmail": false,
					},
					With: map[string]Query{
						"MoodleCourse": {
							With: map[string]Query{
								"MoodlePage":       {},
								"MoodleAssignment": {},
							},
						},
					},
				},
			},
			Where: "User.email = ?",
			OrderBy: map[string]string{
				"gpa": "asc",
			},
		},
	},
}
