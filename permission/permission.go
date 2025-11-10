package permission

type PermissionCode string

const (
	// Admin
	PermissionAdmin PermissionCode = "admin"
	// Users
	PermissionViewUser   PermissionCode = "view_user"
	PermissionCreateUser PermissionCode = "create_user"
	PermissionUpdateUser PermissionCode = "update_user"
	// Permissions
	PermissionViewPermission   PermissionCode = "view_permission"
	PermissionCreatePermission PermissionCode = "create_permission"
	PermissionUpdatePermission PermissionCode = "update_permission"
	// Roles
	PermissionViewRole   PermissionCode = "view_role"
	PermissionCreateRole PermissionCode = "create_role"
	PermissionUpdateRole PermissionCode = "update_role"
	// Tenants
	PermissionViewTenant   PermissionCode = "view_tenant"
	PermissionCreateTenant PermissionCode = "create_tenant"
	PermissionUpdateTenant PermissionCode = "update_tenant"
)
