package model

// import (
//	"time"
//	// "gorm.io/gorm"
// )

// // TeamService represents the TeamService service.
// type TeamService interface {
//	ListTeamIdsAndUUIDs(*wrapperspb.Empty) (*ListTeamIdUuid, error)
//	ListTeams(*wrapperspb.Empty) (*ListTeamsResponse, error)
//	ListTeamsMetadata(*wrapperspb.Empty) (*ListTeamsMetadataResponse, error)
//	ListTeamsContent(*wrapperspb.Empty) (*ListTeamsContentResponse, error)
//	GetTeam(*GetTeamByUuidRequest) (*Team, error)
//	CreateTeam(*CreateTeamRequest) (*Team, error)
//	UpdateTeam(*UpdateTeamRequest) (*Team, error)
//	DeleteTeam(*GetTeamByUuidRequest) (*Team, error)
//	GetTeamById(*GetTeamByIdRequest) (*Team, error)
//	GetTeamByName(*GetTeamByNameRequest) (*Team, error)
//	GetTeamByEmail(*GetTeamByEmailRequest) (*Team, error)
//	UpdateTeamMetadata(*UpdateTeamMetadataRequest) (*TeamMetadataResponse, error)
//	UpdateTeamContent(*UpdateTeamContentRequest) (*TeamContentResponse, error)
//	GetTeamMetadata(*GetTeamByUuidRequest) (*TeamMetadataResponse, error)
//	GetTeamContent(*GetTeamByUuidRequest) (*TeamContentResponse, error)
// }

// GetTeamByEmailRequest represents the GetTeamByEmailRequest message.
type GetTeamByEmailRequest struct {
	Email string `json:"email"`
}

// GetTeamByIdRequest represents the GetTeamByIdRequest message.
type GetTeamByIdRequest struct {
	ID string `json:"ID"`
}

// GetTeamByNameRequest represents the GetTeamByNameRequest message.
type GetTeamByNameRequest struct {
	Name string `json:"name"`
}

// GetTeamByTeamNameRequest represents the GetTeamByTeamNameRequest message.
type GetTeamByTeamNameRequest struct {
	TeamName string `json:"teamName"`
}

// GetTeamByUuidRequest represents the GetTeamByUuidRequest message.
type GetTeamByUuidRequest struct {
	UUID string `json:"UUID"`
}

// ListTeamIdUuid represents the ListTeamIdUuid message.
type ListTeamIdUuid struct {
	TeamIdUuids []*IdUuid `json:"teamIdUuids"`
}

// ListTeamsResponse represents the ListTeamsResponse message.
type ListTeamsResponse struct {
	Teams []*Team `json:"teams"`
}

// ListTeamsMetadataResponse represents the ListTeamsMetadataResponse message.
type ListTeamsMetadataResponse struct {
	TeamMetadataResponses []*TeamMetadataResponse `json:"teamMetadataResponses"`
}

// ListTeamsContentResponse represents the ListTeamsContentResponse message.
type ListTeamsContentResponse struct {
	TeamContentResponses []*TeamContentResponse `json:"teamContentResponses"`
}

// CreateTeamRequest represents the CreateTeamRequest message.
type CreateTeamRequest struct {
	ID		 string		   `json:"ID"`
	UUID	 string		   `json:"UUID"`
	Metadata *TeamMetadata `json:"metadata"`
	Content	 *TeamContent  `json:"content"`
}

// UpdateTeamRequest represents the UpdateTeamRequest message.
type UpdateTeamRequest struct {
	UUID	 string		   `json:"UUID"`
	Metadata *TeamMetadata `json:"metadata"`
	Content	 *TeamContent  `json:"content"`
}

// UpdateTeamMetadataRequest represents the UpdateTeamMetadataRequest message.
type UpdateTeamMetadataRequest struct {
	UUID	 string		   `json:"UUID"`
	Metadata *TeamMetadata `json:"metadata"`
}

// UpdateTeamContentRequest represents the UpdateTeamContentRequest message.
type UpdateTeamContentRequest struct {
	UUID	string		 `json:"UUID"`
	Content *TeamContent `json:"content"`
}

// TeamMetadataResponse represents the TeamMetadataResponse message.
type TeamMetadataResponse struct {
	ID		 string		   `json:"ID"`
	UUID	 string		   `json:"UUID"`
	Metadata *TeamMetadata `json:"metadata"`
}

// TeamContentResponse represents the TeamContentResponse message.
type TeamContentResponse struct {
	ID		string		 `json:"ID"`
	UUID	string		 `json:"UUID"`
	Content *TeamContent `json:"content"`
}

// Team represents the Team message.
type Team struct {
	ID		 string		   `gorm:"unique" json:"ID"`
	UUID	 string		   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"UUID"`
	Metadata *TeamMetadata `gorm:"embedded;embeddedPrefix:metadata_" json:"metadata"`
	Content	 *TeamContent  `gorm:"embedded;embeddedPrefix:content_" json:"content"`
}

// TeamMetadata represents the TeamMetadata message.
type TeamMetadata struct {
	Name  string	  `json:"name"`
	Dates *CommonDate `gorm:"embedded;embeddedPrefix:dates_" json:"dates"`
}

// TeamContent represents the TeamContent message.
type TeamContent struct {
	Email		 string	   `gorm:"unique" json:"email"`
	Members		 []*IdUuid `gorm:"-" json:"members"`
	ProductOwner *IdUuid   `gorm:"-" json:"productOwner"`
	ScrumMaster	 *IdUuid   `gorm:"-" json:"scrumMaster"`
}
