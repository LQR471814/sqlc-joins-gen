// Code generated by sqlc-joins-gen. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

// Table: User
type GetUserData struct {
	Gpa              float64
	Email            string
	PSUserCourse     []GetUserData0
	MoodleUserCourse []GetUserData0
}

// Table: PSUserCourse
type GetUserData0 struct {
	CourseName       string
	UserEmail        string
	PSUserAssignment []GetUserData1
	PSUserMeeting    []GetUserData1
}

// Table: MoodleUserCourse
type GetUserData1 struct {
	CourseId     string
	UserEmail    string
	MoodleCourse GetUserData00
}

// Table: PSUserAssignment
type GetUserData00 struct {
	UserEmail      string
	AssignmentName string
	CourseName     string
	Missing        int
	Collected      int
	Scored         sql.NullFloat64
	Total          sql.NullFloat64
	PSAssignment   GetUserData01
}

// Table: PSUserMeeting
type GetUserData01 struct {
	UserEmail  string
	CourseName string
	StartTime  int
	EndTime    int
}

// Table: MoodleCourse
type GetUserData10 struct {
	Id               string
	CourseName       string
	Teacher          sql.NullString
	Zoom             sql.NullString
	MoodlePage       []GetUserData000
	MoodleAssignment []GetUserData000
}

// Table: PSAssignment
type GetUserData000 struct {
	Name               string
	CourseName         string
	AssignmentTypeName string
	Description        sql.NullString
	Duedate            int
	Category           string
}

// Table: MoodlePage
type GetUserData100 struct {
	CourseId string
	Url      string
	Content  string
}

// Table: MoodleAssignment
type GetUserData101 struct {
	Name        string
	CourseId    string
	Description sql.NullString
	Duedate     int
	Category    sql.NullString
}

const queryGetUserData = `select
User.gpa as User_gpa
User.email as User_email
PSUserCourse.courseName as PSUserCourse_courseName
PSUserCourse.userEmail as PSUserCourse_userEmail
PSUserAssignment.userEmail as PSUserAssignment_userEmail
PSUserAssignment.assignmentName as PSUserAssignment_assignmentName
PSUserAssignment.courseName as PSUserAssignment_courseName
PSUserAssignment.missing as PSUserAssignment_missing
PSUserAssignment.collected as PSUserAssignment_collected
PSUserAssignment.scored as PSUserAssignment_scored
PSUserAssignment.total as PSUserAssignment_total
PSAssignment.name as PSAssignment_name
PSAssignment.courseName as PSAssignment_courseName
PSAssignment.assignmentTypeName as PSAssignment_assignmentTypeName
PSAssignment.description as PSAssignment_description
PSAssignment.duedate as PSAssignment_duedate
PSAssignment.category as PSAssignment_category
PSUserMeeting.userEmail as PSUserMeeting_userEmail
PSUserMeeting.courseName as PSUserMeeting_courseName
PSUserMeeting.startTime as PSUserMeeting_startTime
PSUserMeeting.endTime as PSUserMeeting_endTime
MoodleUserCourse.courseId as MoodleUserCourse_courseId
MoodleUserCourse.userEmail as MoodleUserCourse_userEmail
MoodleCourse.id as MoodleCourse_id
MoodleCourse.courseName as MoodleCourse_courseName
MoodleCourse.teacher as MoodleCourse_teacher
MoodleCourse.zoom as MoodleCourse_zoom
MoodlePage.courseId as MoodlePage_courseId
MoodlePage.url as MoodlePage_url
MoodlePage.content as MoodlePage_content
MoodleAssignment.name as MoodleAssignment_name
MoodleAssignment.courseId as MoodleAssignment_courseId
MoodleAssignment.description as MoodleAssignment_description
MoodleAssignment.duedate as MoodleAssignment_duedate
MoodleAssignment.category as MoodleAssignment_category
from GetUserData
inner join PSUserCourse on PSUserCourse.userEmail = User.email
inner join PSUserMeeting on PSUserMeeting.userEmail = PSUserCourse.userEmail and PSUserMeeting.courseName = PSUserCourse.courseName
inner join PSUserAssignment on PSUserAssignment.courseName = PSUserCourse.courseName and PSUserAssignment.userEmail = PSUserCourse.userEmail
inner join PSAssignment on PSUserAssignment.assignmentName = PSAssignment.name and PSUserAssignment.courseName = PSAssignment.courseName
inner join MoodleUserCourse on MoodleUserCourse.userEmail = User.email
inner join MoodleCourse on MoodleUserCourse.courseId = MoodleCourse.id
inner join MoodlePage on MoodlePage.courseId = MoodleCourse.id
inner join MoodleAssignment on MoodleAssignment.courseId = MoodleCourse.id
where User.email = ?
order by
PSAssignment.name asc,
PSAssignment.courseName asc,
PSUserAssignment.userEmail asc,
PSUserAssignment.assignmentName asc,
PSUserAssignment.courseName asc,
PSUserMeeting.userEmail asc,
PSUserMeeting.courseName asc,
PSUserMeeting.startTime asc,
PSUserCourse.userEmail asc,
PSUserCourse.courseName asc,
MoodlePage.url asc,
MoodlePage.courseId asc,
MoodleAssignment.name asc,
MoodleAssignment.courseId asc,
MoodleCourse.id asc,
MoodleUserCourse.courseId asc,
MoodleUserCourse.userEmail asc,
User.email asc,
User.gpa asc`

