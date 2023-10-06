package model

// import (
//	"time"
//	// "gorm.io/gorm"
// )

// // UserService represents the UserService service.
// type UserService interface {
//	ListUserIdsAndUUIDs(*wrapperspb.Empty) (*ListUserIdUuid, error)
//	ListUsers(*wrapperspb.Empty) (*ListUsersResponse, error)
//	ListUsersMetadata(*wrapperspb.Empty) (*ListUsersMetadataResponse, error)
//	ListUsersContent(*wrapperspb.Empty) (*ListUsersContentResponse, error)
//	GetUser(*GetUserByUuidRequest) (*User, error)
//	CreateUser(*CreateUserRequest) (*User, error)
//	UpdateUser(*UpdateUserRequest) (*User, error)
//	DeleteUser(*GetUserByUuidRequest) (*User, error)
//	GetUserById(*GetUserByIdRequest) (*User, error)
//	GetUserByName(*GetUserByUsernameRequest) (*User, error)
//	GetUserByEmail(*GetUserByEmailRequest) (*User, error)
//	UpdateUserMetadata(*UpdateUserMetadataRequest) (*UserMetadataResponse, error)
//	UpdateUserContent(*UpdateUserContentRequest) (*UserContentResponse, error)
//	GetUserMetadata(*GetUserByUuidRequest) (*UserMetadataResponse, error)
//	GetUserContent(*GetUserByUuidRequest) (*UserContentResponse, error)
// }

// GetUserByEmailRequest represents the GetUserByEmailRequest message.
type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

// GetUserByIdRequest represents the GetUserByIdRequest message.
type GetUserByIdRequest struct {
	ID string `json:"ID"`
}

// GetUserByNameRequest represents the GetUserByNameRequest message.
type GetUserByNameRequest struct {
	Name string `json:"name"`
}

// GetUserByUsernameRequest represents the GetUserByUsernameRequest message.
type GetUserByUsernameRequest struct {
	Username string `json:"username"`
}

// GetUserByUuidRequest represents the GetUserByUuidRequest message.
type GetUserByUuidRequest struct {
	UUID string `json:"UUID"`
}

// ListUserIdUuid represents the ListUserIdUuid message.
type ListUserIdUuid struct {
	UserIdUuids []*IdUuid `json:"userIdUuids"`
}

// ListUsersResponse represents the ListUsersResponse message.
type ListUsersResponse struct {
	Users []*User `json:"users"`
}

// ListUsersMetadataResponse represents the ListUsersMetadataResponse message.
type ListUsersMetadataResponse struct {
	UserMetadataResponses []*UserMetadataResponse `json:"userMetadataResponses"`
}

// ListUsersContentResponse represents the ListUsersContentResponse message.
type ListUsersContentResponse struct {
	UserContentResponses []*UserContentResponse `json:"userContentResponses"`
}

// CreateUserRequest represents the CreateUserRequest message.
type CreateUserRequest struct {
	ID		 string		   `json:"ID"`
	UUID	 string		   `json:"UUID"`
	Metadata *UserMetadata `json:"metadata"`
	Content	 *UserContent  `json:"content"`
}

// UpdateUserRequest represents the UpdateUserRequest message.
type UpdateUserRequest struct {
	UUID	 string		   `json:"UUID"`
	Metadata *UserMetadata `json:"metadata"`
	Content	 *UserContent  `json:"content"`
}

// UpdateUserMetadataRequest represents the UpdateUserMetadataRequest message.
type UpdateUserMetadataRequest struct {
	UUID	 string		   `json:"UUID"`
	Metadata *UserMetadata `json:"metadata"`
}

// UpdateUserContentRequest represents the UpdateUserContentRequest message.
type UpdateUserContentRequest struct {
	UUID	string		 `json:"UUID"`
	Content *UserContent `json:"content"`
}

// UserMetadataResponse represents the UserMetadataResponse message.
type UserMetadataResponse struct {
	ID		 string		   `json:"ID"`
	UUID	 string		   `json:"UUID"`
	Metadata *UserMetadata `json:"metadata"`
}

// UserContentResponse represents the UserContentResponse message.
type UserContentResponse struct {
	ID		string		 `json:"ID"`
	UUID	string		 `json:"UUID"`
	Content *UserContent `json:"content"`
}

// User represents the User message.
type User struct {
	ID		 string		   `gorm:"unique" json:"ID"`
	UUID	 string		   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"UUID"`
	Metadata *UserMetadata `gorm:"embedded;embeddedPrefix:metadata_" json:"metadata"`
	Content	 *UserContent  `gorm:"embedded;embeddedPrefix:content_" json:"content"`
}

// UserMetadata represents the UserMetadata message.
type UserMetadata struct {
	Name  string	  `json:"name"`
	Dates *CommonDate `gorm:"embedded;embeddedPrefix:dates_" json:"dates"`
}

// UserContent represents the UserContent message.
type UserContent struct {
	Email		 string				  `gorm:"unique" json:"email"`
	Phone		 string				  `gorm:"unique" json:"phone"`
	FirstName	 string				  `gorm:"column:first_name" json:"firstName"`
	LastName	 string				  `gorm:"column:last_name" json:"lastName"`
	ProjectRoles []PROJECT_ROLE_TYPES `gorm:"-" json:"projectRoles"`
	ScrumRoles	 []SCRUM_ROLE_TYPES	  `gorm:"-" json:"scrumRoles"`
	Password	 string				  `gorm:"-" json:"-"`
}
