-- name: GetUserData :many
select * from User
where email = ?
inner join PSUserCourse on User.email = PSUserCourse.userEmail
inner join PSUserAssignment on PSUserCourse.courseName = 
inner join MoodleUserCourse on User.email = MoodleUserCourse.userEmail
inner join MoodleCourse on MoodleUserCourse.courseId = MoodleCourse.id;
