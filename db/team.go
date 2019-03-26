package db

import (
	"context"
	"fmt"
	"time"

	"github.com/A9u/function_junction/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/A9u/function_junction/constant"
)

type Team struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	EventID     primitive.ObjectID `json:"event_id"`
	CreatorID   primitive.ObjectID `json:"creator_id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	ShowcaseUrl string             `json:"showcase_url"`
	Description string             `json:"description"`
}

type TeamInfo struct {
	*Team
	CreatorInfo UserInfo          `json:"created_by"`
	Members     []*TeamMemberInfo `json:"members"`
}

func (s *store) CreateTeam(ctx context.Context, collection *mongo.Collection, team *Team) (createdTeam *TeamInfo, err error) {
	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now
	res, err := collection.InsertOne(ctx, team)
	if err != nil {
		fmt.Println("Error in team creation ", err, team)
		return
	}
	id := res.InsertedID
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&team)
	creatorInfo, _ := FindUserInfoByID(ctx, team.CreatorID)
	teamMember := TeamMember{TeamID: team.ID, Status: constant.Accepted, InviteeID: team.CreatorID, EventID: team.EventID}
	s.CreateTeamMember(ctx, app.GetCollection("team_members"), &teamMember)
	members, err := s.ListTeamMember(ctx, team.ID, team.EventID, app.GetCollection("team_members"), app.GetCollection("users"), app.GetCollection("events"), app.GetCollection("teams"))
	teamInfo := TeamInfo{Team: team, CreatorInfo: creatorInfo, Members: members}
	return &teamInfo, err
}

func (s *store) ListTeams(ctx context.Context, collection *mongo.Collection, eventID primitive.ObjectID) (teamsInfo []*TeamInfo, err error) {
	// findOptions := options.Find()
	fmt.Println(collection)
	cur, err := collection.Find(ctx, bson.D{{"eventid", eventID}})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem Team
		err = cur.Decode(&elem)
		creatorInfo, _ := FindUserInfoByID(ctx, elem.CreatorID)
		members, err := s.ListTeamMember(ctx, elem.ID, eventID, app.GetCollection("team_members"), app.GetCollection("users"), app.GetCollection("events"), app.GetCollection("teams"))
		fmt.Println(err)
		teamInfo := TeamInfo{Team: &elem, CreatorInfo: creatorInfo, Members: members}
		teamsInfo = append(teamsInfo, &teamInfo)
	}
	if err := cur.Err(); err != nil {
	}
	return teamsInfo, err
}

func (s *store) FindTeamByID(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (team *Team, err error) {
	err = collection.FindOne(ctx, bson.D{{"_id", teamID}}).Decode(&team)

	if err != nil {
		fmt.Println("Error in FindTeamByID: ", err)
		return
	}
	return team, err
}

func (s *store) FindTeamByEventIDAndName(ctx context.Context, eventID primitive.ObjectID, name string, collection *mongo.Collection) (team *Team, err error) {

	err = collection.FindOne(ctx, bson.D{{"eventid", eventID}}).Decode(&team)

	if err != nil {
		fmt.Println("Error in FindTeamByEventIDAndName: ", err)
		return
	}
	return
}

func (s *store) UpdateTeam(ctx context.Context, id primitive.ObjectID, team Team) (createdTeam TeamInfo, err error) {
	collection := app.GetCollection("teams")

	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set",
		bson.D{{"name", team.Name},
			{"showcaseurl", team.ShowcaseUrl},
			{"description", team.Description},
			{"updatedat", time.Now()}}}})

	if err != nil {
		fmt.Println("Error in team update ", err, team)
		return
	}

	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&team)

	if err != nil {
		fmt.Println("Error in team update ", err, team)
		return
	}

	creatorInfo, _ := FindUserInfoByID(ctx, team.CreatorID)

	createdTeam = TeamInfo{Team: &team, CreatorInfo: creatorInfo}
	return
}
