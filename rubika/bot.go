package rubika

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/m4hdi1020/rubigo/encryption"
)

const (
	chatUpdatesMethod           = "getChatsUpdates"
	appVersion                  = "4.2.0"
	appName                     = "Main"
	platform                    = "Web"
	packAge                     = "web.rubika.ir"
	langcode                    = "fa"
	jsonContentType             = "application/json"
	apiVersion                  = "5"
	sendMessage                 = "sendMessage"
	editMessage                 = "editMessage"
	deleteMessage               = "deleteMessages"
	createPollMethod            = "createPoll"
	sendFileMethod              = "requestSendFile"
	webSocketMethod             = "handShake"
	getLink                     = "getLinkFromAppUrl"
	joinGroup                   = "joinGroup"
	leaveGroup                  = "leaveGroup"
	removeMemberMethod          = "banGroupMember"
	pinMessageMethod            = "setPinMessage"
	getUserInfoMethod           = "getUserInfo"
	blockUserMethod             = "setBlockUser"
	block                       = "Block"
	Unblock                     = "Unblock"
	deleteUserChat              = "deleteUserChat"
	forwardMessageMethod        = "forwardMessages"
	getGroupInfoMethod          = "getGroupInfo"
	deleteChatHistoryMethod     = "deleteChatHistory"
	getInfoByIdMethod           = "getObjectByUsername"
	getChannelInfoMethod        = "getChannelInfo"
	getGroupAdminMembersMethod  = "getGroupAdminMembers"
	getAllGroupMembersMethod    = "getGroupAllMembers"
	getChannelAllMembersMethod  = "getChannelAllMembers"
	addGroupAdminMethod         = "setGroupAdmin"
	AdminChangeInfoAccess       = "ChangeInfo"
	AdminPinMessageAccess       = "PinMessages"
	AdminDeleteGlobalMessage    = "DeleteGlobalAllMessages"
	AdminBanMember              = "BanMember"
	AdminSetJoinLink            = "SetJoinLink"
	AdminSetAdmin               = "SetAdmin"
	AdminSetMemberAccess        = "SetMemberAccess"
	groupAccessMethod           = "setGroupDefaultAccess"
	AccessGroupAddMember        = "AddMember"
	AccessGroupViewAdmins       = "ViewAdmins"
	AccessGroupSendMessage      = "SendMessages"
	AccessGroupViewMembers      = "ViewMembers"
	getGroupLinkMethod          = "getGroupLink"
	getChannelAdminsMethod      = "getChannelAdminMembers"
	getMessageInfoMethod        = "getMessagesByID"
	getBlockedUsersListMethod   = "getBlockedUsers"
	getBannedGroupMembersMethod = "getBannedGroupMembers"
)

var (
	rubikaAPIList = [4]string{"https://messengerg2c32.iranlms.ir/", "https://messengerg2c201.iranlms.ir/", "https://messengerg2c171.iranlms.ir/", "https://messengerg2c146.iranlms.ir/"}
	clientVal     = clientValue{AppName: appName, AppVersion: appVersion, Platform: platform, Package: packAge, LangCode: langcode}
)

type bot struct {
	Auth string
}

type send struct {
	ApiVersion string `json:"api_version"`
	Auth       string `json:"auth"`
	DataEnc    string `json:"data_enc"`
}

type getResponseChatUpdates struct {
	Status    string  `json:"status"`
	StatusDet string  `json:"status_det"`
	Data      getData `json:"data"`
}

type getData struct {
	Chats     []getChats `json:"chats"`
	Newstate  int        `json:"new_state"`
	Status    string     `json:"status"`
	Timestamp string     `json:"timestamp"`
}

type getChats struct {
	Guid                string         `json:"object_guid"`
	Access              []string       `json:"access"`
	CountUnseen         int            `json:"count_unseen"`
	IsMute              bool           `json:"is_mute"`
	IsPinned            bool           `json:"is_pinned"`
	TimeString          string         `json:"time_string"`
	LastMessage         getLastMessage `json:"last_message"`
	LastSeenMyMid       string         `json:"last_seen_my_mid"`
	LastSeenPeerMid     string         `json:"455862862947497"`
	Status              string         `json:"status"`
	Time                int            `json:"time"`
	AbsObject           getAbsObject   `json:"abs_object"`
	IsBlocked           bool           `json:"is_blocked"`
	LastMessageId       string         `json:"last_message_id"`
	LastDeletedMid      string         `json:"last_deleted_mid"`
	SlowModeDuration    int            `json:"slow_mode_duration"`
	GroupMyLastSendTime int            `json:"group_my_last_send_time"`
}

type getLastMessage struct {
	MessageId        string `json:"message_id"`
	Type             string `json:"type"`
	Text             string `json:"text"`
	AuthorObjectGuid string `json:"author_object_guid"`
	IsMine           bool   `json:"is_mine"`
	AuthorTitle      string `json:"author_title"`
	AuthorType       string `json:"author_type"`
}

type getAbsObject struct {
	ObjectGuid      string             `json:"object_guid"`
	Type            string             `json:"type"`
	Title           string             `json:"title"`
	AvatarThumbnail getAvatarThumbnail `json:"avatar_thumbnail"`
	IsVerified      bool               `json:"is_verified"`
	IsDeleted       bool               `json:"is_deleted"`
}

type getAvatarThumbnail struct {
	FileId        string `json:"avatar_thumbnail"`
	Mime          string `json:"mime"`
	DcId          string `json:"dc_id"`
	AccessHashRec string `json:"access_hash_rec"`
}

type SendReqChatUpdates struct {
	Method string    `json:"method"`
	Input  inputStr  `json:"input"`
	Client clientStr `json:"client"`
}

type inputStr struct {
	State int `json:"state"`
}

type clientStr struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	Platform   string `json:"platform"`
	Package    string `json:"package"`
	LangCode   string `json:"lang_code"`
}

type sendMessagePayload struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGuid     string `json:"object_guid"`
		Rnd            string `json:"rnd"`
		Text           string `json:"text,omitempty"`
		ReplyToMessage string `json:"reply_to_message_id,omitempty"`
	} `json:"input"`
	Clinet struct {
		AppName    string `json:"app_name"`
		AppVersion string `json:"app_version"`
		Platform   string `json:"platform"`
		Package    string `json:"package"`
		LangCode   string `json:"lang_code"`
	} `json:"client"`
}

