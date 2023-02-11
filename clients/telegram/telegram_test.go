package telegram

import (
	"encoding/json"
	"testing"
)

func TestUserStruct(t *testing.T) {
	testCases := []struct {
		name    string
		data    string
		want    User
		wantErr bool
	}{
		{
			name: "test_user1",
			data: `{"id":23,"first_name":"User 1"}`,
			want: User{
				ID:        23,
				FirstName: "User 1",
			},
			wantErr: false,
		},
		{
			name: "test_user2",
			data: `{"id":24,"first_name":"User 2"}`,
			want: User{
				ID:        24,
				FirstName: "User 2",
			},
			wantErr: false,
		},
		{
			name: "test_user3",
			data: `{"id":25,"first_name":"User 3`,
			want: User{
				ID:        0,
				FirstName: "",
			},
			wantErr: true,
		},
		{
			name: "test_user4",
			data: `{"id":"","first_name":"User 4"}`,
			want: User{
				ID:        0,
				FirstName: "User 4",
			},
			wantErr: true,
		},
		{
			name: "test_user5",
			data: `{"id":55,"first_name":55}`,
			want: User{
				ID:        55,
				FirstName: "",
			},
			wantErr: true,
		},
	}

	for _, tst := range testCases {
		t.Run(tst.name, func(t *testing.T) {
			var u User
			if err := json.Unmarshal([]byte(tst.data), &u); err != nil && !tst.wantErr {
				t.Fatalf("failed to unarshal data: %s", err)
			}
			if u.FirstName != tst.want.FirstName || u.ID != tst.want.ID {
				t.Fatal("incorrect result")
			}
		})
	}
}
