package main

import "time"

// User represents the main user entity
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	RoleID    string    `json:"role_id"`
	TeamID    string    `json:"team_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role represents user roles for RBAC
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Avatar      string    `json:"avatar"`
	Bio         string    `json:"bio"`
	PhoneNumber string    `json:"phone_number"`
	Location    string    `json:"location"`
	Company     string    `json:"company"`
	Website     string    `json:"website"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Team represents a group of users
type Team struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamMember represents the relationship between users and teams
type TeamMember struct {
	ID       string    `json:"id"`
	TeamID   string    `json:"team_id"`
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"` // admin, member, viewer
	JoinedAt time.Time `json:"joined_at"`
}

// AuditLog represents system audit trail
type AuditLog struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Action      string                 `json:"action"`
	ResourceID  string                 `json:"resource_id"`
	ResourceType string                `json:"resource_type"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Status      string                 `json:"status"` // success, failure
	Details     map[string]interface{} `json:"details,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// PasswordReset represents password reset tokens
type PasswordReset struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents user sessions
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Token        string    `json:"token"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserPreferences represents user-specific settings
type UserPreferences struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Theme        string                 `json:"theme"` // light, dark, auto
	Language     string                 `json:"language"`
	Timezone     string                 `json:"timezone"`
	Notifications map[string]bool       `json:"notifications"`
	Settings     map[string]interface{} `json:"settings"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ActivityLog represents user activity tracking
type ActivityLog struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	ActivityType string                `json:"activity_type"` // login, logout, view, edit, create, delete
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	IPAddress   string                 `json:"ip_address"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Invitation represents team/system invitations
type Invitation struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	TeamID     string    `json:"team_id,omitempty"`
	RoleID     string    `json:"role_id"`
	InvitedBy  string    `json:"invited_by"`
	Token      string    `json:"token"`
	Status     string    `json:"status"` // pending, accepted, expired, revoked
	ExpiresAt  time.Time `json:"expires_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Permission represents individual permissions
type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"` // create, read, update, delete, manage
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserPermission represents custom permissions assigned to users
type UserPermission struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	PermissionID string    `json:"permission_id"`
	GrantedBy    string    `json:"granted_by"`
	GrantedAt    time.Time `json:"granted_at"`
}
