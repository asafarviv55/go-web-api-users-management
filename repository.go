package main

import (
	"errors"
	"sync"
	"time"
)

// In-memory storage for demonstration purposes
var (
	users           = make(map[string]*User)
	roles           = make(map[string]*Role)
	profiles        = make(map[string]*UserProfile)
	teams           = make(map[string]*Team)
	teamMembers     = make(map[string][]*TeamMember)
	auditLogs       = []*AuditLog{}
	passwordResets  = make(map[string]*PasswordReset)
	sessions        = make(map[string]*Session)
	preferences     = make(map[string]*UserPreferences)
	activityLogs    = []*ActivityLog{}
	invitations     = make(map[string]*Invitation)
	permissions     = make(map[string]*Permission)
	userPermissions = make(map[string][]*UserPermission)

	mu sync.RWMutex
)

// Initialize with default data
func initializeData() {
	// Default roles
	adminRole := &Role{
		ID:          "role-1",
		Name:        "Admin",
		Description: "Full system access",
		Permissions: []string{"*"},
		CreatedAt:   time.Now(),
	}
	userRole := &Role{
		ID:          "role-2",
		Name:        "User",
		Description: "Standard user access",
		Permissions: []string{"read:own", "write:own"},
		CreatedAt:   time.Now(),
	}
	roles[adminRole.ID] = adminRole
	roles[userRole.ID] = userRole

	// Default permissions
	defaultPerms := []*Permission{
		{ID: "perm-1", Name: "users.create", Resource: "users", Action: "create", Description: "Create new users", CreatedAt: time.Now()},
		{ID: "perm-2", Name: "users.read", Resource: "users", Action: "read", Description: "View users", CreatedAt: time.Now()},
		{ID: "perm-3", Name: "users.update", Resource: "users", Action: "update", Description: "Update users", CreatedAt: time.Now()},
		{ID: "perm-4", Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users", CreatedAt: time.Now()},
		{ID: "perm-5", Name: "teams.manage", Resource: "teams", Action: "manage", Description: "Manage teams", CreatedAt: time.Now()},
	}
	for _, p := range defaultPerms {
		permissions[p.ID] = p
	}
}

// UserRepository methods
func createUser(user *User) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[user.ID]; exists {
		return errors.New("user already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	users[user.ID] = user
	return nil
}

func getUser(id string) (*User, error) {
	mu.RLock()
	defer mu.RUnlock()

	user, exists := users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func getAllUsers() []*User {
	mu.RLock()
	defer mu.RUnlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList
}

func updateUser(id string, updatedUser *User) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[id]; !exists {
		return errors.New("user not found")
	}

	updatedUser.UpdatedAt = time.Now()
	users[id] = updatedUser
	return nil
}

// RoleRepository methods
func getRole(id string) (*Role, error) {
	mu.RLock()
	defer mu.RUnlock()

	role, exists := roles[id]
	if !exists {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func getAllRoles() []*Role {
	mu.RLock()
	defer mu.RUnlock()

	roleList := make([]*Role, 0, len(roles))
	for _, role := range roles {
		roleList = append(roleList, role)
	}
	return roleList
}

func createRole(role *Role) error {
	mu.Lock()
	defer mu.Unlock()

	role.CreatedAt = time.Now()
	roles[role.ID] = role
	return nil
}

// ProfileRepository methods
func createProfile(profile *UserProfile) error {
	mu.Lock()
	defer mu.Unlock()

	profile.UpdatedAt = time.Now()
	profiles[profile.ID] = profile
	return nil
}

func getProfileByUserID(userID string) (*UserProfile, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, profile := range profiles {
		if profile.UserID == userID {
			return profile, nil
		}
	}
	return nil, errors.New("profile not found")
}

func updateProfile(id string, updatedProfile *UserProfile) error {
	mu.Lock()
	defer mu.Unlock()

	updatedProfile.UpdatedAt = time.Now()
	profiles[id] = updatedProfile
	return nil
}

// TeamRepository methods
func createTeam(team *Team) error {
	mu.Lock()
	defer mu.Unlock()

	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()
	teams[team.ID] = team
	return nil
}

func getTeam(id string) (*Team, error) {
	mu.RLock()
	defer mu.RUnlock()

	team, exists := teams[id]
	if !exists {
		return nil, errors.New("team not found")
	}
	return team, nil
}

func getAllTeams() []*Team {
	mu.RLock()
	defer mu.RUnlock()

	teamList := make([]*Team, 0, len(teams))
	for _, team := range teams {
		teamList = append(teamList, team)
	}
	return teamList
}

func addTeamMember(member *TeamMember) error {
	mu.Lock()
	defer mu.Unlock()

	member.JoinedAt = time.Now()
	teamMembers[member.TeamID] = append(teamMembers[member.TeamID], member)

	// Update team member count
	if team, exists := teams[member.TeamID]; exists {
		team.MemberCount++
		team.UpdatedAt = time.Now()
	}

	return nil
}

func getTeamMembers(teamID string) []*TeamMember {
	mu.RLock()
	defer mu.RUnlock()

	return teamMembers[teamID]
}

// AuditLogRepository methods
func createAuditLog(log *AuditLog) error {
	mu.Lock()
	defer mu.Unlock()

	log.CreatedAt = time.Now()
	auditLogs = append(auditLogs, log)
	return nil
}

func getAuditLogs(limit int) []*AuditLog {
	mu.RLock()
	defer mu.RUnlock()

	if limit > len(auditLogs) {
		limit = len(auditLogs)
	}

	// Return most recent logs
	start := len(auditLogs) - limit
	if start < 0 {
		start = 0
	}

	return auditLogs[start:]
}

// PasswordResetRepository methods
func createPasswordReset(reset *PasswordReset) error {
	mu.Lock()
	defer mu.Unlock()

	reset.CreatedAt = time.Now()
	passwordResets[reset.Token] = reset
	return nil
}

func getPasswordResetByToken(token string) (*PasswordReset, error) {
	mu.RLock()
	defer mu.RUnlock()

	reset, exists := passwordResets[token]
	if !exists {
		return nil, errors.New("reset token not found")
	}
	return reset, nil
}

func markPasswordResetAsUsed(token string) error {
	mu.Lock()
	defer mu.Unlock()

	if reset, exists := passwordResets[token]; exists {
		reset.Used = true
		return nil
	}
	return errors.New("reset token not found")
}

// SessionRepository methods
func createSession(session *Session) error {
	mu.Lock()
	defer mu.Unlock()

	session.CreatedAt = time.Now()
	session.LastActivity = time.Now()
	sessions[session.Token] = session
	return nil
}

func getSessionByToken(token string) (*Session, error) {
	mu.RLock()
	defer mu.RUnlock()

	session, exists := sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func getUserSessions(userID string) []*Session {
	mu.RLock()
	defer mu.RUnlock()

	var userSessions []*Session
	for _, session := range sessions {
		if session.UserID == userID {
			userSessions = append(userSessions, session)
		}
	}
	return userSessions
}

func deleteSession(token string) error {
	mu.Lock()
	defer mu.Unlock()

	delete(sessions, token)
	return nil
}

// PreferencesRepository methods
func createPreferences(prefs *UserPreferences) error {
	mu.Lock()
	defer mu.Unlock()

	prefs.UpdatedAt = time.Now()
	preferences[prefs.UserID] = prefs
	return nil
}

func getPreferencesByUserID(userID string) (*UserPreferences, error) {
	mu.RLock()
	defer mu.RUnlock()

	prefs, exists := preferences[userID]
	if !exists {
		return nil, errors.New("preferences not found")
	}
	return prefs, nil
}

func updatePreferences(userID string, updatedPrefs *UserPreferences) error {
	mu.Lock()
	defer mu.Unlock()

	updatedPrefs.UpdatedAt = time.Now()
	preferences[userID] = updatedPrefs
	return nil
}

// ActivityLogRepository methods
func createActivityLog(log *ActivityLog) error {
	mu.Lock()
	defer mu.Unlock()

	log.CreatedAt = time.Now()
	activityLogs = append(activityLogs, log)
	return nil
}

func getUserActivityLogs(userID string, limit int) []*ActivityLog {
	mu.RLock()
	defer mu.RUnlock()

	var userLogs []*ActivityLog
	for i := len(activityLogs) - 1; i >= 0 && len(userLogs) < limit; i-- {
		if activityLogs[i].UserID == userID {
			userLogs = append(userLogs, activityLogs[i])
		}
	}
	return userLogs
}

// InvitationRepository methods
func createInvitation(invitation *Invitation) error {
	mu.Lock()
	defer mu.Unlock()

	invitation.CreatedAt = time.Now()
	invitations[invitation.Token] = invitation
	return nil
}

func getInvitationByToken(token string) (*Invitation, error) {
	mu.RLock()
	defer mu.RUnlock()

	invitation, exists := invitations[token]
	if !exists {
		return nil, errors.New("invitation not found")
	}
	return invitation, nil
}

func updateInvitationStatus(token string, status string) error {
	mu.Lock()
	defer mu.Unlock()

	if invitation, exists := invitations[token]; exists {
		invitation.Status = status
		if status == "accepted" {
			now := time.Now()
			invitation.AcceptedAt = &now
		}
		return nil
	}
	return errors.New("invitation not found")
}

func getPendingInvitations() []*Invitation {
	mu.RLock()
	defer mu.RUnlock()

	var pending []*Invitation
	for _, inv := range invitations {
		if inv.Status == "pending" {
			pending = append(pending, inv)
		}
	}
	return pending
}

// PermissionRepository methods
func getAllPermissions() []*Permission {
	mu.RLock()
	defer mu.RUnlock()

	permList := make([]*Permission, 0, len(permissions))
	for _, perm := range permissions {
		permList = append(permList, perm)
	}
	return permList
}

func grantUserPermission(userPerm *UserPermission) error {
	mu.Lock()
	defer mu.Unlock()

	userPerm.GrantedAt = time.Now()
	userPermissions[userPerm.UserID] = append(userPermissions[userPerm.UserID], userPerm)
	return nil
}

func getUserPermissions(userID string) []*UserPermission {
	mu.RLock()
	defer mu.RUnlock()

	return userPermissions[userID]
}

func revokeUserPermission(userID, permissionID string) error {
	mu.Lock()
	defer mu.Unlock()

	perms := userPermissions[userID]
	for i, perm := range perms {
		if perm.PermissionID == permissionID {
			userPermissions[userID] = append(perms[:i], perms[i+1:]...)
			return nil
		}
	}
	return errors.New("permission not found for user")
}
