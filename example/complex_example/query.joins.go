// Code generated by sqlc-joins-gen. DO NOT EDIT.

package complex_example

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type queryMap[T any] struct {
	dict map[string]*T
	list []T
}

func newQueryMap[T any]() queryMap[T] {
	return queryMap[T]{
		dict: make(map[string]*T),
	}
}

// Table: User
type GetUserData struct {
	Gpa              float64
	Email            string
	PSUserCourse     []GetUserData0
	MoodleUserCourse []GetUserData1
}

// Table: PSUserCourse
type GetUserData0 struct {
	CourseName       string
	UserEmail        string
	PSUserAssignment []GetUserData00
	PSUserMeeting    []GetUserData01
}

// Table: MoodleUserCourse
type GetUserData1 struct {
	CourseId     string
	UserEmail    string
	MoodleCourse GetUserData10
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
	PSAssignment   GetUserData000
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
	MoodlePage       []GetUserData100
	MoodleAssignment []GetUserData101
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

const getUserDataQuery = `select
"User"."gpa" as "User_gpa",
"User"."email" as "User_email",
"PSUserCourse"."courseName" as "PSUserCourse_courseName",
"PSUserCourse"."userEmail" as "PSUserCourse_userEmail",
"PSUserAssignment"."userEmail" as "PSUserAssignment_userEmail",
"PSUserAssignment"."assignmentName" as "PSUserAssignment_assignmentName",
"PSUserAssignment"."courseName" as "PSUserAssignment_courseName",
"PSUserAssignment"."missing" as "PSUserAssignment_missing",
"PSUserAssignment"."collected" as "PSUserAssignment_collected",
"PSUserAssignment"."scored" as "PSUserAssignment_scored",
"PSUserAssignment"."total" as "PSUserAssignment_total",
"PSAssignment"."name" as "PSAssignment_name",
"PSAssignment"."courseName" as "PSAssignment_courseName",
"PSAssignment"."assignmentTypeName" as "PSAssignment_assignmentTypeName",
"PSAssignment"."description" as "PSAssignment_description",
"PSAssignment"."duedate" as "PSAssignment_duedate",
"PSAssignment"."category" as "PSAssignment_category",
"PSUserMeeting"."userEmail" as "PSUserMeeting_userEmail",
"PSUserMeeting"."courseName" as "PSUserMeeting_courseName",
"PSUserMeeting"."startTime" as "PSUserMeeting_startTime",
"PSUserMeeting"."endTime" as "PSUserMeeting_endTime",
"MoodleUserCourse"."courseId" as "MoodleUserCourse_courseId",
"MoodleUserCourse"."userEmail" as "MoodleUserCourse_userEmail",
"MoodleCourse"."id" as "MoodleCourse_id",
"MoodleCourse"."courseName" as "MoodleCourse_courseName",
"MoodleCourse"."teacher" as "MoodleCourse_teacher",
"MoodleCourse"."zoom" as "MoodleCourse_zoom",
"MoodlePage"."courseId" as "MoodlePage_courseId",
"MoodlePage"."url" as "MoodlePage_url",
"MoodlePage"."content" as "MoodlePage_content",
"MoodleAssignment"."name" as "MoodleAssignment_name",
"MoodleAssignment"."courseId" as "MoodleAssignment_courseId",
"MoodleAssignment"."description" as "MoodleAssignment_description",
"MoodleAssignment"."duedate" as "MoodleAssignment_duedate",
"MoodleAssignment"."category" as "MoodleAssignment_category"
from "User"
inner join "PSUserCourse" on "PSUserCourse"."userEmail" = "User"."email"
inner join (select * from "PSUserAssignment" where PSUserAssignment.scored != null and PSUserAssignment.total != null) as "PSUserAssignment" on "PSUserAssignment"."courseName" = "PSUserCourse"."courseName" and "PSUserAssignment"."userEmail" = "PSUserCourse"."userEmail"
inner join "PSAssignment" on "PSUserAssignment"."assignmentName" = "PSAssignment"."name" and "PSUserAssignment"."courseName" = "PSAssignment"."courseName"
inner join "PSUserMeeting" on "PSUserMeeting"."userEmail" = "PSUserCourse"."userEmail" and "PSUserMeeting"."courseName" = "PSUserCourse"."courseName"
inner join (select * from "MoodleUserCourse" where MoodleCourse.teacher = $1) as "MoodleUserCourse" on "MoodleUserCourse"."userEmail" = "User"."email"
inner join (select * from "MoodleCourse" where MoodleCourse.id in ($2)) as "MoodleCourse" on "MoodleUserCourse"."courseId" = "MoodleCourse"."id"
inner join "MoodlePage" on "MoodlePage"."courseId" = "MoodleCourse"."id"
inner join "MoodleAssignment" on "MoodleAssignment"."courseId" = "MoodleCourse"."id"
where User.email = $0
order by "User"."gpa" asc
`

func (q *Queries) GetUserData(ctx context.Context, userEmail string, teacher sql.NullString, ids []int) (*GetUserData, error) {
	queryStr := getUserDataQuery

	userEmailStr := `"` + userEmail + `"`
	queryStr = strings.Replace(queryStr, "$0", userEmailStr, 1)

	teacherStr := "null"
	if teacher.Valid {
		teacherStr = `"` + teacher.String + `"`
	}
	queryStr = strings.Replace(queryStr, "$1", teacherStr, 1)

	idsJoined := ""
	for i, e := range ids {
		if i > 0 {
			idsJoined += ", "
		}
		idsStr := fmt.Sprint(e)
		idsJoined += idsStr
	}
	queryStr = strings.Replace(queryStr, "$2", idsJoined, 1)

	rows, err := q.db.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	getUserDataMap := newQueryMap[GetUserData]()
	getUserData0Map := newQueryMap[GetUserData0]()
	getUserData1Map := newQueryMap[GetUserData1]()
	getUserData00Map := newQueryMap[GetUserData00]()
	getUserData01Map := newQueryMap[GetUserData01]()
	getUserData10Map := newQueryMap[GetUserData10]()
	getUserData000Map := newQueryMap[GetUserData000]()
	getUserData100Map := newQueryMap[GetUserData100]()
	getUserData101Map := newQueryMap[GetUserData101]()

	for rows.Next() {
		var getUserData GetUserData
		var getUserData0 GetUserData0
		var getUserData1 GetUserData1
		var getUserData00 GetUserData00
		var getUserData01 GetUserData01
		var getUserData10 GetUserData10
		var getUserData000 GetUserData000
		var getUserData100 GetUserData100
		var getUserData101 GetUserData101

		err := rows.Scan(
			&getUserData.Gpa,
			&getUserData.Email,
			&getUserData0.CourseName,
			&getUserData0.UserEmail,
			&getUserData00.UserEmail,
			&getUserData00.AssignmentName,
			&getUserData00.CourseName,
			&getUserData00.Missing,
			&getUserData00.Collected,
			&getUserData00.Scored,
			&getUserData00.Total,
			&getUserData000.Name,
			&getUserData000.CourseName,
			&getUserData000.AssignmentTypeName,
			&getUserData000.Description,
			&getUserData000.Duedate,
			&getUserData000.Category,
			&getUserData01.UserEmail,
			&getUserData01.CourseName,
			&getUserData01.StartTime,
			&getUserData01.EndTime,
			&getUserData1.CourseId,
			&getUserData1.UserEmail,
			&getUserData10.Id,
			&getUserData10.CourseName,
			&getUserData10.Teacher,
			&getUserData10.Zoom,
			&getUserData100.CourseId,
			&getUserData100.Url,
			&getUserData100.Content,
			&getUserData101.Name,
			&getUserData101.CourseId,
			&getUserData101.Description,
			&getUserData101.Duedate,
			&getUserData101.Category,
		)
		if err != nil {
			return nil, err
		}

		getUserDataPkey := fmt.Sprint(getUserData.Email)
		getUserDataExisting, ok := getUserDataMap.dict[getUserDataPkey]
		if !ok {
			getUserDataMap.list = append(getUserDataMap.list, getUserData)
			getUserDataMap.dict[getUserDataPkey] = &getUserDataMap.list[len(getUserDataMap.list)-1]
		}

		getUserData0Pkey := fmt.Sprint(getUserData0.CourseName, getUserData0.UserEmail)
		getUserData0Existing, ok := getUserData0Map.dict[getUserData0Pkey]
		if !ok {
			getUserData0Map.list = append(getUserData0Map.list, getUserData0)
			getUserData0Map.dict[getUserData0Pkey] = &getUserData0Map.list[len(getUserData0Map.list)-1]
			getUserDataExisting.PSUserCourse = append(getUserDataExisting.PSUserCourse, *getUserData0Existing)
		}

		getUserData1Pkey := fmt.Sprint(getUserData1.CourseId, getUserData1.UserEmail)
		getUserData1Existing, ok := getUserData1Map.dict[getUserData1Pkey]
		if !ok {
			getUserData1Map.list = append(getUserData1Map.list, getUserData1)
			getUserData1Map.dict[getUserData1Pkey] = &getUserData1Map.list[len(getUserData1Map.list)-1]
			getUserDataExisting.MoodleUserCourse = append(getUserDataExisting.MoodleUserCourse, *getUserData1Existing)
		}

		getUserData00Pkey := fmt.Sprint(getUserData00.UserEmail, getUserData00.AssignmentName, getUserData00.CourseName)
		getUserData00Existing, ok := getUserData00Map.dict[getUserData00Pkey]
		if !ok {
			getUserData00Map.list = append(getUserData00Map.list, getUserData00)
			getUserData00Map.dict[getUserData00Pkey] = &getUserData00Map.list[len(getUserData00Map.list)-1]
			getUserData0Existing.PSUserAssignment = append(getUserData0Existing.PSUserAssignment, *getUserData00Existing)
		}

		getUserData01Pkey := fmt.Sprint(getUserData01.UserEmail, getUserData01.CourseName, getUserData01.StartTime)
		getUserData01Existing, ok := getUserData01Map.dict[getUserData01Pkey]
		if !ok {
			getUserData01Map.list = append(getUserData01Map.list, getUserData01)
			getUserData01Map.dict[getUserData01Pkey] = &getUserData01Map.list[len(getUserData01Map.list)-1]
			getUserData0Existing.PSUserMeeting = append(getUserData0Existing.PSUserMeeting, *getUserData01Existing)
		}

		getUserData10Pkey := fmt.Sprint(getUserData10.Id)
		getUserData10Existing, ok := getUserData10Map.dict[getUserData10Pkey]
		if !ok {
			getUserData10Map.list = append(getUserData10Map.list, getUserData10)
			getUserData10Map.dict[getUserData10Pkey] = &getUserData10Map.list[len(getUserData10Map.list)-1]
			getUserData1Existing.MoodleCourse = *getUserData10Existing
		}

		getUserData000Pkey := fmt.Sprint(getUserData000.Name, getUserData000.CourseName)
		getUserData000Existing, ok := getUserData000Map.dict[getUserData000Pkey]
		if !ok {
			getUserData000Map.list = append(getUserData000Map.list, getUserData000)
			getUserData000Map.dict[getUserData000Pkey] = &getUserData000Map.list[len(getUserData000Map.list)-1]
			getUserData00Existing.PSAssignment = *getUserData000Existing
		}

		getUserData100Pkey := fmt.Sprint(getUserData100.CourseId, getUserData100.Url)
		getUserData100Existing, ok := getUserData100Map.dict[getUserData100Pkey]
		if !ok {
			getUserData100Map.list = append(getUserData100Map.list, getUserData100)
			getUserData100Map.dict[getUserData100Pkey] = &getUserData100Map.list[len(getUserData100Map.list)-1]
			getUserData10Existing.MoodlePage = append(getUserData10Existing.MoodlePage, *getUserData100Existing)
		}

		getUserData101Pkey := fmt.Sprint(getUserData101.Name, getUserData101.CourseId)
		getUserData101Existing, ok := getUserData101Map.dict[getUserData101Pkey]
		if !ok {
			getUserData101Map.list = append(getUserData101Map.list, getUserData101)
			getUserData101Map.dict[getUserData101Pkey] = &getUserData101Map.list[len(getUserData101Map.list)-1]
			getUserData10Existing.MoodleAssignment = append(getUserData10Existing.MoodleAssignment, *getUserData101Existing)
		}
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(getUserDataMap.list) == 0 {
		return nil, nil
	}
	return &getUserDataMap.list[0], nil
}
