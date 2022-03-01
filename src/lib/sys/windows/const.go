package windows

const (
	WM_KEYDOWN = 0x0100
	WM_CHAR    = 0x0102
)

/**
 * GetWindow() Constants
 */
const (
	GW_HWNDFIRST = 0
	GW_HWNDLAST  = 1
	GW_HWNDNEXT  = 2
	GW_HWNDPREV  = 3
	GW_OWNER     = 4
	GW_CHILD     = 5
)

/*
 * ShowWindow() Commands
 */
const (
	SW_HIDE            = 0
	SW_SHOWNORMAL      = 1
	SW_NORMAL          = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_MAXIMIZE        = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
	SW_MAX             = 11
)
