package rubika

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/m4hdi1020/rubigo/encryption"
)

func (b bot) SendMessage(text string, guid string, replyToMessageID string) error {
	if text == "" {
		return fmt.Errorf("error:Text is empty")
	}
	if guid == "" {
		return fmt.Errorf("error:Guid is empty")
	}
	dataEnc, err := newSendMessage(text, guid, replyToMessageID)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	responseText, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response_Text map[string]interface{}
	err = json.Unmarshal(responseText, &response_Text)
	if err != nil {
		log.Fatalln(err)
	}
	if response_Text["status"] != "OK" {
		return fmt.Errorf("error in sending text Message:\n%v", response_Text)
	}
	return nil
}

func (b bot) EditMessage(text string, guid string, messageId string) error {
	if text == "" {
		return fmt.Errorf("error: Text is empty")
	}
	if guid == "" {
		return fmt.Errorf("error: Guid is empty")
	}
	data, err := newEditText(b.Auth, text, guid)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := newSend(b.Auth, data)
	if err != nil {
		return err
	}
	respDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var getDecodeResp map[string]interface{}
	err = json.Unmarshal(respDecode, &getDecodeResp)
	if err != nil {
		return err
	}
	if getDecodeResp["status"] != "OK" {
		return fmt.Errorf("error in editing text message:\n%v", getDecodeResp)
	}
	return nil
}

func (b bot) DeleteMessage(guid string, messageIds ...string) error {
	if guid == "" {
		return fmt.Errorf("error: Guid is empty")
	}
	dataEnc, err := newDeleteMessage(guid, messageIds...)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	decodeBody, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var bodyJson map[string]interface{}
	err = json.Unmarshal(decodeBody, &bodyJson)
	if err != nil {
		return err
	}
	if bodyJson["status"] != "OK" {
		return fmt.Errorf("error in deleting messages:\n%v", bodyJson)
	}
	return nil
}

func (b bot) CreatePoll(guid string, isAnonymous bool, multipleAnswers bool, question string, options ...string) error {
	dataEnc, err := newPoll(guid, isAnonymous, multipleAnswers, question, options...)
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
	var responseBody map[string]interface{}
	err = json.Unmarshal(bodyDecode, &responseBody)
	if err != nil {
		return err
	}
	if responseBody["status"] != "OK" {
		return fmt.Errorf("error in sending poll:\n%v", responseBody)
	}
	return nil
}

func (b bot) SendFile(guid string, fileName string, data io.Reader, caption string, replyToMessageID string) error {
	if guid == "" {
		return fmt.Errorf("error: Guid is empty")
	}
	if fileName == "" || fileName == " "{
		fileName = "rubigo"
	}
	var buf bytes.Buffer
	i, err := io.Copy(&buf, data)
	if err != nil {
		return err
	}
	fileBytes := buf.Bytes()
	size := int(i)
	sizeStr := strconv.Itoa(size)
	id, dcId, hashAccess, url, err := getInfoSendFile(fileName, size, b.Auth)
	if err != nil {
		return err
	}
	if size <= 131072 {
		hashReq, err := sendPartFile(url, id, dcId, fileName, hashAccess, b.Auth, sizeStr, "1", "1", sizeStr, bytes.NewBuffer(fileBytes))
		if err != nil {
			return err
		}
		dataEnc, err := newSendFile(caption, guid, dcId, id, fileName, size, hashReq.Data.AccessHashRec, replyToMessageID)
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
		var responseBody map[string]interface{}
		err = json.Unmarshal(bodyDecode, &responseBody)
		if err != nil {
			return err
		}
		if responseBody["status"] != "OK" {
			return fmt.Errorf("error in sending file:\n%v", err)
		}
		return nil
	} else {
		totalPart := int(math.Ceil(float64(size) / float64(131072)))
		s := 0
		e := 131072
		for i := 1; i <= totalPart; i++ {
			if i == totalPart {
				hashReq, err := sendPartFile(url, id, dcId, fileName, hashAccess, b.Auth, sizeStr, strconv.Itoa(i), strconv.Itoa(totalPart), strconv.Itoa(len(fileBytes[s:])), bytes.NewBuffer(fileBytes[s:]))
				if err != nil {
					return err
				}
				dataEnc, err := newSendFile(caption, guid, dcId, id, fileName, size, hashReq.Data.AccessHashRec, replyToMessageID)
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
				var responseBody map[string]interface{}
				err = json.Unmarshal(bodyDecode, &responseBody)
				if err != nil {
					return err
				}
				if responseBody["status"] != "OK" {
					return fmt.Errorf("error in sending file:\n%v", err)
				}
				return nil
			} else {
				_, err := sendPartFile(url, id, dcId, fileName, hashAccess, b.Auth, sizeStr, strconv.Itoa(i), strconv.Itoa(totalPart), strconv.Itoa(len(fileBytes[s:e])), bytes.NewBuffer(fileBytes[s:e]))
				if err != nil {
					return err
				}
				s = e
				e += 131072
			}
		}
	}
	return nil
}
func getInfoSendFile(fileName string, fileSize int, auth string) (string, string, string, string, error) {
	dataEnc, err := newSendInfoFile(fileName, fileSize)
	if err != nil {
		return "", "", "", "", err
	}
	body, err := newSend(auth, dataEnc)
	if err != nil {
		return "", "", "", "", err
	}
	bodydecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return "", "", "", "", err
	}
	var bodyJson getRespSendFile
	err = json.Unmarshal(bodydecode, &bodyJson)
	if err != nil {
		return "", "", "", "", err
	}
	return bodyJson.Data.ID, bodyJson.Data.DcID, bodyJson.Data.AccessHashSend, bodyJson.Data.UploadURL, nil
}
func sendPartFile(url string, id string, dcId string, fileName string, hashAccess string, auth string, size string, partNumber string, totalPart string, chunkSize string, body io.Reader) (getHashReq, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("access-hash-send", hashAccess)
	req.Header.Set("auth", auth)
	req.Header.Set("chunk-size", chunkSize)
	req.Header.Set("file-id", id)
	req.Header.Set("total-part", totalPart)
	req.Header.Set("part-number", partNumber)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var hashReq getHashReq
	err = json.NewDecoder(resp.Body).Decode(&hashReq)
	if err != nil {
		log.Fatalln(err)
	}
	return hashReq, nil
}
func sendImageData(url string, auth string, id string, hashAccess string, totalPart, partNumber, chunkSize string, body io.Reader) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("access-hash-send", hashAccess)
	req.Header.Set("file-id", id)
	req.Header.Set("auth", auth)
	req.Header.Set("chunk-size", chunkSize)
	req.Header.Set("total-part", totalPart)
	req.Header.Set("part-number", partNumber)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if partNumber == totalPart {
		var response getHashReq
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return "", err
		}
		return response.Data.AccessHashRec, nil
	}
	return "", nil
}

