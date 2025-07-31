package tests

import (
	"testing"

	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	sn "github.com/Araks1255/mangacage_protos/gen/site_notifications"
)

func TestNotifyAboutNewModerationRequest(t *testing.T) {
	for i := 1; i <= 7; i++ {
		if _, err := env.SiteNotificationsClient.NotifyAboutNewModerationRequest(
			env.Ctx, &sn.ModerationRequest{
				EntityOnModeration: enums.EntityOnModeration(i),
				ID:                 0,
			},
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNotifyAboutNewRole(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := env.SiteNotificationsClient.NotifyUserAboutNewRole(
		env.Ctx, &sn.NewRole{
			UserID: uint64(userID),
			RoleID: 1,
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyAboutNewUserOnVerification(t *testing.T) {
	if _, err := env.SiteNotificationsClient.NotifyAboutNewUserOnVerification(
		env.Ctx, &sn.UserOnVerification{ID: 0},
	); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyAboutSubmittedTeamJoinRequest(t *testing.T) {
	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{
		Roles:    []string{"team_leader"},
		TgUserID: env.TgUserID,
	})

	if err != nil {
		t.Fatal(err)
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, teamLeaderID)
	if err != nil {
		t.Fatal(err)
	}

	if err := testhelpers.AddUserToTeam(env.DB, teamLeaderID, teamID); err != nil {
		t.Fatal(err)
	}

	if _, err := env.SiteNotificationsClient.NotifyAboutSubmittedTeamJoinRequest(
		env.Ctx, &sn.TeamJoinRequest{
			TeamID:      uint64(teamID),
			CandidateID: uint64(candidateID),
			Message:     "message",
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyAboutTeamJoinRequestResponse(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := env.SiteNotificationsClient.NotifyAboutTeamJoinRequestResponse(
		env.Ctx, &sn.TeamJoinRequestResponse{
			TeamID: uint64(teamID),
			UserID: uint64(userID),
			Result: enums.ResultOfTeamJoinRequest_RESULT_OF_TEAM_JOIN_REQUEST_APPROVED,
		},
	); err != nil {
		t.Fatal(err)
	}

	if _, err := env.SiteNotificationsClient.NotifyAboutTeamJoinRequestResponse(
		env.Ctx, &sn.TeamJoinRequestResponse{
			TeamID: uint64(teamID),
			UserID: uint64(userID),
			Result: enums.ResultOfTeamJoinRequest_RESULT_OF_TEAM_JOIN_REQUEST_DECLINED,
			Reason: "reason",
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyAboutTitleTranslateRequest(t *testing.T) {
	if _, err := env.SiteNotificationsClient.NotifyAboutTitleTranslateRequest(
		env.Ctx, &sn.TitleTranslateRequest{
			TitleID:  0,
			SenderID: 0,
			Message:  "",
		},
	); err != nil {
		t.Fatal(err)
	}
}
