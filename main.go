package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize default data
	initializeData()

	router := gin.Default()

	// User routes
	router.POST("/users", createUserHandler)
	router.GET("/users", getAllUsersHandler)
	router.GET("/users/:id", getUserHandler)
	router.PUT("/users/:id", updateUserHandler)

	// Role routes (RBAC)
	router.POST("/roles", createRoleHandler)
	router.GET("/roles", getAllRolesHandler)
	router.GET("/roles/:id", getRoleHandler)

	// Profile routes
	router.POST("/profiles", createProfileHandler)
	router.GET("/profiles/user/:userId", getProfileHandler)
	router.PUT("/profiles/user/:userId", updateProfileHandler)

	// Team routes
	router.POST("/teams", createTeamHandler)
	router.GET("/teams", getAllTeamsHandler)
	router.GET("/teams/:id", getTeamHandler)
	router.POST("/teams/:id/members", addTeamMemberHandler)
	router.GET("/teams/:id/members", getTeamMembersHandler)

	// Audit log routes
	router.GET("/audit-logs", getAuditLogsHandler)

	// Password reset routes
	router.POST("/password-reset/request", requestPasswordResetHandler)
	router.POST("/password-reset/reset", resetPasswordHandler)

	// Session routes
	router.POST("/sessions", createSessionHandler)
	router.GET("/sessions/user/:userId", getUserSessionsHandler)
	router.DELETE("/sessions/:token", deleteSessionHandler)

	// Preferences routes
	router.POST("/preferences", createPreferencesHandler)
	router.GET("/preferences/user/:userId", getPreferencesHandler)
	router.PUT("/preferences/user/:userId", updatePreferencesHandler)

	// Activity log routes
	router.POST("/activity-logs", createActivityLogHandler)
	router.GET("/activity-logs/user/:userId", getUserActivityLogsHandler)

	// Invitation routes
	router.POST("/invitations", createInvitationHandler)
	router.GET("/invitations/:token", getInvitationHandler)
	router.POST("/invitations/:token/accept", acceptInvitationHandler)
	router.GET("/invitations/pending", getPendingInvitationsHandler)

	// Permission routes
	router.GET("/permissions", getAllPermissionsHandler)
	router.POST("/users/:userId/permissions", grantUserPermissionHandler)
	router.GET("/users/:userId/permissions", getUserPermissionsHandler)
	router.DELETE("/users/:userId/permissions/:permissionId", revokeUserPermissionHandler)

	router.Run("localhost:8080")
}
