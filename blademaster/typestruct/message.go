package typestruct

const (
	MessageCongratulate    = 11
	MessageNoticeImportant = 20
	MessageDialogBox       = 21
	MessageNotice          = 22
	MessageNotice2         = 23
	MessageDialogBoxExit   = 60
)

var (
	GAME_SERVER_ERROR                    = []byte("#CSO2_ServerMessage_ServerDown")
	GAME_LOGIN_ALREADY                   = []byte("您的账号当前已经有人登录！如果你的账号已泄露请联系管理员！")
	GAME_LOGIN_EXIT_FORCE                = []byte("您的账号当前有人登录！您已被强制退出！")
	GAME_LOGIN_ERROR                     = []byte("#CSO2_LOGIN_ERROR")
	GAME_ROOM_COUNT_MODE_ERROR           = []byte("无法开始游戏！请检查你的房间设置，比如是否有2人及以上或者开启bot")
	GAME_ROOM_JOIN_ERROR                 = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_INVALID_GAME_IP")
	GAME_GM_ADD_ALLWEAPONS               = []byte("所有武器已存入你的仓库!")
	GAME_CHANNEL_MESSAGE_NOT_IN          = []byte("你未进入任何频道!")
	GAME_GM_NO_AUTHORIZE                 = []byte("你没有权限进行此操作!")
	GAME_USER_NEW_ITEM                   = []byte("你获得了管理员给予的物品!")
	GAME_USER_NEW_ITEM_RESTART           = []byte("重启游戏以获得奖励的物品!")
	GAME_CHANNEL_MESSAGE                 = []byte("#CSO2_UI_MEGAPHONE_TYPE_02_01")
	GAME_ROOM_JOIN_FAILED_CLOSED         = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_CLOSED")
	GAME_ROOM_JOIN_FAILED_FULL           = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_FULL")
	GAME_ROOM_JOIN_FAILED_BAD_PASSWORD   = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_INVALID_PASSWD")
	GAME_ROOM_CHANGETEAM_FAILED          = []byte("#CSO2_POPUP_ROOM_CHANGETEAM_FAILED")
	GAME_ROOM_COUNTDOWN_FAILED_NOENEMIES = []byte("#CSO2_UI_ROOM_COUNTDOWN_FAILED_NOENEMY")
	GAME_LOGIN_BAD_USERINFO              = []byte("#CSO2_LoginAuth_Certify_NoPassport")
	GAME_LOGIN_BAD_PASSWORD              = []byte("#CSO2_LoginAuth_WrongPassword")
	GAME_LOGIN_BAD_USERNAME              = []byte("#CSO2_LoginAuth_UserNotExists")
	GAME_LOGIN_TENTH_FAILED              = []byte("#CSO2_LoginAuth_TempBlockedByLoginFail")
	GAME_LOGIN_INVALID_USERINFO          = []byte("#CSO2_ServerMessage_INVALID_USERINFO")
	GAME_REPORT_USER_SUCCEED             = []byte("#CSO2_POPUP_USER_REPORT_SUCCEED")
	CSO2_BUY_FAIL_CN_0X11_NO_CASH        = []byte("#CSO2_BUY_FAIL_CN_0X11_NO_CASH")
	CSO2_BUY_FAIL_NO_POINT               = []byte("#CSO2_BUY_FAIL_NO_POINT")
	CSO2_BUY_FAIL_NO_MILEAGE             = []byte("#CSO2_BUY_FAIL_NO_MILEAGE")
	CSO2_BUY_SUCCEED                     = []byte("#CSO2_BUY_SUCCEED")
	CSO2_BUY_FAIL_CN_0X12_DB_ERROR       = []byte("#CSO2_BUY_FAIL_CN_0X12_DB_ERROR")
	CSO2_UI_RAMDOMBOX_ALERT_000          = []byte("#CSO2_UI_RAMDOMBOX_ALERT_000")
	CSO2_POPUP_ROOM_LEAVE_PENALTY        = []byte("CSO2_POPUP_ROOM_LEAVE_PENALTY")
)
