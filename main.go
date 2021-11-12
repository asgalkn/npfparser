package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const URLADDR = "https://npfb.ru/grafik-vyplaty-pensii.php"

func main() {
	client := buildHttpClient()

	data, err := getData(client)
	if err != nil {
		log.Println("Не удается установить соединение:", err)
		return
	}

	title, err := getTitle(data)
	if err != nil {
		log.Println("Ошибка получения названия страницы:", err)
		return
	}

	payments, err := getPaymentsInfo(data)
	if err != nil {
		log.Println("Ошибка получения информации о выплатах:", err)
		return
	}

	fmt.Printf("%s\n%s\n", title, payments)
}

// Создает HTTP-клиент с дополнительными параметрами
func buildHttpClient() *http.Client {
	netTransport := &http.Transport{
		Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}

	return client
}

// Получает данные с удаленного сервера
func getData(client *http.Client) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", URLADDR, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Получает название страницы с месяцем начисления выплаты
func getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("h1.page-title").Text()
	if title == "" {
		return "", errors.New("getTitle(): data not found")
	}
	return title, nil
}

// Получает строку с указанием филиала и даты выплаты
func getPaymentsInfo(doc *goquery.Document) (string, error) {
	sel := doc.Find("table.pension-payments tbody tr td")
	if sel.Nodes == nil {
		return "", errors.New("getPaymentsInfo(): data not found")
	}

	var info string

	for i := range sel.Nodes {
		if i == 10 {
			text := sel.Eq(i).Text()
			info = strings.Join(strings.Fields(text), " ")
		}
	}

	return info, nil
}
