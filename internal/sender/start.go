package sender

import (
	"errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Sender) Start() {
	var sendedSingleMessagesCount int

	ticker := time.NewTicker(time.Second / 25)
	defer ticker.Stop()

	for {
		<-ticker.C // В любом случае ждем тика

		var (
			err       error                   // Переменная для ошибки отправки
			missedMsg *tgbotapi.MessageConfig // Переменная для сообщения, которое не отправилосб из-за ошибки с кодом 429
		)

		if sendedSingleMessagesCount < 2 {
			// Сообщения для одного пользователя и для множества пользователей отправляются в соотношении 2:1

			select {

			case msg := <-s.singleMessages:
				// Основной сценарий, отправляем сообщение для одного пользователя
				if _, err = s.bot.Send(msg); err != nil {
					missedMsg = msg // Если произошла ошибка, пишем сообщение в переменную для пропущенного сообщения
				}
				sendedSingleMessagesCount++

			default:
				// Если надо отправить одному пользователю, а канал пуст, проверяем канал с сообщениями для множества пользователей, чтобы не простаивал тик
				select {

				case msg := <-s.massMessages:
					if _, err = s.bot.Send(msg); err != nil {
						missedMsg = msg // Если произошла ошибка, пишем сообщение в переменную для пропущенного сообщения
					}

				default:
					// Если оба пусты, ничего не делаем
				}
			}

		} else {
			// 2 сообщения для одного пользователя уже отправлены

			select {

			case msg := <-s.massMessages:
				// Основной сценарий, отправляем одно сообщение для множества пользователей
				if _, err = s.bot.Send(msg); err != nil {
					missedMsg = msg // Если произошла ошибка, пишем сообщение в переменную для пропущенного сообщения
				}
				sendedSingleMessagesCount = 0 // Счетчик на 0, чтобы следующие 2 отправились одиночные

			default:
				// Если канал с сообщениями для множества пользователей пуст, проверяем канал с одиночными, чтобы не простаивать тик
				select {

				case msg := <-s.singleMessages:
					if _, err = s.bot.Send(msg); err != nil {
						missedMsg = msg // Если произошла ошибка, пишем сообщение в переменную для пропущенного сообщения
					}
					// С счетчиком ничего не делаем, чтобы следующее сообщение пошло по ветке для множества пользователей

				default:
					// Оба пусты, ничего не делаем
				}
			}
		}

		if err != nil {
			// При получении ошибки проверяем её, если это 429 код (слишком много сообщений), уводим горутину в сон на время, указанное в RetryAfter
			var tgErr tgbotapi.Error

			if errors.As(err, &tgErr) && tgErr.Code == 429 {

				log.Printf("429 code. waiting for %d seconds", tgErr.RetryAfter)
				time.Sleep(time.Duration(tgErr.RetryAfter) * time.Second)
				s.bot.Send(missedMsg) // Отправка того самого пропущенного сообщения

			}
		}
	}
}
