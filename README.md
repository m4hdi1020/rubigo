# rubigo
**‌کتابخونه rubigo برای ساخت بات روبیکا با زبان Go**

# برای نصب:

```sh
go get -u github.com/darkecho2022/rubigo@v1.5.4
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
data , err := bot.DownloadFile("guid" , "MessageId")
```

متد هایی هم برای فوروارد و پین کردن پیام. حذف یک کاربر از گروه و تغییر دسترسی های گروه و همچنین دسترسی به لیست اعضا و ادمین های گروه هاو کانال ها وغیره در کتابخونه موجوده و روش استفاده هم خیلی سادست اکثرا به guid , messageId نیاز دارن

### مشکلی داخل کتابخونه مشاهده کردید؟ لطفا به من اطلاع بدید
+ Rubika: @go_lang
+ Email: m4hdi2022@gmail.com



### teammate: [shayan ghosi](https://github.com/shadowcoder2020)
