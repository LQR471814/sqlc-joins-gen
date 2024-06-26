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
	MoodleUserCourse []GetUserData0
	PSUserCourse     []GetUserData1
}

// Table: MoodleUserCourse
type GetUserData0 struct {
	CourseId     string
	UserEmail    string
	MoodleCourse GetUserData00
}

// Table: PSUserCourse
type GetUserData1 struct {
	CourseName       string
	UserEmail        string
	PSUserAssignment []GetUserData10
	PSUserMeeting    []GetUserData11
}

// Table: MoodleCourse
type GetUserData00 struct {
	Id               string
	CourseName       string
	Teacher          sql.NullString
	Zoom             sql.NullString
	MoodlePage       []GetUserData000
	MoodleAssignment []GetUserData001
}

// Table: PSUserAssignment
type GetUserData10 struct {
	UserEmail      string
	AssignmentName string
	CourseName     string
	Missing        int
	Collected      int
	Scored         sql.NullFloat64
	Total          sql.NullFloat64
	PSAssignment   GetUserData100
}

// Table: PSUserMeeting
type GetUserData11 struct {
	UserEmail  string
	CourseName string
	StartTime  int
	EndTime    int
}

// Table: MoodlePage
type GetUserData000 struct {
	CourseId string
	Url      string
	Content  string
}

// Table: MoodleAssignment
type GetUserData001 struct {
	Name        string
	CourseId    string
	Description sql.NullString
	Duedate     int
	Category    sql.NullString
}

// Table: PSAssignment
type GetUserData100 struct {
	Name               string
	CourseName         string
	AssignmentTypeName string
	Description        sql.NullString
	Duedate            int
	Category           string
}

