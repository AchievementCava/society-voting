package httpcore

import (
	"errors"
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/events"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiListElections(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	elections, err := database.GetAllElections()
	if err != nil {
		return fmt.Errorf("apiListElections get all elections: %w", err)
	}

	for _, election := range elections {
		if err := election.PopulateCandidates(); err != nil {
			return fmt.Errorf("apiListElections: %w", err)
		}
	}

	return ctx.JSON(elections)
}

func (endpoints) apiElectionsSSE(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	ctx.Set("Content-Type", "text/event-stream")
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		events.AsStreamWriter(events.NewReceiver(events.TopicElectionStarted)),
	)

	return nil
}

func (endpoints) apiGetActiveElectionInformation(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	election, err := database.GetActiveElection()
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	ballot, err := database.GetAllBallotEntriesForElection(election.ID)
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation get ballot: %w", err)
	}

	var response = struct {
		Election *database.Election      `json:"election"`
		Ballot   []*database.BallotEntry `json:"ballot"`
	}{
		Election: election,
		Ballot:   ballot,
	}

	return ctx.JSON(response)
}

func (endpoints) apiVote(ctx *fiber.Ctx) error {
	user, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Vote []int  `json:"vote" validate:"unique"`
		Code string `json:"code" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	election, err := database.GetActiveElection()
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election that you can vote in.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	hasVotedAlready, err := database.HasUserVotedInElection(user.StudentID, election.ID)
	if err != nil {
		return fmt.Errorf("apiVote check if user %s has already voted: %w", user.StudentID, err)
	}

	if hasVotedAlready {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "You have already voted. Go away :)",
		}
	}

	if err := (&database.Vote{
		ElectionID: election.ID,
		UserID:     user.StudentID,
		Choices:    request.Vote,
	}).Insert(); err != nil {
		return fmt.Errorf("apiVote insert user vote: %w", err)
	}

	events.SendEvent(events.TopicVoteReceived, nil)

	ctx.Status(fiber.StatusNoContent)
	return nil
}
