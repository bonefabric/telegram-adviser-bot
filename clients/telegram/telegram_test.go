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

func TestMessageStruct(t *testing.T) {
	testCases := []struct {
		name    string
		data    string
		want    Message
		wantErr bool
	}{
		{
			name: "test_message1",
			data: `{"message_id":11}`,
			want: Message{
				ID:   11,
				From: nil,
				Text: nil,
			},
			wantErr: false,
		},
		{
			name: "test_message1",
			data: `{"message_id":"11"}`,
			want: Message{
				ID:   0,
				From: nil,
				Text: nil,
			},
			wantErr: true,
		},
		{
			name: "test_message1",
			data: `{"message_id":24,"from":{"id":55,"first_name":"user542"}}`,
			want: Message{
				ID: 24,
				From: &User{
					ID:        55,
					FirstName: "user542",
				},
				Text: nil,
			},
			wantErr: false,
		},
	}

	for _, tst := range testCases {
		t.Run(tst.name, func(t *testing.T) {
			var m Message
			if err := json.Unmarshal([]byte(tst.data), &m); err != nil && !tst.wantErr {
				t.Fatalf("failed to unarshal data: %s", err)
			}
			if m.ID != tst.want.ID {
				t.Fatal("incorrect result - ID")
			}
			if (m.Text != nil && tst.want.Text == nil) || (m.Text == nil && tst.want.Text != nil) {
				t.Fatal("incorrect result - Text")
			}
			if (m.From != nil && tst.want.From == nil) || (m.From == nil && tst.want.From != nil) {
				t.Fatal("incorrect result - From")
			}
			if m.From != nil && tst.want.From != nil {
				if m.From.ID != tst.want.From.ID {
					t.Fatal("incorrect result - From.ID")
				}
				if m.From.FirstName != tst.want.From.FirstName {
					t.Fatal("incorrect result - From.FirstName")
				}
			}
		})
	}
}
