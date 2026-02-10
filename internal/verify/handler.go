package verify

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"project/configs"
	"project/pkg/request"
	"slices"
	"strings"
	"sync"

	"github.com/jordan-wright/email"
)

type handler struct {
	*configs.MailConfig
}

const hashesSep = ";"

func getHashesFromFile() []string {
	s, err := os.ReadFile("hashes.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла с хэшами:", err.Error())
		return []string{}
	}

	return strings.Split(string(s), hashesSep)
}

func writeHashesToFile(hashes []string) {
	hashesString := strings.Join(hashes, hashesSep)
	hashesStringBytes := []byte(hashesString)

	if err := os.WriteFile("hashes.txt", hashesStringBytes, 0666); err != nil {
		fmt.Println("Ошибка записи файла хэшей:", err.Error())
	}
}

func (h *handler) send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[SendEmailRequest](&w, r)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)

		hash := md5.Sum([]byte(body.Email))
		hashString := hex.EncodeToString(hash[:])
		wg := sync.WaitGroup{}

		wg.Add(2)
		// читаем хэши из файла, добавляем к ним новый по емэйлу из запроса, записываем всё это в файл
		go func() {
			defer wg.Done()

			hashes := getHashesFromFile()

			if slices.Contains(hashes, hashString) {
				return
			}

			hashes = append(hashes, hashString)
			writeHashesToFile(hashes)
		}()
		// отправляем на email из запроса ссылку подтверждения почты
		go func() {
			defer wg.Done()

			e := email.NewEmail()
			e.From = h.Mail
			e.To = []string{body.Email}
			e.Subject = "Confirm your email - go to link"
			link := fmt.Sprintf(
				"<a href=\"http://localhost:3000/verify/%s\">Confirm your email</a>",
				hashString,
			)
			e.HTML = []byte(link)
			err := e.Send(
				fmt.Sprintf("%s:%s", h.SmtpHost, h.SmtpPort),
				smtp.PlainAuth("", h.Mail, h.Password, h.SmtpHost),
			)
			if err != nil {
				fmt.Println("Ошибка отправки почты:", err.Error())
				return
			}
		}()
		wg.Wait()
	}
}

func (h *handler) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		hash := r.PathValue("hash")
		hashes := getHashesFromFile()
		if slices.Contains(hashes, hash) {
			ind := slices.Index(hashes, hash)
			hashes = slices.Delete(hashes, ind, ind+1)
			writeHashesToFile(hashes)
		}
	}
}

func NewVerifyHandler(r *http.ServeMux, mailConf *configs.MailConfig) {
	h := handler{mailConf}
	r.HandleFunc("POST /send", h.send())
	r.HandleFunc("/verify/{hash}", h.verify())
}