func (b bot) JoinGroupByLink(link string) (string, error) {
	if link == "" {
		return "", fmt.Errorf("error: link invalid")
	}
	hashLink := strings.Replace(link, "https://rubika.ir/joing/", "", 1)
	dataEncJoinGroup, err := newJoinGroup(hashLink)
	if err != nil {
		return "", err
	}
	body, err := newSend(b.Auth, dataEncJoinGroup)
	if err != nil {
		return "", err
	}
	bodyDecode, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return "", err
	}
	var join joinGroupReponse
	err = json.Unmarshal(bodyDecode, &join)
	if err != nil {
		return "", err
	}
	if join.Status != "OK" {
		return "", fmt.Errorf("error joinig to group:\nReponse: %+v", bodyDecode)
	}
	return join.Data.Group.GroupGUID, nil
}

func (b bot) LeaveGroup(guid string) error {
	if guid == "" {
		return fmt.Errorf("error: Guid is empty")
	}
	dataEnc, err := newLeaveGroup(guid)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyJson, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var reponse map[string]interface{}

	err = json.Unmarshal(bodyJson, &reponse)
	if err != nil {
		return err
	}
	if reponse["status"] != "OK" {
		return fmt.Errorf("error leaving group:\n%v", reponse)
	}
	return nil
}

func (b bot) RemoveMember(groupGuid string, memberGuid string) error {
	if groupGuid == "" {
		return fmt.Errorf("error: GroupGuid is empty")
	}
	if memberGuid == "" {
		return fmt.Errorf("error: MemberGuid is empty")
	}
	dataEnc, err := newRemoveMember(groupGuid, memberGuid, "Set")
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
		return fmt.Errorf("error: %v", response)
	}
	return nil
}

func (b bot) PinMessage(groupGuid, messageId string) error {
	if groupGuid == "" {
		return fmt.Errorf("error: GroupGuid is empty")
	}
	if messageId == "" {
		return fmt.Errorf("error: MessageId is empty")
	}
	dataEnc, err := newPinMessage(groupGuid, messageId)
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
		return fmt.Errorf("error pinning message:\nResponse: %v", response)
	}
	return nil
}

func (b bot) ForwardMessages(fromGuid string, toGuid string, messageIds ...string) error {
	if fromGuid == "" {
		return fmt.Errorf("error: FromGuid is empty")
	}
	if toGuid == "" {
		return fmt.Errorf("error: ToGuid is empty")
	}
	var messageIdList []string
	messageIdList = append(messageIdList, messageIds...)
	dataEnc, err := newForwardMessage(fromGuid, toGuid, messageIdList)
	if err != nil {
		return err
	}
	body, err := newSend(b.Auth, dataEnc)
	if err != nil {
		return err
	}
	bodyDeocde, err := encryption.Decrypt(body["data_enc"])
	if err != nil {
		return err
	}
	var response map[string]interface{}
	err = json.Unmarshal(bodyDeocde, &response)
	if err != nil {
		return err
	}
	if response["status"] != "OK" {
		return fmt.Errorf("error forwarding Messages >>>\nResponse: %v", response)
	}
	return nil
}

