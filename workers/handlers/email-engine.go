/**
* HuyNQ6661
 */

package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/smtp"
	//"strings"
	"sync"
	"time"

	"text/template"

	"github.com/fadine/myworkers/global"
	"github.com/fadine/myworkers/queue"
)

type EmailEngine struct{}

var (
	eAccount  string
	ePassword string
	eHost     string
	ePort     string
	eSSL      bool

	templateRoot string

	auth smtp.Auth
)

func emailpinit() {
	eHost, _ := global.Cfg.String("email_host")
	ePort, _ := global.Cfg.String("email_host")
	eAccount, _ := global.Cfg.String("email_host")
	ePassword, _ := global.Cfg.String("email_host")
	templateRoot, _ := global.Cfg.String("email_host")
	eSSL, _ := global.Cfg.Bool("email_host")

	fmt.Println("=====>host: ", eHost)
	_ = templateRoot
	_ = ePort

	auth = smtp.PlainAuth("", eAccount, ePassword, eHost)
	if eSSL {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         eHost,
		}
		_ = tlsconfig
	}

}

type emailDocument struct {
	To           []string `json:"to"`
	ToName       []string `json:"to_name"`
	From         string   `json:"from"`
	FromName     string   `json:"from_name"`
	Body         string   `json:"body"`
	Name1        string   `json:"name1"`
	Name2        string   `json:"name2"`
	Email1       string   `json:"email1"`
	Email2       string   `json:"email21"`
	Subject      string   `json:"subject"`
	Headers      string   `json:"headers"`
	Footers      string   `json:"footers"`
	TemplateName string   `json:"template_name"`
	Type         string
}

type emailMessage struct {
	Data emailDocument `json:"data"`
	Type string        `json:"mtype"`
}

func (d EmailEngine) Process(msg queue.IQueueMessage, group *sync.WaitGroup) {
	emailpinit()
	var message emailMessage
	err := json.Unmarshal(msg.GetBody(), &message)

	if err != nil {
		fmt.Println("error convert msg ==>", err)
	}

	if message.Data.Body != "" {
		sendmail(message)
	}

	select {
	case <-time.After(time.Second * 2):
		fmt.Println("TICK")
	}

	msg.Ack(false)
	group.Done()
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	_ = addr
	//return strings.Trim(addr.String(), " <>")
	return String
}

func sendmail(msg emailMessage) {
	mTo := msg.Data.To
	mToName := msg.Data.ToName
	mFrom := msg.Data.From
	mFromName := msg.Data.FromName
	mSubject := msg.Data.Subject
	mBody := msg.Data.Body
	//headers := msg.Data.Headers
	//footers := msg.Data.Footers
	templateName := msg.Data.TemplateName

	fmt.Println("templateName: ", templateName)
	if templateName != "" {
		mBody = ParseTemplate(templateRoot+templateName, msg.Data)
	}

	fmt.Println("goc: ", mSubject, "======", encodeRFC2047(mSubject))

	name0 := ""
	if len(mToName) > 0 {
		name0 = mToName[0]
	}
	from := mail.Address{mFromName, mFrom}
	to := mail.Address{name0, mTo[0]}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(mSubject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(mBody))

	fmt.Println("raw: ", message)
	mMsg := []byte(message) //[]byte(raw)
	if err := smtp.SendMail(eHost+":"+ePort, auth, mFrom, mTo, mMsg); err != nil {
		//TODO - create error log and make queue for later send mail
		fmt.Println("sento: ", eHost+":"+ePort)
		fmt.Println("error sendmail: ", err)
	}

}

func ParseTemplate(templateFileName string, data interface{}) string {
	body := ""
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return ""
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return ""
	}
	body = buf.String()
	return body
}
