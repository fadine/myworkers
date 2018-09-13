# myworkers
This is a microservice build base on golang, rabbitmq, this project have only send-email service. You can pack this project by docker:
```bash
docker build -t myworkers .
```
You can change settings in conf/app.json

# Input sample (json, in rabbitmq queue)

```json
{
	"action": "EmailEngine",
	"mtype": "task",
	"data": {
		"to": ["huy010579@gmail.com"],
		"from": "mr.jun3@gmail.com",
		"body": "mail tu payspray",
		"subject": "test subject",
		"headers": "===header====",
		"footers": "----footer------",
		"template_name": "test.html"
	}
}
```