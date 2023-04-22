package rubika

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/m4hdi1020/rubigo/encryption"
)

func (b bot) GetMessageAll() ([]getChats, error) {
	values, err := newChatUpdates(b.Auth)
	if err != nil {
		return nil, err
	}
	body, err := newSend(b.Auth, values)
	if err != nil {
		return nil, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return nil, err
	}
	var response getResponseChatUpdates
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf("error getting data >>>\nServer Rerponse:\n%+v", response)
	}
	return response.Data.Chats, nil
}

func (b bot) GetPvMessageAll() ([]getChats, error) {
	messages, err := b.GetMessageAll()
	if err != nil {
		return nil, err
	}
	var messageList []getChats
	for i := range messages {
		if messages[i].AbsObject.ObjectGuid[0:1] == "u" {
			messageList = append(messageList, messages[i])
		}
	}
	return messageList, nil
}

func (b bot) GetGroupMessageAll() ([]getChats, error) {
	messages, err := b.GetMessageAll()
	if err != nil {
		return nil, err
	}
	var messageList []getChats
	for i := range messages {
		if messages[i].AbsObject.ObjectGuid[0:1] == "g" {

			messageList = append(messageList, messages[i])
		}
	}
	return messageList, nil
}

func (b bot) GetMessageAllWebSocket(index int) ([]WebSocketResponse, error) {
	if index <= 0 {
		return nil, fmt.Errorf("error: please enter a number greater than zero")
	}
	resp, err := newWebSocket(b.Auth, index)
	if err != nil {
		log.Fatalln(err)
	}
	return resp, nil
}

func (b bot) WebSocket() (*websocket.Conn, error) {
	data := webSocketData{
		Method: webSocketMethod,
		Input:  struct{}{},
		Clinet: clientVal,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataEnc, err := encryption.Encrypt(dataJson)
	if err != nil {
		return nil, err
	}
	webSocketURL := [4]string{"wss://jsocket2.iranlms.ir:80/", "wss://jsocket2.iranlms.ir:80/", "wss://nsocket6.iranlms.ir:80/", "wss://msocket1.iranlms.ir:80/"}
	send := send{
		ApiVersion: "5",
		Auth:       b.Auth,
		DataEnc:    dataEnc,
	}
	if err != nil {
		log.Fatalln(err)
	}
	// rand.Seed(time.Now().Unix())
	conn, _, err := websocket.DefaultDialer.Dial(webSocketURL[rand.Intn(3-0+1)+0], nil)
	if err != nil {
		log.Fatalln(err)
	}
	_, _, err = conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	err = conn.WriteJSON(send)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (b bot) GetUserInfo(userGuid string) (userInfo, error) {
	if userGuid[0:1] != "u" {
		return userInfo{}, fmt.Errorf("error: your GUID is invalid")
	}
	dataEnc, err := newUserInfo(userGuid)
	if err != nil {
		return userInfo{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return userInfo{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return userInfo{}, err
	}
	var response userInfoData

	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return userInfo{}, err
	}
	if response.Status != "OK" {
		return userInfo{}, fmt.Errorf("error getting user info >>>\nUser Guid: %v\nServer Response:\n%+v", userGuid, response)
	}
	return response.Data, nil
}

func (b bot) BlockUser(userGuid string) error {
	if userGuid[0:1] != "u" {
		return fmt.Errorf("error: your GUID is invalid")
	}
	dataEnc, err := newBlockUser(userGuid, block)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response map[string]interface{}
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return err
	}
	if response["status"] != "OK" {
		return fmt.Errorf("error blocking user >>>\nUser Guid: %v\nServer Response:\n%+v", userGuid, response)
	}
	return nil
}

func (b bot) UnblockUser(userGuid string) error {
	if userGuid[0:1] != "u" {
		return fmt.Errorf("error: your GUID is invalid")
	}
	dataEnc, err := newBlockUser(userGuid, Unblock)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response map[string]interface{}
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return err
	}
	if response["status"] != "OK" {
		return fmt.Errorf("error Unblocking user >>>\nUser Guid: %v\nServer Response:\n%+v", userGuid, response)
	}
	return nil
}

func (b bot) DeleteUserChat(userGuid, lastMessageId string) error {
	dataEnc, err := newDeleteUserChat(userGuid, lastMessageId)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response map[string]interface{}

	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return err
	}
	fmt.Println(response)
	if response["status"] != "OK" {
		return fmt.Errorf("error delete user chat >>> \nResponse: %v", response)
	}
	return nil
}