func (b bot) AddAdminToGroup(groupGuid, memberGuid string, adminAccessList ...string) error {
	if groupGuid[0:1] != "g" {
		return fmt.Errorf("error: your guid is invalid")
	}
	if memberGuid[0:1] != "u" {
		return fmt.Errorf("error: your guid is invalid")
	}
	var accessList []string
	accessList = append(accessList, adminAccessList...)

	dataEnc, err := newAddGroupAdmin(groupGuid, memberGuid, accessList)
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
		return fmt.Errorf("error adding admin to the group >>>\nGroup Guid: %v\nAdmin Guid: %v\nServer Response: %v", groupGuid, memberGuid, response)
	}
	return nil
}

func (b bot) RemoveAdminGroup(groupGuid string, memberGuid string) error {
	if groupGuid[0:1] != "g" {
		return fmt.Errorf("error: your guid is invalid")
	}
	if memberGuid[0:1] != "u" {
		return fmt.Errorf("error: your guid is invalid")
	}
	dataEnc, err := newRemoveGroupAdmin(groupGuid, memberGuid)
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
		return fmt.Errorf("error removing admin from group\nGroupGuid: %v\nAdmin Guid: %v\nServer Response: %v", groupGuid, memberGuid, response)
	}
	return nil
}

func (b bot) SetGroupAccess(groupGuid string, access ...string) error {
	var accessList []string
	accessList = append(accessList, access...)

	dataEnc, err := newSetGroupAccess(groupGuid, accessList)
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
		return fmt.Errorf("error Seting Group Access >>>\nGroup Guid: %v\nServer Response: %v", groupGuid, response)
	}
	return nil
}

func (b bot) SendImage(guid string, imageName string, data io.Reader, caption string, replyToMessageID string) error {
	if imageName == "" || imageName == " "{
		imageName = "rubigo"
	}
	var buf bytes.Buffer
	i, err := io.Copy(&buf, data)
	if err != nil {
		return err
	}
	imageBytes := buf.Bytes()
	reader := bufio.NewReader(&buf)
	imageInfo, _, err := image.DecodeConfig(reader)
	if err != nil {
		return err
	}
	size := int(i)
	id, dcId, hashAccess, url, err := getInfoSendFile(imageName, size, b.Auth)
	if err != nil {
		return err
	}
	if size <= 131072 {
		hashReq, err := sendImageData(url, b.Auth, id, hashAccess, "1", "1", strconv.Itoa(size), bytes.NewBuffer(imageBytes))
		if err != nil {
			return err
		}
		dataEnc, err := newSendImage(guid, caption, dcId, id, imageName, size, imageInfo.Width, imageInfo.Height, hashReq, replyToMessageID)
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
			return fmt.Errorf("error sending image >>>\nServer Response:\n%v", response)
		}
		return nil
	} else {
		totalPart := int(math.Ceil(float64(size) / float64(131072)))
		totalPartStr := strconv.Itoa(totalPart)
		s := 0
		e := 131072
		for i := 1; i <= totalPart; i++ {
			if i == totalPart {
				hashReq, err := sendImageData(url, b.Auth, id, hashAccess, totalPartStr, strconv.Itoa(i), strconv.Itoa(len(imageBytes[s:])), bytes.NewBuffer(imageBytes[s:]))
				if err != nil {
					return err
				}
				dataEnc, err := newSendImage(guid, caption, dcId, id, imageName, size, imageInfo.Width, imageInfo.Height, hashReq, replyToMessageID)
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
					return fmt.Errorf("error sending image >>>\nServer Response:\n%v", response)
				}
				return nil
			} else {
				_, err := sendImageData(url, b.Auth, id, hashAccess, totalPartStr, strconv.Itoa(i), strconv.Itoa(len(imageBytes[s:e])), bytes.NewBuffer(imageBytes[s:e]))
				if err != nil {
					return err
				}
				s = e
				e += 131072
			}
		}
	}
	return nil
}

func (b bot) SendFileByLink(link string, guid string, caption string, replyToMessageId string) error {
	u, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("error: your link is invalid\nError: %s", err.Error())
	}
	fileName := path.Base(u.Path)
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("error http request:\nError Message: %s", err.Error())
	}
	defer resp.Body.Close()
	err = b.SendFile(guid, fileName, resp.Body, caption, replyToMessageId)
	if err != nil {
		return err
	}
	return nil
}

func (b bot) SendImageByLink(link string, guid string, caption string, replyToMessageId string) error {
	u, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("error: your link is invalid\nError: %s", err.Error())
	}
	imageName := path.Base(u.Path)
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("error http request:\nError Message: %s", err.Error())
	}
	defer resp.Body.Close()
	err = b.SendImage(guid, imageName, resp.Body, caption, replyToMessageId)
	if err != nil {
		return err
	}
	return nil
}

func (b bot) UnbanGroupMember(groupGuid, memeberGuid string) error {
	dataEnc, err := newRemoveMember(groupGuid, memeberGuid, "Unset")
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
		return fmt.Errorf("error unbannig group member --->\nServer Response:\n%v", response)
	}
	return nil
}
