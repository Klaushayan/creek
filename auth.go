// Copyright (c) 2022-present Shayan <klaushayan@gmail.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package brook

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)


type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Users []User

type AuthManager struct {
	Users Users `json:"users"`
	UserListPath string
}

func NewAuthManager(userListPath string) *AuthManager {
	return &AuthManager{
		UserListPath: userListPath,
	}
}

func (am *AuthManager) Exists() bool {
	if _, err := os.Stat(am.UserListPath); os.IsNotExist(err) || !strings.HasSuffix(am.UserListPath, ".json") {
		return false
	}
	return true
}

func (am *AuthManager) Load() error {
	if !am.Exists() {
		err := errors.New("error: auth.json not found")
		return err
	}
	return LoadJSON(am.UserListPath, &am.Users)
}

//TODO: secure
func (am *AuthManager) Verify(username, password string) bool {
	for _, user := range am.Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

// TODO: move to utils later?
func LoadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}


