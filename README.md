# rubigo
**‌کتابخونه rubigo برای ساخت بات روبیکا با زبان Go**

# برای نصب:

```sh
go get -u github.com/m4hdi1020/rubigo@v1.5.5
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
bot , err := rubika.NewBot("auth")
```

## برای بدست آوردن آخرین پیام های دریافت شده

```go
messages , err := bot.GetMessageAll()
```

**این دستور به شماآخرین پیام های ارسال شده به همراه اطلاعات دیگه ای مثل guid نویسنده پیام، guid جایی که پیام ارسال شده، messageId پیام و... رو میده**

## برای ارسال متن

```go
err = bot.SendMessage("guid" , "text" , "messageId")
```
**توجه کنید که مقدار مسیج آیدی برای ریپلای کردن روی پیامه و اگر نمیخواید ریپلای کنه فقط کافیه مقدار "" رو قرار بدید**

## ارسال فایل

```go
err = bot.SendFile("guid" , "fileName" , data (io.Reader) , "caption" , "MessageId")
```
***اگر نمیخواید روی پیامی ریپلای کنه مقدار مسیج آیدی رو "" قرار بدید***
برای ارسال عکس هم به همین شکل فقط از متد bot.SendImage استفاده کنید

## ارسال فایل و عکس با لینک
```go
err = bot.SendFileByLink("link" , "guid" , data , "caption" , "MessageId")

err = bot.SendImageByLink("link" , "guid" , data , "caption" , "MessageId")
```

## دانلود فایل از روبیکا
```go
fileName , data , err := bot.DownloadFile("guid" , "MessageId")
```

## بقیه متد ها:
```go
bot.EditMessage(newText , guid , messageID)
bot.DeleteMessage(guid , messageID1 , messageID2 , messageID3 , ...)
bot.CreatePoll(guid string, isAnonymous bool, multipleAnswers bool, question string, options ...string)
bot.GetUserInfo("User Guid")
bot.BlockUser("User Guid")
bot.UnblockUser("User Guid")
bot.DeleteUserChat("User Guid" , "last Message Id")
bot.GetGroupInfo("Group Guid")
bot.DeleteChatHistory("Chat Guid" , "Last Message Id")
bot.GetInfoByUsername("Username")
bot.GetChannelInfo("Channel Guid")
bot.GetGroupAdminInfo("Group Guid")
bot.GetAllGroupMembers("Group Guid")
bot.GetChannelAllMembers("Channel Guid")
bot.GetGroupLink("Group Guid")
bot.GetChannelLink("Channel Guid")
bot.GetChannelAdmins("Channel Guid")
bot.GetMessagesInfoByID("Chat Guid" , "Message Id")
bot.EditMessage("New Text" , "Chat Guid" , "Message ID")
bot.DeleteMessage("Chat Guid" , "Message ID 1" , "Message ID 2" , "Message ID 3" , "...")
bot.CreatePoll("Chat Guid" , isAnonymous , multipleAnswers , "question" , "option1" , "option2" , "option3" , "...")
bot.JoinGroupByLink("Group Link")
bot.LeaveGroup("Group Guid")
bot.RemoveMember("Group Guid" , "Member Guid")
bot.PinMessage("Group Guid" , "Message ID")
bot.ForwardMessages("from Guid" , "to Guid" , "message ID 1" , "Message ID 2" , "Message ID 3" , "...")
bot.AddAdminToGroup("Group Guid" , "Member Guid" , AdminAccessList...)
// Admin Access Option => AdminChangeInfoAccess , AdminPinMessageAccess , AdminDeleteGlobalMessage , AdminBanMember , AdminSetJoinLink , AdminSetAdmin , AdminSetMemberAccess
bot.RemoveAdminGroup("Group Guid" , "Admin Guid")
bot.SetGroupAccess("Group Guid" , GroupAccess...)
// Group Acess Option => AccessGroupAddMember , AccessGroupViewAdmins , AccessGroupSendMessage , AccessGroupViewMembers
```

### مشکلی داخل کتابخونه مشاهده کردید؟ لطفا به من اطلاع بدید
+ Rubika: @go_lang
+ Email: m4hdi2022@gmail.com



### teammate: [shayan ghosi](https://github.com/shadowcoder2020)
