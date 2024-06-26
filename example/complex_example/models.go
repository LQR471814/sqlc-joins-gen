// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package complex_example

import (
	"database/sql"
)

type MoodleAssignment struct {
	Name        string
	Courseid    string
	Description sql.NullString
	Duedate     int64
	Category    sql.NullString
}

type MoodleCourse struct {
	ID         string
	Coursename string
	Teacher    sql.NullString
	Zoom       sql.NullString
}

type MoodlePage struct {
	Courseid string
	Url      string
	Content  string
}

type MoodleUserCourse struct {
	Courseid  string
	Useremail string
}

type PSAssignment struct {
	Name               string
	Coursename         string
	Assignmenttypename string
	Description        sql.NullString
	Duedate            int64
	Category           string
}

type PSAssignmentType struct {
	Coursename string
	Name       string
}

type PSCourse struct {
	Name string
}

type PSUserAssignment struct {
	Useremail      string
	Assignmentname string
	Coursename     string
	Missing        int64
	Collected      int64
	Scored         sql.NullFloat64
	Total          sql.NullFloat64
}

type PSUserCourse struct {
	Useremail  string
	Coursename string
}

type PSUserMeeting struct {
	Useremail  string
	Coursename string
	Starttime  int64
	Endtime    int64
}

type User struct {
	Email string
	Gpa   float64
}

type WeightCourse struct {
	Name string
}

type WeightCourseAssignmentType struct {
	Coursename string
	Name       string
	Weight     float64
}
