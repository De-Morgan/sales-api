package usergrp

import (
	"sales-api/business/core/user"
	"sales-api/business/web/v1/response"
)

type userRes struct {
	User AppUser `json:"user"`
}

func userResponse(usr user.User) response.Success[userRes] {
	return response.NewSuccess(userRes{
		User: toAppUser(usr),
	})
}
