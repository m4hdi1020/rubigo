# rubigo
**‌کتابخونه rubigo برای ساخت بات روبیکا با زبان Go**

# برای نصب:

```sh
go get -u github.com/darkecho2022/rubigo
```

# شروع کردن:
## بعد از بدست آوردن شناسه اکانت (auth) با کد زیر بات رو ایجاد کنید
‍‍

```go
bot , err := rubika.NewBot("auth")
```

## برای بدست آوردن آخرین پیام های دریافت شده

```go
messages , err := bot.GetMessageAll()
```

## برای بدست آوردن آخرین پیام ارسال شده در یک گروه یا پیوی خاص

```go
message , err := bot.GetGroupMessage("group guid")
// for user
message , err := bot.GetUserMessage("user guid")
```

متد های دیگه ای هم برای دریافت پیام موجوده که روش استفاده ازشون سادست نیازی به گفتن نیست

## برای ارسال متن

```go
err = bot.SendMessage("guid" , "text")
// for replay to message
err = bot.SendMessageReply("guid" , "text" , "messageId")
```
## برای ارسال فایل

```go
err = bot.SendFile("guid" , "fileName" , data (io.Reader) , "fileType (txt , mp3 , mp4 , ...)")
```
برای ارسال عکس هم به همین شکل فقط از متد bot.SendImage استفاده کنید

متد هایی هم برای فوروارد و پین کردن پیام. حذف یک کاربر از گروه و تغییر دسترسی های گروه و همچنین دسترسی به لیست اعضا و ادمین های گروه هاو کانال ها وغیره در کتابخونه موجوده و روش استفاده هم خیلی سادست اکثرا به guid , messageId نیاز دارن

### مشکلی داخل کتابخونه مشاهده کردید؟ لطفا به من اطلاع بدید
+ rubika: @go_lang
+ email: m4hdi2022@gmail.com

***ممنون از شایان که به من توی توسعه این کتابخونه کمک زیادی کرد***
