package status

var (
	// -1 = UnknownError
	UnknownError = New(-1, "Unknown Error")

	// 0 = OK
	OK = New(0, "OK")

	// -1 and 0 are reversed, and SHOULD NOT BE MODIFIED BY ANYONE!

	// 1~99999: CLI
	// 100000~199999: Server
	// 200000~299999: Coordinator

	// 1~999: CLI Global
	CLINotStarted = New(1, "CLI not started")
	CLIIsStarting = New(2, "CLI is starting")

	// 1000~1999: gsm server install
	CLIInstallRequiresRWAccess                         = New(1000, "Installing server requires both Read and Write access")
	CLIInstallNotAllowedBothSetSymlinkAndInstall       = New(1001, "Not allowed set both symlink and import flag")
	CLIInstallServerDirectoryAlreadyExist              = New(1002, "Server directory is not empty")
	CLIInstallUnableToCreateSymlink                    = New(1003, "Unable to create symlink")
	CLIInstallSymlinkServerIDFoundRecursiveReferencing = New(1004, "Symlink server ID found recursive referencing")
	CLIInstallExplicitlyDeclareExternalServerDirectory = New(1005, "Must explictly declare the external server directory")
	CLIInstallExternalServerDirectoryNotExist          = New(1006, "Unable to find external server directory")
	CLIInstallExplicitlyDeclareWhichGameToInstall      = New(1007, "Explicitly declare which game to install")
	CLIInstallGameNotSupported                         = New(1008, "Game is not supported")
	CLIInstallUnableToFoundSteamCMD                    = New(1009, "Unable to found SteamCMD")
	CLIInstallErrorFoundWhenExecutingSteamCMD          = New(1010, "Error found when executing SteamCMD")
	CLIInstallUnableToFindSteamCMDPreRequisite         = New(1011, "Unable to find SteamCMD's pre-requisite")
	CLIInstallRootSymlinkServerNotExist                = New(1012, "Root symlink server is not exist.")
	CLIInstallSymlinkServerNotExist                    = New(1013, "This symlink server for symlink is not exist.")
	CLIInstallSymlinkServerDeleted                     = New(1014, "The server for symlink is deleted")
	CLIInstallErrorOnGeneratingDirectory               = New(1015, "Error found when generating directory.")

	// 2000~2999: gsm server attach
	CLIAttachUnableEstablishConnectionToCoordinator = New(2000, "Unable to establish connection to coordinator")

	// 3000~3999: gsm server backup
	CLIBackupRequireRWAccess = New(3000, "Backing up server requires Read and Write access")

	// 4000~4999: gsm server remove
	CLIRemoveRequireRWAccess = New(4000, "Removing server requires Read and Write access")

	// 5000~5999: gsm server restart
	CLIRestartFailedToRestartServer = New(5000, "Failed to restart the server")

	// 6000~6999: gsm server send
	CLISendUnableToSendCommandToServer = New(6000, "Failed to send command to the server")

	// 7000~7999: gsm server start
	CLIStartUnableToStartupServer = New(7000, "Failed to startup the server")

	// 8000~8999: gsm server stop
	CLIStopUnableToStopServer = New(8000, "Failed to stop the server")

	// 9000~9999: gsm server update
	CLIUpdateUnableToUpdateServer = New(9000, "Failed to update the server")

	// 10000~10999: gsm server validate
	CLIValidateUnableToValidateServer = New(10000, "Unable to validate the server")

	// 11000~11999: gsm server search
	CLISearchUnableToSearchServer = New(11000, "Unable to search the server")

	// 100000~199999: Server
	ServerStartingInteractiveProcess             = New(100000, "Starting interactive process")
	ServerStoppingInteractiveProcess             = New(100001, "Stopping interactive process")
	ServerStarting                               = New(100002, "Starting server")
	ServerStopping                               = New(100003, "Stopping server")
	ServerCrashed                                = New(100004, "Server crashed")
	ServerRestartCountingDown                    = New(100005, "Restarting counting down")
	ServerFoundLastAbnormalExitProcess           = New(100006, "Found last abnormal exit process")
	ServerProcessExitedButInteratingProcessIsNot = New(100007, "Process has been exited but its iterating process is not")
	ServerSearchingCoordinator                   = New(100008, "Searching coordinator")
	ServerFoundCoordinatorConnecting             = New(100009, "Found a coordinator, connecting")
	ServerConnectedCoordinatorAndLoggingIn       = New(100010, "Connected to the coordinator, logging in")

	// 200000~299999: Coordinator
	CoordinatorStarting                           = New(200000, "Starting coordinator")
	CoordinatorConnectingToUpperNode              = New(200001, "Connecting to the upper coordinator node")
	CoordinatorConnectedAndRedirecting            = New(200002, "Connected to the upper node, redirecting")
	CoordinatorConnectedToTheSpecificAndLoggingIn = New(200003, "Connected to the specific coordinator, logging in")
	CoordinatorReceivedMessageFromCLIAndResolving = New(200004, "Received message from CLI, resolving")
	CoordinatorUnknownCLIMessage                  = New(200005, "Unknown CLI message")
	CoordinatorIncorrectStartupArgs               = New(200006, "Incorrect startup args")
	CoordinatorServerNotFound                     = New(200007, "Server not found")
	CoordinatorCoordinatorNotFound                = New(200008, "Coordinator not found")
	CoordinatorServerOffline                      = New(200009, "Server offline")
	CoordinatorCoordinatorOffline                 = New(200010, "Coordinator offline")
	CoordinatorUnknownActorRole                   = New(200011, "Unknown actor role")
	CoordinatorServerAlreadyExist                 = New(200012, "Server already exists")
)
