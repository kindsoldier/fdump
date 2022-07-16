
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdaapi

import (
    "fdump/dscomm/dsdescr"
)

const AddUserMethod string = "addUser"
type AddUserParams struct {
    Login   string      `msgpack:"login"    json:"login"`
    Pass    string      `msgpack:"pass"    json:"pass"`
    State   string      `msgpack:"state"    json:"state"`
}
type AddUserResult struct {
}
func NewAddUserResult() *AddUserResult {
    return &AddUserResult{}
}
func NewAddUserParams() *AddUserParams {
    return &AddUserParams{}
}

const UpdateUserMethod string = "updateUser"
type UpdateUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
}
type UpdateUserResult struct {
}
func NewUpdateUserResult() *UpdateUserResult {
    return &UpdateUserResult{}
}
func NewUpdateUserParams() *UpdateUserParams {
    return &UpdateUserParams{}
}


const CheckUserMethod string = "checkUser"
type CheckUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
}
type CheckUserResult struct {
    Match   bool        `json:"match"`
}

func NewCheckUserResult() *CheckUserResult {
    return &CheckUserResult{}
}
func NewCheckUserParams() *CheckUserParams {
    return &CheckUserParams{}
}

const ListUsersMethod string = "listUsers"
type ListUsersParams struct {
}
type ListUsersResult struct {
    Users  []*dsdescr.User     `json:"users"`
}
func NewListUsersResult() *ListUsersResult {
    return &ListUsersResult{}
}
func NewListUsersParams() *ListUsersParams {
    return &ListUsersParams{}
}


const DeleteUserMethod string = "deleteUser"
type DeleteUserParams struct {
    Login      string           `json:"login"`
}
type DeleteUserResult struct {
}
func NewDeleteUserResult() *DeleteUserResult {
    return &DeleteUserResult{}
}
func NewDeleteUserParams() *DeleteUserParams {
    return &DeleteUserParams{}
}
