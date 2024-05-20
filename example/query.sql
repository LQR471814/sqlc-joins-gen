-- name: GetUserData :many
select * from User
where email = ?
inner join
    PSUserCourse on User.email = PSUserCourse.userEmail
inner join
    MoodleUserCourse on User.email = MoodleUserCourse.userEmail;
