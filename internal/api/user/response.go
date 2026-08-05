package user

import "learn-gin/internal/infra/database"

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

func toUserResponse(u *database.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   u.Status,
	}
}

func toUserResponseList(list []database.User) []UserResponse {
	res := make([]UserResponse, 0, len(list))
	for i := range list {
		res = append(res, toUserResponse(&list[i]))
	}
	return res
}
