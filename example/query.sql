-- name: GetUserData :many
select * from User
join PSUserCourse on User.email = PSUserCourse.userEmail
join PSUserAssignment on
    PSUserCourse.courseName = PSUserAssignment.courseName and
    PSUserCourse.userEmail = PSUserAssignment.userEmail
join PSAssignment on
    PSUserAssignment.assignmentName = PSAssignment.name and
    PSUserAssignment.courseName = PSAssignment.courseName
join PSUserMeeting on
    PSUserCourse.userEmail = PSUserMeeting.userEmail and
    PSUserCourse.courseName = PSUserMeeting.courseName
join MoodleUserCourse on User.email = MoodleUserCourse.userEmail
join MoodleCourse on MoodleUserCourse.courseId = MoodleCourse.id
join MoodlePage on MoodleCourse.id = MoodlePage.courseId
join MoodleAssignment on MoodleCourse.id = MoodleAssignment.courseId
where email = ?;
