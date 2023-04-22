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

## تمام متد ها:
### کار با متن ها و پیام ها:
```go

```
### کار با فایل ها:
```go

```
### کار با گروه ها و کانال ها:
```go

```
### کار با کاربر ها و اکانت ها:
```go

```



### مشکلی داخل کتابخونه مشاهده کردید؟ لطفا به من اطلاع بدید
+ Rubika: @go_lang
+ Email: m4hdi2022@gmail.com



### teammate: [shayan ghosi](https://github.com/shadowcoder2020)