type EditText struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID string `json:"object_guid"`
		MessageID  string `json:"message_id"`
		Text       string `json:"text"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type deleteMessageStruct struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID string   `json:"object_guid"`
		MessageIds []string `json:"message_ids"`
		Type       string   `json:"type"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type createPoll struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID            string   `json:"object_guid"`
		Options               []string `json:"options"`
		Rnd                   string   `json:"rnd"`
		Question              string   `json:"question"`
		Type                  string   `json:"type"`
		IsAnonymous           bool     `json:"is_anonymous"`
		AllowsMultipleAnswers bool     `json:"allows_multiple_answers"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type infoSendFile struct {
	Method string `json:"method"`
	Input  struct {
		FileName string `json:"file_name"`
		Size     int    `json:"size"`
		Mime     string `json:"mime"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getRespSendFile struct {
	Status    string `json:"status"`
	StatusDet string `json:"status_det"`
	Data      struct {
		ID             string `json:"id"`
		DcID           string `json:"dc_id"`
		AccessHashSend string `json:"access_hash_send"`
		UploadURL      string `json:"upload_url"`
	} `json:"data"`
}

type getHashReq struct {
	Status       string `json:"status"`
	StatusDet    string `json:"status_det"`
	DevMessage   any    `json:"dev_message"`
	UserMesssage any    `json:"user_messsage"`
	Data         struct {
		AccessHashRec string `json:"access_hash_rec"`
	} `json:"data"`
	DataEnc any `json:"data_enc"`
}