func (q *Queries) GetUserData(ctx context.Context, args any) ([]GetUserData, error) {
	rows, err := q.db.QueryContext(ctx, queryGetUserData, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GetUserDataMap map[string]GetUserData
	var GetUserData0Map map[string]GetUserData0
	var GetUserData1Map map[string]GetUserData1
	var GetUserData00Map map[string]GetUserData00
	var GetUserData01Map map[string]GetUserData01
	var GetUserData10Map map[string]GetUserData10
	var GetUserData000Map map[string]GetUserData000
	var GetUserData100Map map[string]GetUserData100
	var GetUserData101Map map[string]GetUserData101

	for rows.Next() {
		var GetUserData GetUserData
		var GetUserData0 GetUserData0
		var GetUserData1 GetUserData1
		var GetUserData00 GetUserData00
		var GetUserData01 GetUserData01
		var GetUserData10 GetUserData10
		var GetUserData000 GetUserData000
		var GetUserData100 GetUserData100
		var GetUserData101 GetUserData101

		err := rows.Scan(
			&GetUserData.Gpa,
			&GetUserData.Email,
			&GetUserData0.CourseName,
			&GetUserData0.UserEmail,
			&GetUserData1.CourseId,
			&GetUserData1.UserEmail,
			&GetUserData00.UserEmail,
			&GetUserData00.AssignmentName,
			&GetUserData00.CourseName,
			&GetUserData00.Missing,
			&GetUserData00.Collected,
			&GetUserData00.Scored,
			&GetUserData00.Total,
			&GetUserData01.UserEmail,
			&GetUserData01.CourseName,
			&GetUserData01.StartTime,
			&GetUserData01.EndTime,
			&GetUserData1.CourseId,
			&GetUserData1.UserEmail,
			&GetUserData00.UserEmail,
			&GetUserData00.AssignmentName,
			&GetUserData00.CourseName,
			&GetUserData00.Missing,
			&GetUserData00.Collected,
			&GetUserData00.Scored,
			&GetUserData00.Total,
			&GetUserData01.UserEmail,
			&GetUserData01.CourseName,
			&GetUserData01.StartTime,
			&GetUserData01.EndTime,
			&GetUserData0.CourseName,
			&GetUserData0.UserEmail,
			&GetUserData1.CourseId,
			&GetUserData1.UserEmail,
			&GetUserData00.UserEmail,
			&GetUserData00.AssignmentName,
			&GetUserData00.CourseName,
			&GetUserData00.Missing,
			&GetUserData00.Collected,
			&GetUserData00.Scored,
			&GetUserData00.Total,
			&GetUserData01.UserEmail,
			&GetUserData01.CourseName,
			&GetUserData01.StartTime,
			&GetUserData01.EndTime,
			&GetUserData1.CourseId,
			&GetUserData1.UserEmail,
			&GetUserData00.UserEmail,
			&GetUserData00.AssignmentName,
			&GetUserData00.CourseName,
			&GetUserData00.Missing,
			&GetUserData00.Collected,
			&GetUserData00.Scored,
			&GetUserData00.Total,
			&GetUserData01.UserEmail,
			&GetUserData01.CourseName,
			&GetUserData01.StartTime,
			&GetUserData01.EndTime,
		)
		if err != nil {
			return nil, err
		}

	}

	var items []GetUserData
	for _, i := range GetUserDataMap {
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
