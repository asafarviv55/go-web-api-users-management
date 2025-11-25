package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Helper function to generate unique IDs
func generateID(prefix string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return prefix + "-" + hex.EncodeToString(bytes)
}

// User Handlers
func createUserHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user.ID = generateID("user")
	user.IsActive = true

	if err := createUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create audit log
	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       user.ID,
		Action:       "user.created",
		ResourceID:   user.ID,
		ResourceType: "user",
		Status:       "success",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusCreated, user)
}

func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, err := getUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getAllUsersHandler(c *gin.Context) {
	users := getAllUsers()
	c.JSON(http.StatusOK, users)
}

func updateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user.ID = id
	if err := updateUser(id, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       id,
		Action:       "user.updated",
		ResourceID:   id,
		ResourceType: "user",
		Status:       "success",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusOK, user)
}

// Role Handlers
func createRoleHandler(c *gin.Context) {
	var role Role
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	role.ID = generateID("role")
	if err := createRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func getAllRolesHandler(c *gin.Context) {
	roles := getAllRoles()
	c.JSON(http.StatusOK, roles)
}

func getRoleHandler(c *gin.Context) {
	id := c.Param("id")
	role, err := getRole(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Profile Handlers
func createProfileHandler(c *gin.Context) {
	var profile UserProfile
	if err := c.BindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	profile.ID = generateID("profile")
	if err := createProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func getProfileHandler(c *gin.Context) {
	userID := c.Param("userId")
	profile, err := getProfileByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func updateProfileHandler(c *gin.Context) {
	userID := c.Param("userId")
	var profile UserProfile
	if err := c.BindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	profile.UserID = userID
	if err := updateProfile(profile.ID, &profile); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Team Handlers
func createTeamHandler(c *gin.Context) {
	var team Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	team.ID = generateID("team")
	team.MemberCount = 0

	if err := createTeam(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       team.OwnerID,
		Action:       "team.created",
		ResourceID:   team.ID,
		ResourceType: "team",
		Status:       "success",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusCreated, team)
}

func getTeamHandler(c *gin.Context) {
	id := c.Param("id")
	team, err := getTeam(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func getAllTeamsHandler(c *gin.Context) {
	teams := getAllTeams()
	c.JSON(http.StatusOK, teams)
}

func addTeamMemberHandler(c *gin.Context) {
	teamID := c.Param("id")
	var member TeamMember
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	member.ID = generateID("member")
	member.TeamID = teamID

	if err := addTeamMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func getTeamMembersHandler(c *gin.Context) {
	teamID := c.Param("id")
	members := getTeamMembers(teamID)
	c.JSON(http.StatusOK, members)
}

// Audit Log Handlers
func getAuditLogsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	logs := getAuditLogs(limit)
	c.JSON(http.StatusOK, logs)
}

// Password Reset Handlers
func requestPasswordResetHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate reset token
	token := generateID("reset")
	reset := &PasswordReset{
		ID:        generateID("pwreset"),
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	// In a real app, find user by email and set UserID
	// For demo purposes, we'll accept any email
	if err := createPasswordReset(reset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
		"token":   token, // In production, this would be sent via email
	})
}

func resetPasswordHandler(c *gin.Context) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	reset, err := getPasswordResetByToken(request.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid reset token"})
		return
	}

	if reset.Used {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token already used"})
		return
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	// In a real app, update user password here
	markPasswordResetAsUsed(request.Token)

	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       reset.UserID,
		Action:       "password.reset",
		ResourceID:   reset.UserID,
		ResourceType: "user",
		Status:       "success",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

// Session Handlers
func createSessionHandler(c *gin.Context) {
	var request struct {
		UserID string `json:"user_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token := generateID("session")
	session := &Session{
		ID:        generateID("sess"),
		UserID:    request.UserID,
		Token:     token,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := createSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createActivityLog(&ActivityLog{
		ID:           generateID("activity"),
		UserID:       request.UserID,
		ActivityType: "login",
		Description:  "User logged in",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusCreated, session)
}

func getUserSessionsHandler(c *gin.Context) {
	userID := c.Param("userId")
	sessions := getUserSessions(userID)
	c.JSON(http.StatusOK, sessions)
}

func deleteSessionHandler(c *gin.Context) {
	token := c.Param("token")
	if err := deleteSession(token); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
}

// Preferences Handlers
func createPreferencesHandler(c *gin.Context) {
	var prefs UserPreferences
	if err := c.BindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	prefs.ID = generateID("pref")
	if err := createPreferences(&prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prefs)
}

func getPreferencesHandler(c *gin.Context) {
	userID := c.Param("userId")
	prefs, err := getPreferencesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preferences not found"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

func updatePreferencesHandler(c *gin.Context) {
	userID := c.Param("userId")
	var prefs UserPreferences
	if err := c.BindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	prefs.UserID = userID
	if err := updatePreferences(userID, &prefs); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preferences not found"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// Activity Log Handlers
func createActivityLogHandler(c *gin.Context) {
	var log ActivityLog
	if err := c.BindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.ID = generateID("activity")
	log.IPAddress = c.ClientIP()

	if err := createActivityLog(&log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, log)
}

func getUserActivityLogsHandler(c *gin.Context) {
	userID := c.Param("userId")
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	logs := getUserActivityLogs(userID, limit)
	c.JSON(http.StatusOK, logs)
}

// Invitation Handlers
func createInvitationHandler(c *gin.Context) {
	var invitation Invitation
	if err := c.BindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	invitation.ID = generateID("invitation")
	invitation.Token = generateID("invite")
	invitation.Status = "pending"
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	if err := createInvitation(&invitation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       invitation.InvitedBy,
		Action:       "invitation.created",
		ResourceID:   invitation.ID,
		ResourceType: "invitation",
		Status:       "success",
		IPAddress:    c.ClientIP(),
	})

	c.JSON(http.StatusCreated, invitation)
}

func getInvitationHandler(c *gin.Context) {
	token := c.Param("token")
	invitation, err := getInvitationByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	c.JSON(http.StatusOK, invitation)
}

func acceptInvitationHandler(c *gin.Context) {
	token := c.Param("token")

	invitation, err := getInvitationByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	if invitation.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation already processed"})
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		updateInvitationStatus(token, "expired")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation expired"})
		return
	}

	if err := updateInvitationStatus(token, "accepted"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"})
}

func getPendingInvitationsHandler(c *gin.Context) {
	invitations := getPendingInvitations()
	c.JSON(http.StatusOK, invitations)
}

// Permission Handlers
func getAllPermissionsHandler(c *gin.Context) {
	permissions := getAllPermissions()
	c.JSON(http.StatusOK, permissions)
}

func grantUserPermissionHandler(c *gin.Context) {
	userID := c.Param("userId")
	var request struct {
		PermissionID string `json:"permission_id"`
		GrantedBy    string `json:"granted_by"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userPerm := &UserPermission{
		ID:           generateID("userperm"),
		UserID:       userID,
		PermissionID: request.PermissionID,
		GrantedBy:    request.GrantedBy,
	}

	if err := grantUserPermission(userPerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createAuditLog(&AuditLog{
		ID:           generateID("audit"),
		UserID:       request.GrantedBy,
		Action:       "permission.granted",
		ResourceID:   userID,
		ResourceType: "user",
		Status:       "success",
		IPAddress:    c.ClientIP(),
		Details: map[string]interface{}{
			"permission_id": request.PermissionID,
			"target_user":   userID,
		},
	})

	c.JSON(http.StatusCreated, userPerm)
}

func getUserPermissionsHandler(c *gin.Context) {
	userID := c.Param("userId")
	permissions := getUserPermissions(userID)
	c.JSON(http.StatusOK, permissions)
}

func revokeUserPermissionHandler(c *gin.Context) {
	userID := c.Param("userId")
	permissionID := c.Param("permissionId")

	if err := revokeUserPermission(userID, permissionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked"})
}
