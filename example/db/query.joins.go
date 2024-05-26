// Code generated by sqlc-joins-gen. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type GetUserDataRow struct {
	Gpa              float64
	Email            string
	PSUserCourse     []GetUserDataRow2
	MoodleUserCourse []GetUserDataRow3
}

type GetUserDataRow2 struct {
	CourseName       string
	UserEmail        string
	PSUserAssignment []GetUserDataRow22
	PSUserMeeting    []GetUserDataRow23
}

type GetUserDataRow3 struct {
	CourseId     string
	UserEmail    string
	MoodleCourse GetUserDataRow32
}

type GetUserDataRow22 struct {
	UserEmail      string
	AssignmentName string
	CourseName     string
	Missing        int
	Collected      int
	Scored         sql.NullFloat64
	Total          sql.NullFloat64
	PSAssignment   GetUserDataRow227
}

type GetUserDataRow23 struct {
	UserEmail  string
	CourseName string
	StartTime  int
	EndTime    int
}

type GetUserDataRow32 struct {
	Id               string
	CourseName       string
	Teacher          sql.NullString
	Zoom             sql.NullString
	MoodlePage       []GetUserDataRow324
	MoodleAssignment []GetUserDataRow325
}

type GetUserDataRow227 struct {
	Name               string
	CourseName         string
	AssignmentTypeName string
	Description        sql.NullString
	Duedate            int
	Category           string
}

type GetUserDataRow324 struct {
	CourseId string
	Url      string
	Content  string
}

type GetUserDataRow325 struct {
	Name        string
	CourseId    string
	Description sql.NullString
	Duedate     int
	Category    sql.NullString
}

const queryGetUserDataRow = `select
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
inner join PSUserAssignment on PSUserAssignment.courseName = PSUserCourse.courseName and PSUserAssignment.userEmail = PSUserCourse.userEmail
inner join PSAssignment on PSUserAssignment.assignmentName = PSAssignment.name and PSUserAssignment.courseName = PSAssignment.courseName
inner join PSUserMeeting on PSUserMeeting.userEmail = PSUserCourse.userEmail and PSUserMeeting.courseName = PSUserCourse.courseName
inner join MoodleUserCourse on MoodleUserCourse.userEmail = User.email
inner join MoodleCourse on MoodleUserCourse.courseId = MoodleCourse.id
inner join MoodleAssignment on MoodleAssignment.courseId = MoodleCourse.id
inner join MoodlePage on MoodlePage.courseId = MoodleCourse.id
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

func (q *Queries) GetUserData(ctx context.Context, args any) ([]GetUserDataRow, error) {
	rows, err := q.db.QueryContext(ctx, queryGetUserDataRow, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GetUserDataRowMap map[string]GetUserDataRow
	var GetUserDataRow325Map map[string]GetUserDataRow325
	var GetUserDataRow2Map map[string]GetUserDataRow2
	var GetUserDataRow3Map map[string]GetUserDataRow3
	var GetUserDataRow22Map map[string]GetUserDataRow22
	var GetUserDataRow23Map map[string]GetUserDataRow23
	var GetUserDataRow32Map map[string]GetUserDataRow32
	var GetUserDataRow227Map map[string]GetUserDataRow227
	var GetUserDataRow324Map map[string]GetUserDataRow324

	for rows.Next() {
		var GetUserDataRow GetUserDataRow
		var GetUserDataRow3 GetUserDataRow3
		var GetUserDataRow22 GetUserDataRow22
		var GetUserDataRow23 GetUserDataRow23
		var GetUserDataRow32 GetUserDataRow32
		var GetUserDataRow227 GetUserDataRow227
		var GetUserDataRow324 GetUserDataRow324
		var GetUserDataRow325 GetUserDataRow325
		var GetUserDataRow2 GetUserDataRow2

		err := rows.Scan(
			&GetUserDataRow.Gpa,
			&GetUserDataRow.Email,
			&GetUserDataRow2.CourseName,
			&GetUserDataRow2.UserEmail,
			&GetUserDataRow22.UserEmail,
			&GetUserDataRow22.AssignmentName,
			&GetUserDataRow22.CourseName,
			&GetUserDataRow22.Missing,
			&GetUserDataRow22.Collected,
			&GetUserDataRow22.Scored,
			&GetUserDataRow22.Total,
			&GetUserDataRow227.Name,
			&GetUserDataRow227.CourseName,
			&GetUserDataRow227.AssignmentTypeName,
			&GetUserDataRow227.Description,
			&GetUserDataRow227.Duedate,
			&GetUserDataRow227.Category,
			&GetUserDataRow23.UserEmail,
			&GetUserDataRow23.CourseName,
			&GetUserDataRow23.StartTime,
			&GetUserDataRow23.EndTime,
			&GetUserDataRow3.CourseId,
			&GetUserDataRow3.UserEmail,
			&GetUserDataRow32.Id,
			&GetUserDataRow32.CourseName,
			&GetUserDataRow32.Teacher,
			&GetUserDataRow32.Zoom,
			&GetUserDataRow324.CourseId,
			&GetUserDataRow324.Url,
			&GetUserDataRow324.Content,
			&GetUserDataRow325.Name,
			&GetUserDataRow325.CourseId,
			&GetUserDataRow325.Description,
			&GetUserDataRow325.Duedate,
			&GetUserDataRow325.Category,
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