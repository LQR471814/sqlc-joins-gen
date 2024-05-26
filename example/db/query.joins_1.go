package db

import (
	"context"
	"fmt"
)

func (q *Queries) GetUserData_other(ctx context.Context, email string) ([]GetUserDataRow, error) {
	rows, err := q.db.QueryContext(ctx, queryGetUserDataRow, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GetUserDataRowMap map[string]GetUserDataRow
	var GetUserDataRow2Map map[string]GetUserDataRow2
	var GetUserDataRow22Map map[string]GetUserDataRow22

	for rows.Next() {
		var User GetUserDataRow
		var PSUserCourse GetUserDataRow2
		var PSUserAssignment GetUserDataRow22
		var PSAssignment GetUserDataRow227

		if err := rows.Scan(
			&User.Gpa,
			&User.Email,
			&PSUserCourse.CourseName,
			&PSUserCourse.UserEmail,
			// ...
		); err != nil {
			return nil, err
		}

		UserPkey := fmt.Sprint(User.Email)
		existingUser, ok := GetUserDataRowMap[UserPkey]
		if !ok {
			GetUserDataRowMap[UserPkey] = User
			existingUser = User
		}

		PSUserCourseMapPkey := fmt.Sprint(PSUserCourse.UserEmail, PSUserCourse.CourseName)
		existingPSUserCourse, ok := GetUserDataRow2Map[PSUserCourseMapPkey]
		if !ok {
			GetUserDataRow2Map[PSUserCourseMapPkey] = PSUserCourse
			*existingUser.PSUserCourse = append(*existingUser.PSUserCourse, PSUserCourse)
		}

		PSUserAssignmentPkey := fmt.Sprint(PSUserAssignment.UserEmail, PSUserAssignment.AssignmentName, PSUserAssignment.CourseName)
		existingPSUserAssignment, ok := GetUserDataRow22Map[PSUserAssignmentPkey]
		if !ok {
			PSUserAssignment.PSAssignment = PSAssignment
			GetUserDataRow22Map[PSUserAssignmentPkey] = PSUserAssignment
			*existingPSUserCourse.PSUserAssignment = append(*existingPSUserCourse.PSUserAssignment, PSUserAssignment)
		}
	}

	var items []GetUserDataRow
	for _, i := range GetUserDataRowMap {
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