type sendFile struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID string `json:"object_guid"`
		Rnd        string `json:"rnd"`
		FileInline struct {
			DcID          string `json:"dc_id"`
			FileID        string `json:"file_id"`
			Type          string `json:"type"`
			FileName      string `json:"file_name"`
			Size          int    `json:"size"`
			Mime          string `json:"mime"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"file_inline"`
		Text             string `json:"text,omitempty"`
		ReplyToMessageId string `json:"reply_to_message_id,omitempty"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type webSocketResponse struct {
	Chats []WebSocketResponse `json:"chat_updates"`
}
type WebSocketResponse struct {
	ObjectGuid string            `json:"object_guid"`
	Action     string            `json:"action"`
	Chat       getWebSocketChats `json:"chat"`
	TimeStamp  string            `json:"timestamp"`
	Type       string            `json:"type"`
}
type getWebSocketChats struct {
	TimeString      string         `json:"time_string"`
	LastMessage     getLastMessage `json:"last_message"`
	LastSeenPeerMid string         `json:"last_seen_peer_mid"`
	Time            int            `json:"time"`
	LastMessageId   string         `json:"last_message_id"`
}

type webSocketPayload struct {
	Method string   `json:"method"`
	Input  struct{} `json:"input"`
	Clinet struct {
		AppName    string `json:"app_name"`
		AppVersion string `json:"app_version"`
		Platform   string `json:"platform"`
		Package    string `json:"package"`
		LangCode   string `json:"lang_code"`
	} `json:"client"`
}

type imageData struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID string `json:"object_guid"`
		Rnd        string `json:"rnd"`
		FileInline struct {
			DcID          string `json:"dc_id"`
			FileID        string `json:"file_id"`
			Type          string `json:"type"`
			FileName      string `json:"file_name"`
			Size          int    `json:"size"`
			Mime          string `json:"mime"`
			ThumbInline   string `json:"thumb_inline"`
			Width         int    `json:"width"`
			Height        int    `json:"height"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"file_inline"`
		Text           string `json:"text,omitempty"`
		ReplyToMessage string `json:"reply_to_message_id,omitempty"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type joinGroupData struct {
	Method string `json:"method"`
	Input  struct {
		HashLink string `json:"hash_link"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type joinGroupReponse struct {
	Status    string `json:"status"`
	StatusDet string `json:"status_det"`
	Data      struct {
		Group struct {
			GroupGUID       string `json:"group_guid"`
			GroupTitle      string `json:"group_title"`
			AvatarThumbnail struct {
				FileID        string `json:"file_id"`
				Mime          string `json:"mime"`
				DcID          string `json:"dc_id"`
				AccessHashRec string `json:"access_hash_rec"`
			} `json:"avatar_thumbnail"`
			CountMembers             int    `json:"count_members"`
			IsDeleted                bool   `json:"is_deleted"`
			IsVerified               bool   `json:"is_verified"`
			SlowMode                 int    `json:"slow_mode"`
			ChatHistoryForNewMembers string `json:"chat_history_for_new_members"`
			EventMessages            bool   `json:"event_messages"`
		} `json:"group"`
		IsValid    bool `json:"is_valid"`
		ChatUpdate struct {
			ObjectGUID string `json:"object_guid"`
			Action     string `json:"action"`
			Chat       struct {
				ObjectGUID  string   `json:"object_guid"`
				Access      []string `json:"access"`
				CountUnseen int      `json:"count_unseen"`
				IsMute      bool     `json:"is_mute"`
				IsPinned    bool     `json:"is_pinned"`
				TimeString  string   `json:"time_string"`
				LastMessage struct {
					MessageID string `json:"message_id"`
					Type      string `json:"type"`
					Text      string `json:"text"`
					IsMine    bool   `json:"is_mine"`
				} `json:"last_message"`
				LastSeenMyMid   string `json:"last_seen_my_mid"`
				LastSeenPeerMid string `json:"last_seen_peer_mid"`
				Status          string `json:"status"`
				Time            int    `json:"time"`
				AbsObject       struct {
					ObjectGUID      string `json:"object_guid"`
					Type            string `json:"type"`
					Title           string `json:"title"`
					AvatarThumbnail struct {
						FileID        string `json:"file_id"`
						Mime          string `json:"mime"`
						DcID          string `json:"dc_id"`
						AccessHashRec string `json:"access_hash_rec"`
					} `json:"avatar_thumbnail"`
					IsVerified bool `json:"is_verified"`
					IsDeleted  bool `json:"is_deleted"`
				} `json:"abs_object"`
				IsBlocked      bool   `json:"is_blocked"`
				LastMessageID  string `json:"last_message_id"`
				LastDeletedMid string `json:"last_deleted_mid"`
			} `json:"chat"`
			UpdatedParameters []any  `json:"updated_parameters"`
			Timestamp         string `json:"timestamp"`
			Type              string `json:"type"`
		} `json:"chat_update"`
		MessageUpdate struct {
			MessageID string `json:"message_id"`
			Action    string `json:"action"`
			Message   struct {
				MessageID string `json:"message_id"`
				Text      string `json:"text"`
				Time      string `json:"time"`
				IsEdited  bool   `json:"is_edited"`
				Type      string `json:"type"`
				EventData struct {
					Type            string `json:"type"`
					PerformerObject struct {
						Type       string `json:"type"`
						ObjectGUID string `json:"object_guid"`
					} `json:"performer_object"`
				} `json:"event_data"`
			} `json:"message"`
			UpdatedParameters []any  `json:"updated_parameters"`
			Timestamp         string `json:"timestamp"`
			PrevMessageID     string `json:"prev_message_id"`
			ObjectGUID        string `json:"object_guid"`
			Type              string `json:"type"`
			State             string `json:"state"`
		} `json:"message_update"`
		Timestamp string `json:"timestamp"`
	} `json:"data"`
}

type leaveGroupData struct {
	Method string `json:"method"`
	Input  struct {
		Guid string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type removeMemberPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGuid  string `json:"group_guid"`
		MemberGuid string `json:"member_guid"`
		Action     string `json:"action"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type pinMessage struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID string `json:"object_guid"`
		MessageID  string `json:"message_id"`
		Action     string `json:"action"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getUserInfoPayload struct {
	Method string `json:"method"`
	Input  struct {
		UserGUID string `json:"user_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type userInfoData struct {
	Status    string   `json:"status"`
	StatusDet string   `json:"status_det"`
	Data      userInfo `json:"data"`
}

type userInfo struct {
	User struct {
		UserGUID   string `json:"user_guid"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Username   string `json:"username"`
		LastOnline int    `json:"last_online"`
		Bio        string `json:"bio"`
		IsDeleted  bool   `json:"is_deleted"`
		IsVerified bool   `json:"is_verified"`
		OnlineTime struct {
			Type              string `json:"type"`
			ApproximatePeriod string `json:"approximate_period"`
		} `json:"online_time"`
	} `json:"user"`
	Chat struct {
		ObjectGUID string   `json:"object_guid"`
		Access     []string `json:"access"`
		IsMute     bool     `json:"is_mute"`
		Status     string   `json:"status"`
		AbsObject  struct {
			ObjectGUID string `json:"object_guid"`
			Type       string `json:"type"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			IsVerified bool   `json:"is_verified"`
			IsDeleted  bool   `json:"is_deleted"`
		} `json:"abs_object"`
		IsBlocked   bool `json:"is_blocked"`
		IsInContact bool `json:"is_in_contact"`
	} `json:"chat"`
	Timestamp         string `json:"timestamp"`
	IsInContact       bool   `json:"is_in_contact"`
	CountCommonGroups int    `json:"count_common_groups"`
	CanReceiveCall    bool   `json:"can_receive_call"`
	CanVideoCall      bool   `json:"can_video_call"`
}

type blockUserPayload struct {
	Method string `json:"method"`
	Input  struct {
		UserGUID string `json:"user_guid"`
		Action   string `json:"action"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type deleteUserChatData struct {
	Method string `json:"method"`
	Input  struct {
		UserGUID             string `json:"user_guid"`
		LastDeletedMessageID string `json:"last_deleted_message_id"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type forwardMessageData struct {
	Method string `json:"method"`
	Input  struct {
		FromObjectGUID string   `json:"from_object_guid"`
		ToObjectGUID   string   `json:"to_object_guid"`
		MessageIds     []string `json:"message_ids"`
		Rnd            string   `json:"rnd"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getGroupInfoData struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getGroupInfo struct {
	Status    string    `json:"status"`
	StatusDet string    `json:"status_det"`
	Data      groupInfo `json:"data"`
}
type groupInfo struct {
	Group struct {
		GroupGUID       string `json:"group_guid"`
		GroupTitle      string `json:"group_title"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		CountMembers             int    `json:"count_members"`
		IsDeleted                bool   `json:"is_deleted"`
		IsVerified               bool   `json:"is_verified"`
		SlowMode                 int    `json:"slow_mode"`
		Description              string `json:"description"`
		ChatHistoryForNewMembers string `json:"chat_history_for_new_members"`
		EventMessages            bool   `json:"event_messages"`
	} `json:"group"`
	Chat struct {
		ObjectGUID  string   `json:"object_guid"`
		Access      []string `json:"access"`
		CountUnseen int      `json:"count_unseen"`
		IsMute      bool     `json:"is_mute"`
		IsPinned    bool     `json:"is_pinned"`
		TimeString  string   `json:"time_string"`
		LastMessage struct {
			MessageID        string `json:"message_id"`
			Type             string `json:"type"`
			Text             string `json:"text"`
			AuthorObjectGUID string `json:"author_object_guid"`
			IsMine           bool   `json:"is_mine"`
			AuthorTitle      string `json:"author_title"`
			AuthorType       string `json:"author_type"`
		} `json:"last_message"`
		LastSeenMyMid   string `json:"last_seen_my_mid"`
		LastSeenPeerMid string `json:"last_seen_peer_mid"`
		Status          string `json:"status"`
		Time            int    `json:"time"`
		PinnedMessageID string `json:"pinned_message_id"`
		AbsObject       struct {
			ObjectGUID      string `json:"object_guid"`
			Type            string `json:"type"`
			Title           string `json:"title"`
			AvatarThumbnail struct {
				FileID        string `json:"file_id"`
				Mime          string `json:"mime"`
				DcID          string `json:"dc_id"`
				AccessHashRec string `json:"access_hash_rec"`
			} `json:"avatar_thumbnail"`
			IsVerified bool `json:"is_verified"`
			IsDeleted  bool `json:"is_deleted"`
		} `json:"abs_object"`
		IsBlocked      bool   `json:"is_blocked"`
		LastMessageID  string `json:"last_message_id"`
		LastDeletedMid string `json:"last_deleted_mid"`
	} `json:"chat"`
	Timestamp string `json:"timestamp"`
}
type deleteChatHistoryPayload struct {
	Method string `json:"method"`
	Input  struct {
		ObjectGUID    string `json:"object_guid"`
		LastMessageID string `json:"last_message_id"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getInfoByUsernamePayload struct {
	Method string `json:"method"`
	Input  struct {
		Username string `json:"username"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getResponseInfoById struct {
	Status    string         `json:"status"`
	StatusDet string         `json:"status_det"`
	Data      infoByUsername `json:"data"`
}

type infoByUsername struct {
	Exist bool   `json:"exist"`
	Type  string `json:"type"`
	User  struct {
		UserGUID        string `json:"user_guid"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Username        string `json:"username"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		LastOnline int    `json:"last_online"`
		Bio        string `json:"bio"`
		IsDeleted  bool   `json:"is_deleted"`
		IsVerified bool   `json:"is_verified"`
		OnlineTime struct {
			Type              string `json:"type"`
			ApproximatePeriod string `json:"approximate_period"`
		} `json:"online_time"`
	} `json:"user"`
	Chat struct {
		ObjectGUID string   `json:"object_guid"`
		Access     []string `json:"access"`
		IsMute     bool     `json:"is_mute"`
		Status     string   `json:"status"`
		AbsObject  struct {
			ObjectGUID      string `json:"object_guid"`
			Type            string `json:"type"`
			FirstName       string `json:"first_name"`
			LastName        string `json:"last_name"`
			AvatarThumbnail struct {
				FileID        string `json:"file_id"`
				Mime          string `json:"mime"`
				DcID          string `json:"dc_id"`
				AccessHashRec string `json:"access_hash_rec"`
			} `json:"avatar_thumbnail"`
			IsVerified bool `json:"is_verified"`
			IsDeleted  bool `json:"is_deleted"`
		} `json:"abs_object"`
		IsBlocked   bool `json:"is_blocked"`
		IsInContact bool `json:"is_in_contact"`
	} `json:"chat"`
	Timestamp   string `json:"timestamp"`
	IsInContact bool   `json:"is_in_contact"`
}

type channelInfoPayload struct {
	Method string `json:"method"`
	Input  struct {
		ChannelGUID string `json:"channel_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type channelInfo struct {
	Status    string          `json:"status"`
	StatusDet string          `json:"status_det"`
	Data      channelInfoData `json:"data"`
}

type channelInfoData struct {
	Channel struct {
		ChannelGUID     string `json:"channel_guid"`
		ChannelTitle    string `json:"channel_title"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		CountMembers int    `json:"count_members"`
		Description  string `json:"description"`
		Username     string `json:"username"`
		IsDeleted    bool   `json:"is_deleted"`
		IsVerified   bool   `json:"is_verified"`
		ShareURL     string `json:"share_url"`
		ChannelType  string `json:"channel_type"`
		SignMessages bool   `json:"sign_messages"`
	} `json:"channel"`
	Chat struct {
		ObjectGUID  string   `json:"object_guid"`
		Access      []string `json:"access"`
		CountUnseen int      `json:"count_unseen"`
		IsMute      bool     `json:"is_mute"`
		IsPinned    bool     `json:"is_pinned"`
		TimeString  string   `json:"time_string"`
		LastMessage struct {
			MessageID string `json:"message_id"`
			Type      string `json:"type"`
			Text      string `json:"text"`
			IsMine    bool   `json:"is_mine"`
		} `json:"last_message"`
		LastSeenMyMid   string `json:"last_seen_my_mid"`
		LastSeenPeerMid string `json:"last_seen_peer_mid"`
		Status          string `json:"status"`
		Time            int    `json:"time"`
		PinnedMessageID string `json:"pinned_message_id"`
		AbsObject       struct {
			ObjectGUID      string `json:"object_guid"`
			Type            string `json:"type"`
			Title           string `json:"title"`
			AvatarThumbnail struct {
				FileID        string `json:"file_id"`
				Mime          string `json:"mime"`
				DcID          string `json:"dc_id"`
				AccessHashRec string `json:"access_hash_rec"`
			} `json:"avatar_thumbnail"`
			IsVerified bool `json:"is_verified"`
			IsDeleted  bool `json:"is_deleted"`
		} `json:"abs_object"`
		IsBlocked      bool   `json:"is_blocked"`
		LastMessageID  string `json:"last_message_id"`
		LastDeletedMid string `json:"last_deleted_mid"`
	} `json:"chat"`
	Timestamp string `json:"timestamp"`
}

type groupAdminMembersPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type adminMembers struct {
	Status    string           `json:"status"`
	StatusDet string           `json:"status_det"`
	Data      adminMembersData `json:"data"`
}

type adminMembersData struct {
	InChatMembers []struct {
		MemberType      string `json:"member_type"`
		MemberGUID      string `json:"member_guid"`
		FirstName       string `json:"first_name"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		IsVerified bool   `json:"is_verified"`
		IsDeleted  bool   `json:"is_deleted"`
		LastOnline int    `json:"last_online"`
		JoinType   string `json:"join_type"`
		Username   string `json:"username"`
		OnlineTime struct {
			Type      string `json:"type"`
			ExactTime int    `json:"exact_time"`
		} `json:"online_time"`
	} `json:"in_chat_members"`
	NextStartID string `json:"next_start_id"`
	HasContinue bool   `json:"has_continue"`
	Timestamp   string `json:"timestamp"`
}

type getAllGroupMembersPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getAllGroupMembers struct {
	Status    string              `json:"status"`
	StatusDet string              `json:"status_det"`
	Data      allGroupMembersData `json:"data"`
}

type allGroupMembersData struct {
	InChatMembers []struct {
		MemberType      string `json:"member_type"`
		MemberGUID      string `json:"member_guid"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name,omitempty"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		IsVerified bool   `json:"is_verified"`
		IsDeleted  bool   `json:"is_deleted"`
		LastOnline int    `json:"last_online"`
		JoinType   string `json:"join_type"`
		Username   string `json:"username"`
	} `json:"in_chat_members"`
	NextStartID string `json:"next_start_id"`
	HasContinue bool   `json:"has_continue"`
	Timestamp   string `json:"timestamp"`
}

type channelAllMembersPayload struct {
	Method string `json:"method"`
	Input  struct {
		ChannelGUID string `json:"channel_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getChannelAllMembers struct {
	Status    string             `json:"status"`
	StatusDet string             `json:"status_det"`
	Data      channelMembersData `json:"data"`
}

type channelMembersData struct {
	InChatMembers []struct {
		MemberType      string `json:"member_type"`
		MemberGUID      string `json:"member_guid"`
		FirstName       string `json:"first_name"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		IsVerified bool   `json:"is_verified"`
		IsDeleted  bool   `json:"is_deleted"`
		LastOnline int    `json:"last_online"`
		JoinType   string `json:"join_type"`
		Username   string `json:"username"`
		OnlineTime struct {
			Type      string `json:"type"`
			ExactTime int    `json:"exact_time"`
		} `json:"online_time"`
	} `json:"in_chat_members"`
	NextStartID string `json:"next_start_id"`
	HasContinue bool   `json:"has_continue"`
	Timestamp   string `json:"timestamp"`
}

type addGroupAdminPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID  string   `json:"group_guid"`
		MemberGUID string   `json:"member_guid"`
		Action     string   `json:"action"`
		AccessList []string `json:"access_list"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type removeGroupAdminPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID  string `json:"group_guid"`
		MemberGUID string `json:"member_guid"`
		Action     string `json:"action"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type groupAccessPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID  string   `json:"group_guid"`
		AccessList []string `json:"access_list"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getGroupLinkPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGUID string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getGroupLink struct {
	Status    string `json:"status"`
	StatusDet string `json:"status_det"`
	Data      struct {
		JoinLink string `json:"join_link"`
	} `json:"data"`
}

type getChannelAdminsPayload struct {
	Method string `json:"method"`
	Input  struct {
		ChannelGUID string `json:"channel_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type channelAdminsInfo struct {
	Status    string        `json:"status"`
	StatusDet string        `json:"status_det"`
	Data      channelAdmins `json:"data"`
}

type channelAdmins struct {
	InChatMembers []struct {
		MemberType      string `json:"member_type"`
		MemberGUID      string `json:"member_guid"`
		FirstName       string `json:"first_name"`
		AvatarThumbnail struct {
			FileID        string `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          string `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
		} `json:"avatar_thumbnail"`
		IsVerified bool   `json:"is_verified"`
		IsDeleted  bool   `json:"is_deleted"`
		LastOnline int    `json:"last_online"`
		JoinType   string `json:"join_type"`
		Username   string `json:"username"`
		OnlineTime struct {
			Type      string `json:"type"`
			ExactTime int    `json:"exact_time"`
		} `json:"online_time"`
	} `json:"in_chat_members"`
	NextStartID string `json:"next_start_id"`
	HasContinue bool   `json:"has_continue"`
	Timestamp   string `json:"timestamp"`
}

type getMessageByIDPayload struct {
	Method string `json:"method"`
	Input  struct {
		Guid        string   `json:"object_guid"`
		Message_Ids []string `json:"message_ids"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type getMessageInfoByID struct {
	Status    string             `json:"status"`
	StatusDet string             `json:"status_det"`
	Data      getMessageInfoData `json:"data"`
}

type getMessageInfoData struct {
	Messages []struct {
		MessageID  string `json:"message_id"`
		Text       string `json:"text"`
		FileInline struct {
			FileID        int64  `json:"file_id"`
			Mime          string `json:"mime"`
			DcID          int    `json:"dc_id"`
			AccessHashRec string `json:"access_hash_rec"`
			FileName      string `json:"file_name"`
			Width         int    `json:"width"`
			Height        int    `json:"height"`
			Time          int    `json:"time"`
			Size          int    `json:"size"`
			Type          string `json:"type"`
		} `json:"file_inline"`
		Time          string `json:"time"`
		IsEdited      bool   `json:"is_edited"`
		ForwardedFrom struct {
			TypeFrom   string `json:"type_from"`
			MessageID  string `json:"message_id"`
			ObjectGUID string `json:"object_guid"`
		} `json:"forwarded_from"`
		Type             string `json:"type"`
		AuthorType       string `json:"author_type"`
		AuthorObjectGUID string `json:"author_object_guid"`
	} `json:"messages"`
	Timestamp string `json:"timestamp"`
}

type clientValue struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	Platform   string `json:"platform"`
	Package    string `json:"package"`
	LangCode   string `json:"lang_code"`
}

type getBlockedUsersPayload struct {
	Method string      `json:"method"`
	Input  struct{}    `json:"input"`
	Client clientValue `json:"client"`
}

type getBlockedUsersResponse struct {
	Status    string `json:"status"`
	StatusDet string `json:"status_det"`
	Data      struct {
		AbsUsers    []blockedUsers `json:"abs_users"`
		NextStartID string         `json:"next_start_id"`
		HasContinue bool           `json:"has_continue"`
		Timestamp   string         `json:"timestamp"`
	} `json:"data"`
}

type blockedUsers struct {
	ObjectGUID      string `json:"object_guid"`
	Type            string `json:"type"`
	FirstName       string `json:"first_name"`
	IsVerified      bool   `json:"is_verified"`
	IsDeleted       bool   `json:"is_deleted"`
	LastName        string `json:"last_name,omitempty"`
	AvatarThumbnail struct {
		FileID        string `json:"file_id"`
		Mime          string `json:"mime"`
		DcID          string `json:"dc_id"`
		AccessHashRec string `json:"access_hash_rec"`
	} `json:"avatar_thumbnail,omitempty"`
}

type bannedGroupMembersPayload struct {
	Method string `json:"method"`
	Input  struct {
		GroupGuid string `json:"group_guid"`
	} `json:"input"`
	Client clientValue `json:"client"`
}

type bannedGroupMembersResp struct {
	Status    string `json:"status"`
	StatusDet string `json:"status_det"`
	Data      struct {
		InChatMembers []bannedList `json:"in_chat_members"`
		NextStartID   string       `json:"next_start_id"`
		HasContinue   bool         `json:"has_continue"`
		Timestamp     string       `json:"timestamp"`
	} `json:"data"`
}

type bannedList struct {
	MemberType      string `json:"member_type"`
	MemberGUID      string `json:"member_guid"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	AvatarThumbnail struct {
		FileID        string `json:"file_id"`
		Mime          string `json:"mime"`
		DcID          string `json:"dc_id"`
		AccessHashRec string `json:"access_hash_rec"`
	} `json:"avatar_thumbnail"`
	IsVerified          bool   `json:"is_verified"`
	IsDeleted           bool   `json:"is_deleted"`
	LastOnline          int    `json:"last_online"`
	RemovedByObjectGUID string `json:"removed_by_object_guid"`
	RemovedByObjectType string `json:"removed_by_object_type"`
	JoinType            string `json:"join_type"`
	Username            string `json:"username"`
	OnlineTime          struct {
		Type              string `json:"type"`
		ApproximatePeriod string `json:"approximate_period"`
	} `json:"online_time,omitempty"`
	OnlineTime0 struct {
		Type      string `json:"type"`
		ExactTime int    `json:"exact_time"`
	} `json:"online_time,omitempty"`
	OnlineTime1 struct {
		Type      string `json:"type"`
		ExactTime int    `json:"exact_time"`
	} `json:"online_time,omitempty"`
	OnlineTime2 struct {
		Type      string `json:"type"`
		ExactTime int    `json:"exact_time"`
	} `json:"online_time,omitempty"`
	OnlineTime3 struct {
		Type      string `json:"type"`
		ExactTime int    `json:"exact_time"`
	} `json:"online_time,omitempty"`
}

func newBannedListReq(guid string) (string, error) {
	data := bannedGroupMembersPayload{
		Method: getBannedGroupMembersMethod,
		Input: struct {
			GroupGuid string `json:"group_guid"`
		}{GroupGuid: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetBlockedUsersPayload() (string, error) {
	data := getBlockedUsersPayload{
		Method: getBlockedUsersListMethod,
		Input:  struct{}{},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}
func newGetMessageInfoByID(guid string, messageIds ...string) (string, error) {
	var messageIdsList []string
	messageIdsList = append(messageIdsList, messageIds...)

	data := getMessageByIDPayload{
		Method: getMessageInfoMethod,
		Input: struct {
			Guid        string   "json:\"object_guid\""
			Message_Ids []string "json:\"message_ids\""
		}{Guid: guid, Message_Ids: messageIdsList},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetChannelAdmins(guid string) (string, error) {
	data := getChannelAdminsPayload{
		Method: getChannelAdminsMethod,
		Input: struct {
			ChannelGUID string "json:\"channel_guid\""
		}{ChannelGUID: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetGroupLink(guid string) (string, error) {
	data := getGroupLinkPayload{
		Method: getGroupLinkMethod,
		Input: struct {
			GroupGUID string "json:\"group_guid\""
		}{GroupGUID: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newSetGroupAccess(groupGuid string, AccessList []string) (string, error) {
	data := groupAccessPayload{
		Method: groupAccessMethod,
		Input: struct {
			GroupGUID  string   "json:\"group_guid\""
			AccessList []string "json:\"access_list\""
		}{GroupGUID: groupGuid, AccessList: AccessList},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func NewBot(auth string) (bot, error) {
	if len(auth) != 32 {
		return bot{}, fmt.Errorf("error: your auth is invalid :(")
	}
	encryption.Secret(auth)
	return bot{Auth: auth}, nil
}

func newRemoveGroupAdmin(groupGuid, memberGuid string) (string, error) {
	data := removeGroupAdminPayload{
		Method: addGroupAdminMethod,
		Input: struct {
			GroupGUID  string "json:\"group_guid\""
			MemberGUID string "json:\"member_guid\""
			Action     string "json:\"action\""
		}{GroupGUID: groupGuid, MemberGUID: memberGuid, Action: "UnsetAdmin"},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newAddGroupAdmin(groupGuid string, memberGuid string, accessList []string) (string, error) {
	data := addGroupAdminPayload{
		Method: addGroupAdminMethod,
		Input: struct {
			GroupGUID  string   "json:\"group_guid\""
			MemberGUID string   "json:\"member_guid\""
			Action     string   "json:\"action\""
			AccessList []string "json:\"access_list\""
		}{GroupGUID: groupGuid, MemberGUID: memberGuid, Action: "SetAdmin", AccessList: accessList},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetChannelAllMembers(guid string) (string, error) {
	data := channelAllMembersPayload{
		Method: getChannelAllMembersMethod,
		Input: struct {
			ChannelGUID string "json:\"channel_guid\""
		}{ChannelGUID: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetAllGroupMembers(guid string) (string, error) {
	data := getAllGroupMembersPayload{
		Method: getAllGroupMembersMethod,
		Input: struct {
			GroupGUID string "json:\"group_guid\""
		}{guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetGroupAdminMembers(guid string) (string, error) {
	data := groupAdminMembersPayload{
		Method: getGroupAdminMembersMethod,
		Input: struct {
			GroupGUID string "json:\"group_guid\""
		}{guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newChannelInfo(guid string) (string, error) {
	data := channelInfoPayload{
		Method: getChannelInfoMethod,
		Input: struct {
			ChannelGUID string "json:\"channel_guid\""
		}{ChannelGUID: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGetInfoById(username string) (string, error) {
	data := getInfoByUsernamePayload{
		Method: getInfoByIdMethod,
		Input: struct {
			Username string "json:\"username\""
		}{Username: username},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newDeleteChatHistory(guid string, lastMessageId string) (string, error) {
	data := deleteChatHistoryPayload{
		Method: deleteChatHistoryMethod,
		Input: struct {
			ObjectGUID    string "json:\"object_guid\""
			LastMessageID string "json:\"last_message_id\""
		}{guid, lastMessageId},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newGroupInfo(groupGuid string) (string, error) {
	data := getGroupInfoData{
		Method: getGroupInfoMethod,
		Input: struct {
			GroupGUID string "json:\"group_guid\""
		}{groupGuid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newForwardMessage(fromGuid, ToGuid string, messageIds []string) (string, error) {
	data := forwardMessageData{
		Method: forwardMessageMethod,
		Input: struct {
			FromObjectGUID string   "json:\"from_object_guid\""
			ToObjectGUID   string   "json:\"to_object_guid\""
			MessageIds     []string "json:\"message_ids\""
			Rnd            string   "json:\"rnd\""
		}{FromObjectGUID: fromGuid, ToObjectGUID: ToGuid, MessageIds: messageIds, Rnd: randNum()},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newDeleteUserChat(userGuid, lastDeletedMessageId string) (string, error) {
	data := deleteUserChatData{
		Method: deleteUserChat,
		Input: struct {
			UserGUID             string "json:\"user_guid\""
			LastDeletedMessageID string "json:\"last_deleted_message_id\""
		}{UserGUID: userGuid, LastDeletedMessageID: lastDeletedMessageId},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newBlockUser(userGuid string, action string) (string, error) {
	data := blockUserPayload{
		Method: blockUserMethod,
		Input: struct {
			UserGUID string "json:\"user_guid\""
			Action   string "json:\"action\""
		}{UserGUID: userGuid, Action: action},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newUserInfo(userGuid string) (string, error) {
	data := getUserInfoPayload{
		Method: getUserInfoMethod,
		Input: struct {
			UserGUID string "json:\"user_guid\""
		}{UserGUID: userGuid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newPinMessage(groupGuid, messageId string) (string, error) {
	data := pinMessage{
		Method: pinMessageMethod,
		Input: struct {
			ObjectGUID string "json:\"object_guid\""
			MessageID  string "json:\"message_id\""
			Action     string "json:\"action\""
		}{groupGuid, messageId, "Pin"},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		log.Fatalln(err)
	}
	return dataEnc, nil
}

func newRemoveMember(groupGuid string, memberGuid string, action string) (string, error) {
	data := removeMemberPayload{
		Method: removeMemberMethod,
		Input: struct {
			GroupGuid  string "json:\"group_guid\""
			MemberGuid string "json:\"member_guid\""
			Action     string "json:\"action\""
		}{GroupGuid: groupGuid, MemberGuid: memberGuid, Action: action},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newLeaveGroup(guid string) (string, error) {
	data := leaveGroupData{
		Method: leaveGroup,
		Input: struct {
			Guid string "json:\"group_guid\""
		}{Guid: guid},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newJoinGroup(hashLink string) (string, error) {
	data := joinGroupData{
		Method: joinGroup,
		Input: struct {
			HashLink string "json:\"hash_link\""
		}{HashLink: hashLink},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newSendFile(text, guid, dcId, fileId, fileName string, size int, accessHashRec string, messageId string) (string, error) {
	data := sendFile{
		Method: sendMessage,
		Input: struct {
			ObjectGUID string "json:\"object_guid\""
			Rnd        string "json:\"rnd\""
			FileInline struct {
				DcID          string "json:\"dc_id\""
				FileID        string "json:\"file_id\""
				Type          string "json:\"type\""
				FileName      string "json:\"file_name\""
				Size          int    "json:\"size\""
				Mime          string "json:\"mime\""
				AccessHashRec string "json:\"access_hash_rec\""
			} "json:\"file_inline\""
			Text             string "json:\"text,omitempty\""
			ReplyToMessageId string "json:\"reply_to_message_id,omitempty\""
		}{ObjectGUID: guid, Rnd: randNum(), FileInline: struct {
			DcID          string "json:\"dc_id\""
			FileID        string "json:\"file_id\""
			Type          string "json:\"type\""
			FileName      string "json:\"file_name\""
			Size          int    "json:\"size\""
			Mime          string "json:\"mime\""
			AccessHashRec string "json:\"access_hash_rec\""
		}{DcID: dcId, FileID: fileId, Type: "File", FileName: fileName, Size: size, Mime: "mime", AccessHashRec: accessHashRec}},
		Client: clientVal,
	}
	if text != "" && text != " " {
		data.Input.Text = text
	}
	if messageId != "" && messageId != " " {
		data.Input.ReplyToMessageId = messageId
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newSendImage(guid string, caption string, dcId string, id string, fileName string, size int, width int, height int, accessHashReq string, messageId string) (string, error) {
	th := "/9j/4AAQSkZJRgABAQAAAQABAAD/4gIoSUNDX1BST0ZJTEUAAQEAAAIYAAAAAAIQAABtbnRyUkdCIFhZWiAAAAAAAAAAAAAAAABhY3NwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADTLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlkZXNjAAAA8AAAAHRyWFlaAAABZAAAABRnWFlaAAABeAAAABRiWFlaAAABjAAAABRyVFJDAAABoAAAAChnVFJDAAABoAAAAChiVFJDAAABoAAAACh3dHB0AAAByAAAABRjcHJ0AAAB3AAAADxtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAFgAAAAcAHMAUgBHAEIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFhZWiAAAAAAAABvogAAOPUAAAOQWFlaIAAAAAAAAGKZAAC3hQAAGNpYWVogAAAAAAAAJKAAAA+EAAC2z3BhcmEAAAAAAAQAAAACZmYAAPKnAAANWQAAE9AAAApbAAAAAAAAAABYWVogAAAAAAAA9tYAAQAAAADTLW1sdWMAAAAAAAAAAQAAAAxlblVTAAAAIAAAABwARwBvAG8AZwBsAGUAIABJAG4AYwAuACAAMgAwADEANv/bAEMAGxIUFxQRGxcWFx4cGyAoQisoJSUoUTo9MEJgVWVkX1VdW2p4mYFqcZBzW12FtYaQnqOrratngLzJuqbHmairpP/bAEMBHB4eKCMoTisrTqRuXW6kpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpKSkpP/AABEIAEEALgMBIgACEQEDEQH/xAAaAAADAQEBAQAAAAAAAAAAAAADBAUAAgYB/8QAJhAAAgICAgEDBAMAAAAAAAAAAQIAAwQREiExBRNBFCIyMzRRgf/EABgBAQEBAQEAAAAAAAAAAAAAAAECAAME/8QAHBEAAwEBAQADAAAAAAAAAAAAAAERAiFBAxIx/9oADAMBAAIRAxEAPwCnhYyqOe9iFysiupD2IC7LTEQ1iedy8qy6w/cdQShetPTpWT1JbmKN/kfrd3TSeJ5IMV7HmVvSc91bi/iEGxDd9TITvzFQhcncoZTFxyUbii1O566ktHfO3KxnNx/cHPW5FuxiG6E9WVQVHlJGQyFyFEs86JQxWI8R7Dw+JBMIltYGm8xmlw2gJimhusIE00BzrNraIAnGTd7Y0fmTrH4tsHzFEN+Ip5bOicSYiWCDkxjeRYHPZgbMYWpDwtRPogzK79mUqAAgYHxJFtL12Edwld1vHgNwKcQ3n2+70vxEWs30fiH3wrPI9xQciSSJn+Gyl9h0uWI3HK7NACfEx1FXIxeyzgDHPET8ndcC5HEdnW4mzMmyB1PjXcvJ7ndlq/TlD+RiTHCcbntu0DHSukH9weNi625jLAaEJSstZ6x8/wAcSbkeJpphFF/YIa78xNNAV6NJ+ucGaaWjkz//2Q=="
	data := imageData{
		Method: sendMessage,
		Input: struct {
			ObjectGUID string "json:\"object_guid\""
			Rnd        string "json:\"rnd\""
			FileInline struct {
				DcID          string "json:\"dc_id\""
				FileID        string "json:\"file_id\""
				Type          string "json:\"type\""
				FileName      string "json:\"file_name\""
				Size          int    "json:\"size\""
				Mime          string "json:\"mime\""
				ThumbInline   string "json:\"thumb_inline\""
				Width         int    "json:\"width\""
				Height        int    "json:\"height\""
				AccessHashRec string "json:\"access_hash_rec\""
			} "json:\"file_inline\""
			Text           string "json:\"text,omitempty\""
			ReplyToMessage string "json:\"reply_to_message_id,omitempty\""
		}{ObjectGUID: guid, Rnd: randNum(), FileInline: struct {
			DcID          string "json:\"dc_id\""
			FileID        string "json:\"file_id\""
			Type          string "json:\"type\""
			FileName      string "json:\"file_name\""
			Size          int    "json:\"size\""
			Mime          string "json:\"mime\""
			ThumbInline   string "json:\"thumb_inline\""
			Width         int    "json:\"width\""
			Height        int    "json:\"height\""
			AccessHashRec string "json:\"access_hash_rec\""
		}{DcID: dcId, FileID: id, Type: "Image", FileName: fileName, Size: size, Mime: "mime", ThumbInline: th, Width: width, Height: height, AccessHashRec: accessHashReq}},
		Client: clientVal,
	}
	if caption != "" && caption != " " {
		data.Input.Text = caption
	}
	if messageId != "" && messageId != " " {
		data.Input.ReplyToMessage = messageId
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newSend(auth string, dataEnc string) (map[string]string, error) {
	data := send{
		ApiVersion: apiVersion,
		Auth:       auth,
		DataEnc:    dataEnc,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error in Marshaling data:\n%v", err)
	}
	var body map[string]string
	for i := 0; i < 4; i++ {
		resp, err := http.Post(rubikaAPIList[i], jsonContentType, bytes.NewBuffer(dataJson))
		if err != nil {
			if i == 3 {
				return nil, err
			}
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Panicln("error in closing response body:", err)
			}
		}(resp.Body)
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
		break
	}
	return body, nil
}

func newEditText(text string, guid string, messageId string) (string, error) {
	data := EditText{
		Method: editMessage,
		Input: struct {
			ObjectGUID string "json:\"object_guid\""
			MessageID  string "json:\"message_id\""
			Text       string "json:\"text\""
		}{Text: text, ObjectGUID: guid, MessageID: messageId},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in Marshaling data:\n%v", err)
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func randNum() string {
	return strconv.Itoa(rand.Intn(999999))
}
func newSendMessage(text string, guid string, messageId string) (string, error) {
	data := sendMessagePayload{
		Method: sendMessage,
		Input: struct {
			ObjectGuid     string "json:\"object_guid\""
			Rnd            string "json:\"rnd\""
			Text           string "json:\"text,omitempty\""
			ReplyToMessage string "json:\"reply_to_message_id,omitempty\""
		}{ObjectGuid: guid, Rnd: randNum(), Text: text},
		Clinet: clientVal,
	}
	if messageId != "" && messageId != " " {
		data.Input.ReplyToMessage = messageId
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in Marshaling data:\n%v", err)
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newDeleteMessage(guid string, messageId ...string) (string, error) {
	var messageIds []string
	messageIds = append(messageIds, messageId...)
	data := deleteMessageStruct{
		Method: deleteMessage,
		Input: struct {
			ObjectGUID string   "json:\"object_guid\""
			MessageIds []string "json:\"message_ids\""
			Type       string   "json:\"type\""
		}{guid, messageIds, "Global"},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newChatUpdates(auth string) (string, error) {
	defaultValues := SendReqChatUpdates{
		Method: chatUpdatesMethod,
		Input: inputStr{
			State: int(time.Now().Unix()) - 200,
		},
		Client: clientStr{
			AppName:    appName,
			AppVersion: appVersion,
			Platform:   platform,
			Package:    packAge,
			LangCode:   langcode,
		},
	}
	methodValueJson, err := json.Marshal(defaultValues)
	if err != nil {
		return "", fmt.Errorf("error in marshaling default values to json. line(43)")
	}
	methodValueEncode, err := encryption.Encrypt(methodValueJson)
	if err != nil {
		return "", err
	}
	return methodValueEncode, nil
}

func newPoll(guid string, isAnonymous bool, multipleAnswers bool, question string, options ...string) (string, error) {
	var optionsList []string
	optionsList = append(optionsList, options...)
	data := createPoll{
		Method: createPollMethod,
		Input: struct {
			ObjectGUID            string   "json:\"object_guid\""
			Options               []string "json:\"options\""
			Rnd                   string   "json:\"rnd\""
			Question              string   "json:\"question\""
			Type                  string   "json:\"type\""
			IsAnonymous           bool     "json:\"is_anonymous\""
			AllowsMultipleAnswers bool     "json:\"allows_multiple_answers\""
		}{ObjectGUID: guid, Options: optionsList, Rnd: randNum(), Question: question, Type: "Regular", IsAnonymous: isAnonymous, AllowsMultipleAnswers: multipleAnswers},
		Client: struct {
			AppName    string "json:\"app_name\""
			AppVersion string "json:\"app_version\""
			Platform   string "json:\"platform\""
			Package    string "json:\"package\""
			LangCode   string "json:\"lang_code\""
		}{AppName: appName, AppVersion: appVersion, Platform: platform, Package: packAge, LangCode: langcode},
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}

func newSendInfoFile(fileName string, size int) (string, error) {
	data := infoSendFile{
		Method: sendFileMethod,
		Input: struct {
			FileName string `json:"file_name"`
			Size     int    `json:"size"`
			Mime     string `json:"mime"`
		}{FileName: fileName, Size: size, Mime: "rubika"},
		Client: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return "", err
	}
	return dataEnc, nil
}
