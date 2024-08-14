package models

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	OnlineStatus bool   `json:"online_status" bson:"online_status"`
}

type UserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	OnlineStatus bool   `json:"online_status"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) ConvertToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		OnlineStatus: u.OnlineStatus,
	}
}
