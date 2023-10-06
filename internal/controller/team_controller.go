package controller

import (
	// "net/http"

	"net/http"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TeamControllerInterface interface {
	ListTeamIdsAndUUIDs(c *gin.Context)
	ListTeams(c *gin.Context)
	ListTeamsMetadata(c *gin.Context)
	ListTeamsContent(c *gin.Context)
	GetTeam(c *gin.Context)
	UpdateTeam(c *gin.Context)
	DeleteTeam(c *gin.Context)
	GetTeamByID(c *gin.Context)
	GetTeamByName(c *gin.Context)
	GetTeamByEmail(c *gin.Context)
	UpdateTeamMetadata(c *gin.Context)
	UpdateTeamContent(c *gin.Context)
	GetTeamMetadata(c *gin.Context)
	GetTeamContent(c *gin.Context)
}

type TeamController struct {
	// ctx	  context.Context
	logger		*logrus.Logger
	teamService *service.TeamService
}

func NewTeamController(logger *logrus.Logger, teamService *service.TeamService) *TeamController {
	return &TeamController{
		// ctx:	   ctx,
		logger:		 logger,
		teamService: teamService,
	}
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/ids-and-uuids | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/ids-and-uuids -H "Authorization: Bearer <token>" | jq
// @Summary Get Team IDs and UUIDs
// @Description Get Team IDs and UUIDs
// @Tags Team
// @Produce json
// @Success 200 {object} model.ListTeamIdUuid
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/ids-and-uuids [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) ListTeamIdsAndUUIDs(c *gin.Context) {
	idUuids, err := tc.teamService.ListTeamIdsAndUUIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idUuidArray := []*model.IdUuid{}
	for _, idUuid := range idUuids {
		idUuidArray = append(idUuidArray, &model.IdUuid{
			ID:	  idUuid.ID,
			UUID: idUuid.UUID,
		})
	}

	tc.logger.Info(idUuidArray)
	c.JSON(http.StatusOK, idUuidArray)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams -H "Authorization: Bearer <token>" | jq
// @Summary Get Teams
// @Description Get Teams
// @Tags Team
// @Produce json
// @Success 200 {object} []model.Team
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) ListTeams(c *gin.Context) {
	teams, err := tc.teamService.ListTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(teams)
	c.JSON(http.StatusOK, teams)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/metadata | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/metadata -H "Authorization: Bearer <token>" | jq
// @Summary Get Teams Metadata
// @Description Get Teams Metadata
// @Tags Team
// @Produce json
// @Success 200 {object} model.ListTeamsMetadataResponse
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/metadata [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) ListTeamsMetadata(c *gin.Context) {
	teams, err := tc.teamService.ListTeamsMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(teams)
	c.JSON(http.StatusOK, teams)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/content | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/content -H "Authorization: Bearer <token>" | jq
// @Summary Get Teams Content
// @Description Get Teams Content
// @Tags Team
// @Produce json
// @Success 200 {object} model.ListTeamsContentResponse
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/content [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) ListTeamsContent(c *gin.Context) {
	teams, err := tc.teamService.ListTeamsContent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(teams)
	c.JSON(http.StatusOK, teams)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" | jq
// @Summary Get Team by UUID
// @Description Get Team by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Success 200 {object} model.Team
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeam(c *gin.Context) {
	UUID := c.Param("UUID")
	team, err := tc.teamService.GetTeam(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(team)
	c.JSON(http.StatusOK, team)
}

// curl -s -X POST -H 'Content-Type: application/json' http://0.0.0.0:3001/teams -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}, "content": {"email": "john.doe@mail.com", "phone": "0912345678", "lastName": "Doe", "firstName": "John", "password": "P@ssw0rd!234" }}' | jq
// curl -s -X POST -H 'Content-Type: application/json' http://0.0.0.0:3001/teams -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "UUID": "00000000-0000-0000-0000-000000000000", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}, "content": {"email": "jane.doe@mail.com", "phone": "0987654321", "lastName": "Doe", "firstName": "Jane", "password": "P@ssw0rd!234" }}' | jq
// @Summary Create Team
// @Description Create Team
// @Tags Team
// @Accept json
// @Produce json
// @Param body body model.CreateTeamRequest true "New team details"
// @Success 200 {object} string "Team created"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams [post]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) CreateTeam(c *gin.Context) {
	var team model.Team

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tc.logger.Info(team)

	createTeamRequest := model.CreateTeamRequest{
		ID:		  team.ID,
		UUID:	  team.UUID,
		Metadata: team.Metadata,
		Content:  team.Content,
	}

	createdTeam, err := tc.teamService.CreateTeam(createTeamRequest)
	tc.logger.Info(team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdTeam)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 -d '{"ID": "john.doe", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}, "content": {"email": "john.doe@example.com", "password": "P@ssw0rd!234"}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}, "content": {"email": "jane.doe@example.com", "password": "P@ssw0rd!234"}}' | jq
// @Summary Update Team by UUID
// @Description Update Team by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Param body body model.UpdateTeamRequest true "Updated team details"
// @Success 200 {object} string "Team updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID} [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) UpdateTeam(c *gin.Context) {
	UUID := c.Param("UUID")
	var team model.Team

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tc.logger.Info(team)

	updateTeamRequest := model.UpdateTeamRequest{
		UUID:	  team.UUID,
		Metadata: team.Metadata,
		Content:  team.Content,
	}

	updatedTeam, err := tc.teamService.UpdateTeam(UUID, updateTeamRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTeam)
}

// curl -s -X DELETE -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 | jq
// curl -s -X DELETE -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" | jq
// @Summary Delete Team by UUID
// @Description Delete Team by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Success 200 {object} string "Team deleted"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID} [delete]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) DeleteTeam(c *gin.Context) {
	UUID := c.Param("UUID")

	_, err := tc.teamService.DeleteTeam(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted"})
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/id/john.doe | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/id/john.doe -H "Authorization: Bearer <token>" | jq
// @Summary Get Team by ID
// @Description Get Team by ID
// @Tags Team
// @Produce json
// @Param ID path string true "Team ID"
// @Success 200 {object} string "Team found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/id/{ID} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeamByID(c *gin.Context) {
	ID := c.Param("ID")
	team, err := tc.teamService.GetTeamByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(team)
	c.JSON(http.StatusOK, team)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/name/John%20Doe | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/name/John%20Doe -H "Authorization: Bearer <token>" | jq
// @Summary Get Team by Name
// @Description Get Team by Name
// @Tags Team
// @Produce json
// @Param name path string true "Team name"
// @Success 200 {object} string "Team found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/name/{name} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeamByName(c *gin.Context) {
	name := c.Param("name")
	team, err := tc.teamService.GetTeamByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(team)
	c.JSON(http.StatusOK, team)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/email/john.doe@mail.com | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/email/john.doe@mail.com -H "Authorization: Bearer <token>" | jq
// @Summary Get Team by Email
// @Description Get Team by Email
// @Tags Team
// @Produce json
// @Param email path string true "Team email"
// @Success 200 {object} string "Team found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/email/{email} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeamByEmail(c *gin.Context) {
	email := c.Param("email")
	team, err := tc.teamService.GetTeamByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(team)
	c.JSON(http.StatusOK, team)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/metadata -d '{"ID": "john.doe", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/metadata -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}}' | jq
// @Summary Update Team Metadata by UUID
// @Description Update Team Metadata by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Param body body model.UpdateTeamMetadataRequest true "Updated team metadata"
// @Success 200 {object} string "Team metadata updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID}/metadata [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) UpdateTeamMetadata(c *gin.Context) {
	UUID := c.Param("UUID")
	var updateTeamMetadataRequest model.UpdateTeamMetadataRequest

	if err := c.ShouldBindJSON(&updateTeamMetadataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tc.logger.Info(updateTeamMetadataRequest)

	updatedTeamMetadataResponse, err := tc.teamService.UpdateTeamMetadata(UUID, updateTeamMetadataRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTeamMetadataResponse)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/content -d '{"ID": "john.doe", "content": {"email": "john.doe@example.com", "phone": "0912345678", "password": "P@ssw0rd!234"}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/content -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "content": {"email": "jane.doe@example.com", "phone": "0987654321", "password": "P@ssw0rd!234"}}' | jq
// @Summary Update Team Content by UUID
// @Description Update Team Content by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Param body body model.UpdateTeamContentRequest true "Updated team content"
// @Success 200 {object} string "Team content updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID}/content [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) UpdateTeamContent(c *gin.Context) {
	UUID := c.Param("UUID")
	var updateTeamContentRequest model.UpdateTeamContentRequest

	if err := c.ShouldBindJSON(&updateTeamContentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tc.logger.Info(updateTeamContentRequest)

	updatedTeamContent, err := tc.teamService.UpdateTeamContent(UUID, updateTeamContentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTeamContent)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/metadata | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/metadata -H "Authorization: Bearer <token>" | jq
// @Summary Get Team Metadata by UUID
// @Description Get Team Metadata by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Success 200 {object} string "Team metadata found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID}/metadata [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeamMetadata(c *gin.Context) {
	UUID := c.Param("UUID")
	teamMetadata, err := tc.teamService.GetTeamMetadata(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(teamMetadata)
	c.JSON(http.StatusOK, teamMetadata)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/content | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/teams/00000000-0000-0000-0000-000000000000/content -H "Authorization: Bearer <token>" | jq
// @Summary Get Team Content by UUID
// @Description Get Team Content by UUID
// @Tags Team
// @Produce json
// @Param UUID path string true "Team UUID"
// @Success 200 {object} string "Team content found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Team not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /teams/{UUID}/content [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (tc *TeamController) GetTeamContent(c *gin.Context) {
	UUID := c.Param("UUID")
	teamContent, err := tc.teamService.GetTeamContent(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.logger.Info(teamContent)
	c.JSON(http.StatusOK, teamContent)
}
