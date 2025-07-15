package group

type createGroupRequest struct {
	Name string `json:"name" validate:"required"`
}

type joinRequest struct {
	Password string `json:"password" validate:"required"`
}

type inviteRequest struct {
	Username string `json:"username" validate:"required"`
}

type setUserTypeRequest struct {
	Type int32 `json:"type" validate:"required"`
}

type changeNameRequest struct {
	Name string `json:"name" validate:"required"`
}

type changePasswordRequest struct {
	Password string `json:"password" validate:"required"`
}
