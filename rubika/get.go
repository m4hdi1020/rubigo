package rubika

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/darkecho2022/rubigo/encryption"
	"github.com/gorilla/websocket"
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
		Clinet: struct {
			AppName    string "json:\"app_name\""
			AppVersion string "json:\"app_version\""
			Platform   string "json:\"platform\""
			Package    string "json:\"package\""
			LangCode   string "json:\"lang_code\""
		}{AppName: appName, AppVersion: apiVersion, Platform: platform, Package: packAge, LangCode: langcode},
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataEnc, err := encryption.Encrypt(key, dataJson)
	if err != nil {
		return nil, err
	}
	webSocketURL := [4]string{"wss://jsocket2.iranlms.ir:80/", "wss://jsocket2.iranlms.ir:80/", "wss://nsocket6.iranlms.ir:80/", "wss://msocket1.iranlms.ir:80/"}
	send := Send{
		ApiVersion: "5",
		Auth:       b.Auth,
		DataEnc:    dataEnc,
	}
	if err != nil {
		log.Fatalln(err)
	}
	rand.Seed(time.Now().Unix())
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

func (b bot) GetUserMessage(userGuid string) (getChats, error) {
	if userGuid[0:1] != "u" {
		return getChats{}, fmt.Errorf("error: your auth is invalid")
	}
	messages, err := b.GetMessageAll()
	if err != nil {
		return getChats{}, err
	}
	for i := range messages {
		if messages[i].AbsObject.ObjectGuid[0:1] == "u" && messages[i].LastMessage.AuthorObjectGuid == userGuid {
			return messages[i], nil
		}
	}
	return getChats{}, nil
}

func (b bot) GetGroupMessage(groupGuid string) (getChats, error) {
	if groupGuid[0:1] != "g" {
		return getChats{}, fmt.Errorf("error: your auth is invalid")
	}
	messages, err := b.GetMessageAll()
	if err != nil {
		log.Fatalln(err)
	}
	for i := range messages {
		if messages[i].AbsObject.ObjectGuid[0:1] == "g" && messages[i].AbsObject.ObjectGuid == groupGuid {
			return messages[i], nil
		}
	}
	return getChats{}, nil
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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

	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDeocde, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
	bodyDecode, err := encryption.Decrypt(key, body["data_enc"])
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
