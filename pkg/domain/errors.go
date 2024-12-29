package domain

import "errors"

var (
	ErrUserNotFound             = errors.New("user doesn't exists")
	ErrCompanyNotFound          = errors.New("company doesn't exists")
	ErrPositionNotFound         = errors.New("position doesn't exists")
	ErrOperatorNotFound         = errors.New("operator doesn't exists")
	ErrSectionNotFound          = errors.New("section doesn't exists")
	ErrRestrictionNotFound      = errors.New("restriction doesn't exists")
	ErrProfileNotFound          = errors.New("profile doesn`t exists")
	ErrRolesNotFound            = errors.New("roles doesn`t exists")
	ErrColleaguesNotFound       = errors.New("colleagues doesn`t exists")
	ErrAvatarNotFound           = errors.New("avatar doesn`t exists")
	ErrUserAvatarNotFound       = errors.New("user_not_found")
	ErrWarehouseNotFound        = errors.New("warehouse doesn`t exists")
	ErrSupplierNotFound         = errors.New("supplier doesn`t exists")
	ErrMaterialNotFound         = errors.New("material doesn`t exists")
	ErrMaterialCategoryNotFound = errors.New("material category doesn`t exists")

	ErrUserAlreadyExists = errors.New("user with such username or email already exists")

	ErrGeneratePassword = errors.New("can`t to generate new password for user")
	ErrGenerateUUID     = errors.New("can`t to generate uuid")
	ErrGenerateJWT      = errors.New("failed to generate JWT")
	ErrGenerateAvatar   = errors.New("failed to generate avatar")

	ErrInvalidInputBody        = errors.New("invalid input body")
	ErrInvalidLimitParam       = errors.New("invalid Limit param")
	ErrInvalidSortParam        = errors.New("invalid Sort param")
	ErrInvalidSortFieldParam   = errors.New("invalid Sort field param")
	ErrInvalidQueryParam       = errors.New("invalid query param")
	ErrInvalidOffsetParam      = errors.New("invalid offset param")
	ErrInvalidIdParam          = errors.New("invalid id param")
	ErrInvalidCountryCodeParam = errors.New("invalid country code param")
	ErrInvalidAdminCodeParam   = errors.New("invalid admin code param")
	ErrInvalidEmailParam       = errors.New("invalid email")
	ErrInvalidAdminToken       = errors.New("invalid admin auth token")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrRefreshTokenNotFound    = errors.New("refresh token not found")
	ErrExpiredRefreshToken     = errors.New("refresh token expired")
	ErrInvalidTimezone         = errors.New("invalid timezone")

	ErrCreateUser    = errors.New("can`t to create new user")
	ErrCreateCompany = errors.New("can`t to create new company")
	ErrCreateRole    = errors.New("can`t to create new role")

	ErrDeleteUser     = errors.New("can`t to delete user")
	ErrDeleteCompany  = errors.New("can`t to delete user")
	ErrDeleteSection  = errors.New("can`t to delete section")
	ErrDeleteChat     = errors.New("can`t to delete chat")
	ErrDeleteOperator = errors.New("can`t to delete operator")
	ErrDeleteRole     = errors.New("can`t to delete role")

	ErrUpdateUser     = errors.New("can`t to update user")
	ErrUpdateCompany  = errors.New("can`t to update company")
	ErrUpdateOperator = errors.New("can`t to update operator")
	ErrUpdateAdmin    = errors.New("can`t to update admin")
	ErrUpdatePassword = errors.New("can`t to update password")
	ErrUpdateCall     = errors.New("can`t to update call")
	ErrUpdateProfile  = errors.New("can`t to update profile")
	ErrUpdateStatus   = errors.New("can`t to update status")

	ErrSmallPAssword = errors.New("password should be at least 6 characters long")
	ErrSmallName     = errors.New("name field cannot be empty")

	ErrMoveAdminToUser = errors.New("can`t move admin user to user")

	ErrSendToNats     = errors.New("error send to nats")
	ErrGetIpAddress   = errors.New("error to get ip address")
	ErrJSONMarshal    = errors.New("the error occurred during data marshaling")
	ErrGRPCLogging    = errors.New("the error occurred in the process of sending logs via grpc")
	ErrGRPCMessage    = errors.New("the error occurred in the process of sending message via grpc")
	ErrGRPCBilling    = errors.New("the error occurred in the process of billing via grpc")
	ErrAuthentication = errors.New("authorization error")
	ErrParseInt       = errors.New("can`t parse int")

	ErrAddSection = errors.New("can`t to add new section")
	ErrAddCompany = errors.New("can`t to add new company")

	ErrCompanyNotApproved = errors.New("company_not_approved")
	ErrCompanyBlocked     = errors.New("company_blocked")
	ErrUserIsNotApproved  = errors.New("user is not approved")
	ErrUserBlocked        = errors.New("user_blocked")
	ErrUserIsNotActive    = errors.New("user is not active")

	ErrNotAllowed         = errors.New("not allowed")
	ErrRoleNotAllowed     = errors.New("role not allowed")
	ErrSectionsNotAllowed = errors.New("sections not allowed")

	ErrSignOut          = errors.New("error occurred during sign out. please try again later or contact support if the problem persists")
	ErrRefreshToken     = errors.New("error occurred during refresh token")
	ErrLoginCredentials = errors.New("invalid login credentials. please check your username and password and try again")

	ErrConvertAvatar = errors.New("failed_to_convert_avatar")
)
