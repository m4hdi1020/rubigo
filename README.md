# rubigo
**‌کتابخونه rubigo برای ساخت بات روبیکا با زبان Go**

# برای نصب:

```sh
go get -u github.com/m4hdi1020/rubigo
```

# شروع کردن:
## وارد کردن کتابخونه:
‍‍‍
```sh
import "github.com/m4hdi1020/rubigo/rubika"
```

## بعد از بدست آوردن شناسه اکانت (auth) با کد زیر بات رو ایجاد کنید
‍‍

```go
bot , err := rubika.NewBot(Auth)
```

## برای بدست آوردن آخرین پیام های دریافت شده

```go
messages , err := bot.GetMessageAll()
```

**این دستور به شماآخرین پیام های ارسال شده به همراه اطلاعات دیگه ای مثل guid نویسنده پیام، guid جایی که پیام ارسال شده، messageId پیام و... رو میده**

## برای ارسال متن

```go
err = bot.SendMessage(Guid string, Text string, ReplyToMessageID string)
```
**توجه داشته باشید که مقدار مسیج آیدی برای ریپلای کردن روی پیامه و اگر نمیخواید ریپلای کنه فقط کافیه مقدار "" رو قرار بدید**

## ارسال فایل

```go
err = bot.SendFile(Guid , FileName , data io.Reader , Caption , MessageID)
```
***اگر نمیخواید روی پیامی ریپلای کنه مقدار مسیج آیدی رو "" قرار بدید***

## تمام متد ها:
### کار با متن ها و پیام ها:
```go
GetMessageAll() ([]getChats, error)
GetMessageAllWebSocket(index int) ([]WebSocketResponse, error)
WebSocket() (*websocket.Conn, error)
GetMessagesInfoByID(guid string, messageIds ...string) (getMessageInfoData, error)
SendMessage(text string, guid string, replyToMessageID string) error
EditMessage(text string, guid string, messageId string) error
DeleteMessage(guid string, messageIds ...string) error
ForwardMessages(fromGuid string, toGuid string, messageIds ...string) error
```
### کار با فایل ها:
```go
SendFile(guid string, fileName string, data io.Reader, caption string, replyToMessageID string) error
SendImage(guid string, imageName string, data io.Reader, caption string, replyToMessageID string) error
SendFileByLink(link string, guid string, caption string, replyToMessageId string) error
SendImageByLink(link string, guid string, caption string, replyToMessageId string) error
DownloadFile(guid string, messageId string) (string, []byte, error)

```
### کار با گروه ها و کانال ها:
```go
GetGroupInfo(groupGuid string) (groupInfo, error)
GetChannelInfo(channelGuid string) (channelInfoData, error)
GetGroupAdminInfo(groupGuid string) (adminMembersData, error)
GetAllGroupMembers(groupGuid string) (allGroupMembersData, error)
GetChannelAllMembers(channelGuid string) (channelMembersData, error)
GetGroupLink(groupGuid string) (string, error)
GetChannelLink(channelGuid string) (string, error)
GetChannelAdmins(channelGuid string) (channelAdmins, error)
GetBannedGroupMembers(groupGuid string) ([]bannedList, error)
JoinGroupByLink(link string) (string, error)
LeaveGroup(guid string) error
RemoveMember(groupGuid string, memberGuid string) error
AddAdminToGroup(groupGuid, memberGuid string, adminAccessList ...string) error
// Admin Accesses: AdminChangeInfoAccess, AdminPinMessageAccess, AdminDeleteGlobalMessage, AdminBanMember , AdminSetJoinLink, AdminSetAdmin
RemoveAdminGroup(groupGuid string, memberGuid string) error
UnbanGroupMember(groupGuid, memeberGuid string) error
SetGroupAccess(groupGuid string, access ...string) error
// Group Accesses: AccessGroupAddMember, AccessGroupViewAdmins  , AccessGroupSendMessage , AccessGroupViewMembers
CreatePoll(guid string, isAnonymous bool, multipleAnswers bool, question string, options ...string) error
```
### کار با کاربر ها و اکانت ها:
```go
GetUserInfo(userGuid string) (userInfo, error)
BlockUser(userGuid string) error
UnblockUser(userGuid string) error
GetBlockedUsersList() ([]blockedUsers, error)
DeleteUserChat(userGuid, lastMessageId string) error
DeleteChatHistory(chatGuid string, lastMessageId string) error
GetInfoByUsername(username string) (infoByUsername, error)
```



### مشکلی داخل کتابخونه مشاهده کردید؟ لطفا به من اطلاع بدید
+ Rubika: @go_lang
+ Email: m4hdi2022@gmail.com



### teammate: [shayan ghosi](https://github.com/shadowcoder2020)