func (b bot) GetGroupInfo(groupGuid string) (groupInfo, error) {
	if groupGuid[0:1] != "g" {
		return groupInfo{}, fmt.Errorf("error: your GUID is invalid")
	}
	dataEnc, err := newGroupInfo(groupGuid)
	if err != nil {
		return groupInfo{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return groupInfo{}, err
	}

	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return groupInfo{}, err
	}
	var response getGroupInfo
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return groupInfo{}, err
	}
	return response.Data, nil
}

func (b bot) DeleteChatHistory(chatGuid string, lastMessageId string) error {
	dataEnc, err := newDeleteChatHistory(chatGuid, lastMessageId)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response map[string]interface{}

	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return err
	}
	if response["status"] != "OK" {
		return fmt.Errorf("error deleting chat history:\nResponse: %v", response)
	}
	return nil
}

func (b bot) GetInfoByUsername(username string) (infoByUsername, error) {
	if found := strings.Contains(username, "@"); found {
		username = strings.Replace(username, "@", "", 1)
	}
	dataEnc, err := newGetInfoById(username)
	if err != nil {
		return infoByUsername{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return infoByUsername{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return infoByUsername{}, err
	}
	var response getResponseInfoById
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return infoByUsername{}, err
	}
	return response.Data, nil
}

func (b bot) GetChannelInfo(channelGuid string) (channelInfoData, error) {
	if channelGuid[0:1] != "c" {
		return channelInfoData{}, fmt.Errorf("error: your GUID is invalid")
	}
	dataEnc, err := newChannelInfo(channelGuid)
	if err != nil {
		return channelInfoData{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return channelInfoData{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return channelInfoData{}, err
	}
	var response channelInfo
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return channelInfoData{}, err
	}
	return response.Data, nil
}

func (b bot) GetGroupAdminInfo(groupGuid string) (adminMembersData, error) {
	if groupGuid[0:1] != "g" {
		return adminMembersData{}, fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newGetGroupAdminMembers(groupGuid)
	if err != nil {
		return adminMembersData{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return adminMembersData{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return adminMembersData{}, err
	}
	var response adminMembers
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return adminMembersData{}, err
	}
	if response.Status != "OK" {
		return adminMembersData{}, fmt.Errorf("error getting group admin Info >>>\nGroup Guid: %v\nServer Response:\n%+v", groupGuid, response)
	}
	return response.Data, nil
}

func (b bot) GetAllGroupMembers(groupGuid string) (allGroupMembersData, error) {
	if groupGuid[0:1] != "g" {
		return allGroupMembersData{}, fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newGetAllGroupMembers(groupGuid)
	if err != nil {
		return allGroupMembersData{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return allGroupMembersData{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return allGroupMembersData{}, err
	}
	var response getAllGroupMembers
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return allGroupMembersData{}, err
	}
	if response.Status != "OK" {
		return allGroupMembersData{}, fmt.Errorf("error getting all group Members >>>\nGroup Guid: %v\nServer Response:\n%+v", groupGuid, response)
	}
	return response.Data, nil
}

func (b bot) GetChannelAllMembers(channelGuid string) (channelMembersData, error) {
	if channelGuid[0:1] != "c" {
		return channelMembersData{}, fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newGetChannelAllMembers(channelGuid)
	if err != nil {
		return channelMembersData{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return channelMembersData{}, err
	}
	bodyDeocde, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return channelMembersData{}, err
	}
	var response getChannelAllMembers

	err = json.Unmarshal(bodyDeocde, &response)
	if err != nil {
		return channelMembersData{}, err
	}
	if response.Status != "OK" {
		return channelMembersData{}, fmt.Errorf("error getting channel all Members >>>\nChannel Guid: %v\nServer Response:\n%+v", channelGuid, response)
	}
	return response.Data, nil
}

func (b bot) GetGroupLink(groupGuid string) (string, error) {
	if groupGuid[0:1] != "g" {
		return "", fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newGetGroupLink(groupGuid)
	if err != nil {
		return "", err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return "", err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return "", err
	}
	var response getGroupLink
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return "", err
	}
	if response.Status != "OK" {
		return "", fmt.Errorf("error Getting Group Link >>>\nGroup Guid: %v\nServer Response:\n%v", groupGuid, response)
	}
	return response.Data.JoinLink, nil
}

func (b bot) GetChannelLink(channelGuid string) (string, error) {
	if channelGuid[0:1] != "c" {
		return "", fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newChannelInfo(channelGuid)
	if err != nil {
		return "", err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return "", err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return "", err
	}
	var response channelInfo
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return "", err
	}
	if response.Status != "OK" {
		return "", fmt.Errorf("error getting channel link >>>\nChannel Guid: %v\nServer Response:\n%+v", channelGuid, response)
	}
	return response.Data.Channel.Username, nil
}

func (b bot) GetChannelAdmins(channelGuid string) (channelAdmins, error) {
	if channelGuid[0:1] != "c" {
		return channelAdmins{}, fmt.Errorf("error: your auth is invalid")
	}
	dataEnc, err := newGetChannelAdmins(channelGuid)
	if err != nil {
		return channelAdmins{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return channelAdmins{}, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return channelAdmins{}, err
	}
	var response channelAdminsInfo
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return channelAdmins{}, err
	}
	if response.Status != "OK" {
		return channelAdmins{}, fmt.Errorf("error getting channel admins >>>\nChannel Guid: %v\nServer Response:\n%+v", channelGuid, response)
	}
	return response.Data, nil
}

func (b bot) GetMessagesInfoByID(guid string, messageIds ...string) (getMessageInfoData, error) {
	dataEnc, err := newGetMessageInfoByID(guid, messageIds...)
	if err != nil {
		return getMessageInfoData{}, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return getMessageInfoData{}, err
	}

	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return getMessageInfoData{}, err
	}
	var response getMessageInfoByID
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return getMessageInfoData{}, err
	}
	if response.Status != "OK" {
		return getMessageInfoData{}, fmt.Errorf("error getting message Info by ID --->\nServer Response:\n%+v", response)
	}
	return response.Data, nil

}

func (b bot) DownloadFile(guid string, messageId string) (string, []byte, error) {
	info, err := b.GetMessagesInfoByID(guid, messageId)
	if err != nil {
		return "", nil, err
	}
	if info.Messages[0].FileInline.FileID == 0 {
		return "", nil, fmt.Errorf("error: your message is not a file")
	}
	data, err := downloader(b.Auth, strconv.Itoa(int(info.Messages[0].FileInline.FileID)), strconv.Itoa(info.Messages[0].FileInline.DcID), info.Messages[0].FileInline.AccessHashRec, info.Messages[0].FileInline.Size)
	if err != nil {
		return "", nil, err
	}
	return info.Messages[0].FileInline.FileName, data, nil
}

func downloader(auth string, fileId string, dcID string, accessHash string, size int) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://messenger%s.iranlms.ir/GetFile.ashx", dcID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("access-hash-rec", accessHash)
	req.Header.Set("auth", auth)
	req.Header.Set("file-id", fileId)

	totalPart := int(math.Ceil(float64(size) / float64(131072)))
	s := 0
	e := 131072
	data := make([]byte, 0, size)
	for i := 1; i <= totalPart; i++ {
		if i == totalPart {
			req.Header.Set("start-index", strconv.Itoa(s))
			req.Header.Set("last-index", strconv.Itoa(size))
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			data = append(data, body...)
			resp.Body.Close()
		} else {
			req.Header.Set("start-index", strconv.Itoa(s))
			req.Header.Set("last-index", strconv.Itoa(e))
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			data = append(data, body...)
			resp.Body.Close()
			s = e + 1
			e += 131072
		}
	}
	return data, nil
}

func (b bot) GetBlockedUsersList() ([]blockedUsers, error) {
	dataEnc, err := newGetBlockedUsersPayload()
	if err != nil {
		return nil, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return nil, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return nil, err
	}
	var response getBlockedUsersResponse
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf("error getting blocked users list --->\nServer Response:\n%+v", response)
	}
	return response.Data.AbsUsers, nil
}

func (b bot) GetBannedGroupMembers(groupGuid string) ([]bannedList, error) {
	dataEnc, err := newBannedListReq(groupGuid)
	if err != nil {
		return nil, err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return nil, err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return nil, err
	}
	var response bannedGroupMembersResp
	err = json.Unmarshal(bodyDecode, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf("error getting banned group members --->\nServer Response:\n%+v", response)
	}
	return response.Data.InChatMembers, nil
}
