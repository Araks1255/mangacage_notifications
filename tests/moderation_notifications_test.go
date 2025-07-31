package tests

import (
	"testing"

	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	mn "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
)

func TestNotifyAboutApprovedModerationRequest(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	entitiesIDs := make([]uint, 8)

	if entitiesIDs[int(enums.Entity_ENTITY_TITLE)], err = testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
		t.Fatal(err)
	}
	if entitiesIDs[int(enums.Entity_ENTITY_CHAPTER)], err = testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
		t.Fatal(err)
	}
	if entitiesIDs[int(enums.Entity_ENTITY_AUTHOR)], err = testhelpers.CreateAuthor(env.DB); err != nil {
		t.Fatal(err)
	}
	if entitiesIDs[int(enums.Entity_ENTITY_TEAM)], err = testhelpers.CreateTeam(env.DB, userID); err != nil {
		t.Fatal(err)
	}
	entitiesIDs[int(enums.Entity_ENTITY_PROFILE)] = userID

	genresIDs, err := testhelpers.CreateGenres(env.DB, 1)
	if err != nil {
		t.Fatal(err)
	}
	tagsIDs, err := testhelpers.CreateTags(env.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	entitiesIDs[int(enums.Entity_ENTITY_GENRE)] = genresIDs[0]
	entitiesIDs[int(enums.Entity_ENTITY_TAG)] = tagsIDs[0]

	for i := 1; i < len(entitiesIDs); i++ {
		if _, err := env.ModerationNotificationsClient.NotifyAboutApprovedModerationRequest(
			env.Ctx, &mn.ApprovedEntity{
				Entity:    enums.Entity(i),
				ID:        uint64(entitiesIDs[i]),
				CreatorID: uint64(userID),
			},
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNotifyAboutNewChapterInTitle(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if err := testhelpers.SubscribeToTitle(env.DB, titleID, userID); err != nil {
		t.Fatal(err)
	}

	if _, err := env.ModerationNotificationsClient.NotifyAboutNewChapterInTitle(
		env.Ctx, &mn.Chapter{ID: uint64(chapterID)},
	); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyUserAboutVerificatedAccount(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := env.ModerationNotificationsClient.NotifyUserAboutVerificatedAccount(
		env.Ctx, &mn.VerificatedUser{ID: uint64(userID)},
	); err != nil {
		t.Fatal(err)
	}
}

func TestSendMessageToUser(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	titleOnModerationID, err := moderation.CreateTitleOnModeration(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}
	chapterOnModerationID, err := moderation.CreateChapterOnModerationWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := env.ModerationNotificationsClient.SendMessageToUser(
		env.Ctx, &mn.MessageFromModerator{
			EntityOnModeration:   enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE,
			EntityOnModerationID: uint64(titleOnModerationID),
			ReceiverID:           uint64(userID),
			Text:                 "message about title",
		},
	); err != nil {
		t.Fatal(err)
	}

	if _, err := env.ModerationNotificationsClient.SendMessageToUser(
		env.Ctx, &mn.MessageFromModerator{
			EntityOnModeration:   enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER,
			EntityOnModerationID: uint64(chapterOnModerationID),
			ReceiverID:           uint64(userID),
			Text:                 "message about chapter",
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestSendModerationRequestDeclineReason(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TgUserID: env.TgUserID})
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 7; i++ {
		if _, err := env.ModerationNotificationsClient.SendModerationRequestDeclineReason(
			env.Ctx, &mn.ModerationRequestDeclineReason{
				EntityOnModeration: enums.EntityOnModeration(i),
				EntityName:         "name",
				CreatorID:          uint64(userID),
				Reason:             "reason",
			},
		); err != nil {
			t.Fatal(err)
		}
	}
}