const getUserDataQuery = `select
"User"."gpa" as "User_gpa",
"User"."email" as "User_email",
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
"MoodleAssignment"."category" as "MoodleAssignment_category",
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
"PSUserMeeting"."endTime" as "PSUserMeeting_endTime"
from "User"
inner join (select * from "MoodleUserCourse" where MoodleCourse.teacher = $1) as "MoodleUserCourse" on "MoodleUserCourse"."userEmail" = "User"."email"
inner join (select * from "MoodleCourse" where MoodleCourse.id in ($2)) as "MoodleCourse" on "MoodleUserCourse"."courseId" = "MoodleCourse"."id"
inner join "MoodlePage" on "MoodlePage"."courseId" = "MoodleCourse"."id"
inner join "MoodleAssignment" on "MoodleAssignment"."courseId" = "MoodleCourse"."id"
inner join "PSUserCourse" on "PSUserCourse"."userEmail" = "User"."email"
inner join (select * from "PSUserAssignment" where PSUserAssignment.scored != null and PSUserAssignment.total != null) as "PSUserAssignment" on "PSUserAssignment"."courseName" = "PSUserCourse"."courseName" and "PSUserAssignment"."userEmail" = "PSUserCourse"."userEmail"
inner join "PSAssignment" on "PSUserAssignment"."assignmentName" = "PSAssignment"."name" and "PSUserAssignment"."courseName" = "PSAssignment"."courseName"
inner join "PSUserMeeting" on "PSUserMeeting"."userEmail" = "PSUserCourse"."userEmail" and "PSUserMeeting"."courseName" = "PSUserCourse"."courseName"
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
	getUserData10Map := newQueryMap[GetUserData10]()
	getUserData11Map := newQueryMap[GetUserData11]()
	getUserData000Map := newQueryMap[GetUserData000]()
	getUserData001Map := newQueryMap[GetUserData001]()
	getUserData100Map := newQueryMap[GetUserData100]()

	for rows.Next() {
		var getUserData GetUserData
		var getUserData0 GetUserData0
		var getUserData1 GetUserData1
		var getUserData00 GetUserData00
		var getUserData10 GetUserData10
		var getUserData11 GetUserData11
		var getUserData000 GetUserData000
		var getUserData001 GetUserData001
		var getUserData100 GetUserData100

		err := rows.Scan(
			&getUserData.Gpa,
			&getUserData.Email,
			&getUserData0.CourseId,
			&getUserData0.UserEmail,
			&getUserData00.Id,
			&getUserData00.CourseName,
			&getUserData00.Teacher,
			&getUserData00.Zoom,
			&getUserData000.CourseId,
			&getUserData000.Url,
			&getUserData000.Content,
			&getUserData001.Name,
			&getUserData001.CourseId,
			&getUserData001.Description,
			&getUserData001.Duedate,
			&getUserData001.Category,
			&getUserData1.CourseName,
			&getUserData1.UserEmail,
			&getUserData10.UserEmail,
			&getUserData10.AssignmentName,
			&getUserData10.CourseName,
			&getUserData10.Missing,
			&getUserData10.Collected,
			&getUserData10.Scored,
			&getUserData10.Total,
			&getUserData100.Name,
			&getUserData100.CourseName,
			&getUserData100.AssignmentTypeName,
			&getUserData100.Description,
			&getUserData100.Duedate,
			&getUserData100.Category,
			&getUserData11.UserEmail,
			&getUserData11.CourseName,
			&getUserData11.StartTime,
			&getUserData11.EndTime,
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

		getUserData0Pkey := fmt.Sprint(getUserData0.CourseId, getUserData0.UserEmail)
		getUserData0Existing, ok := getUserData0Map.dict[getUserData0Pkey]
		if !ok {
			getUserData0Map.list = append(getUserData0Map.list, getUserData0)
			getUserData0Map.dict[getUserData0Pkey] = &getUserData0Map.list[len(getUserData0Map.list)-1]
			getUserDataExisting.MoodleUserCourse = append(getUserDataExisting.MoodleUserCourse, *getUserData0Existing)
		}

		getUserData1Pkey := fmt.Sprint(getUserData1.CourseName, getUserData1.UserEmail)
		getUserData1Existing, ok := getUserData1Map.dict[getUserData1Pkey]
		if !ok {
			getUserData1Map.list = append(getUserData1Map.list, getUserData1)
			getUserData1Map.dict[getUserData1Pkey] = &getUserData1Map.list[len(getUserData1Map.list)-1]
			getUserDataExisting.PSUserCourse = append(getUserDataExisting.PSUserCourse, *getUserData1Existing)
		}

		getUserData00Pkey := fmt.Sprint(getUserData00.Id)
		getUserData00Existing, ok := getUserData00Map.dict[getUserData00Pkey]
		if !ok {
			getUserData00Map.list = append(getUserData00Map.list, getUserData00)
			getUserData00Map.dict[getUserData00Pkey] = &getUserData00Map.list[len(getUserData00Map.list)-1]
			getUserData0Existing.MoodleCourse = *getUserData00Existing
		}

		getUserData10Pkey := fmt.Sprint(getUserData10.UserEmail, getUserData10.AssignmentName, getUserData10.CourseName)
		getUserData10Existing, ok := getUserData10Map.dict[getUserData10Pkey]
		if !ok {
			getUserData10Map.list = append(getUserData10Map.list, getUserData10)
			getUserData10Map.dict[getUserData10Pkey] = &getUserData10Map.list[len(getUserData10Map.list)-1]
			getUserData1Existing.PSUserAssignment = append(getUserData1Existing.PSUserAssignment, *getUserData10Existing)
		}

		getUserData11Pkey := fmt.Sprint(getUserData11.UserEmail, getUserData11.CourseName, getUserData11.StartTime)
		getUserData11Existing, ok := getUserData11Map.dict[getUserData11Pkey]
		if !ok {
			getUserData11Map.list = append(getUserData11Map.list, getUserData11)
			getUserData11Map.dict[getUserData11Pkey] = &getUserData11Map.list[len(getUserData11Map.list)-1]
			getUserData1Existing.PSUserMeeting = append(getUserData1Existing.PSUserMeeting, *getUserData11Existing)
		}

		getUserData000Pkey := fmt.Sprint(getUserData000.CourseId, getUserData000.Url)
		getUserData000Existing, ok := getUserData000Map.dict[getUserData000Pkey]
		if !ok {
			getUserData000Map.list = append(getUserData000Map.list, getUserData000)
			getUserData000Map.dict[getUserData000Pkey] = &getUserData000Map.list[len(getUserData000Map.list)-1]
			getUserData00Existing.MoodlePage = append(getUserData00Existing.MoodlePage, *getUserData000Existing)
		}

		getUserData001Pkey := fmt.Sprint(getUserData001.Name, getUserData001.CourseId)
		getUserData001Existing, ok := getUserData001Map.dict[getUserData001Pkey]
		if !ok {
			getUserData001Map.list = append(getUserData001Map.list, getUserData001)
			getUserData001Map.dict[getUserData001Pkey] = &getUserData001Map.list[len(getUserData001Map.list)-1]
			getUserData00Existing.MoodleAssignment = append(getUserData00Existing.MoodleAssignment, *getUserData001Existing)
		}

		getUserData100Pkey := fmt.Sprint(getUserData100.Name, getUserData100.CourseName)
		getUserData100Existing, ok := getUserData100Map.dict[getUserData100Pkey]
		if !ok {
			getUserData100Map.list = append(getUserData100Map.list, getUserData100)
			getUserData100Map.dict[getUserData100Pkey] = &getUserData100Map.list[len(getUserData100Map.list)-1]
			getUserData10Existing.PSAssignment = *getUserData100Existing
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
