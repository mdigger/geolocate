package geolocate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	RequestTimeout = time.Second * 5 // максимальное время ожидания ответа
	IgnoreIPMethod = true            // не использовать определение по IP-адресу
)

// Ошибки, возвращаемые при запросе данных стандартного сервиса геолокации.
var (
	ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest)) // неверный формат данных запроса или плохой ключ
	ErrForbidden  = errors.New(http.StatusText(http.StatusForbidden))  // исчерпан лимит запросов
	ErrNotFound   = errors.New(http.StatusText(http.StatusNotFound))   // информация не найдена
)

// URL сервисов геолокации.
const (
	Mozilla = "https://location.services.mozilla.com/v1/geolocate"
	Google  = "https://www.googleapis.com/geolocation/v1/geolocate"
	Yandex  = "http://api.lbs.yandex.net/geolocation"
)

// Locator описывает интерфейс, поддерживаемый всеми типами сервисов геолокации.
type Locator interface {
	Get(req Request) (*Response, error)
}

// base описывает информацию о сервисе геолокации, использующем стандартный тип
// запросов, такие как Mozilla и Google Locator.
type base struct {
	serviceUrl string       // адрес для запроса сервиса
	client     *http.Client // HTTP-клиент
}

// New возвращает новый инициализированный сервис геолокации.
func New(serviceUrl, apiKey string) (locator Locator, err error) {
	if serviceUrl == Yandex { // для Яндекса возвращаем отдельный обрабочик
		return &yandex{
			apiKey: apiKey,
			client: &http.Client{
				Timeout: RequestTimeout,
			},
		}, nil
	}
	if apiKey != "" { // добавляем ключ к запросу
		serviceUrl += "?key=" + url.QueryEscape(apiKey)
	}
	// проверяем, что URL валиден
	if _, err := url.ParseRequestURI(serviceUrl); err != nil {
		return nil, err
	}
	return &base{ // возвращаем базовый обработчик геолокации
		serviceUrl: serviceUrl,
		client: &http.Client{
			Timeout: RequestTimeout,
		},
	}, nil
}

// Get передает данные на сервер геолокации и возвращает от него разобранный ответ или ошибку.
func (l *base) Get(req Request) (*Response, error) {
	req.ConsiderIp = !IgnoreIPMethod
	if IgnoreIPMethod {
		req.Fallbacks = &Fallbacks{
			LAC: false,
			IP:  false,
		}
	}
	if req.RadioType == "" {
		req.RadioType = "gsm" // Mozilla не находит данные, если не указано
	}
	req.IPAddress = "" // не используется в этих запросах
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(data))
	resp, err := l.client.Post(l.serviceUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200: // все хорошо — данные получены
	case 400: // неверный формат данных запроса или плохой ключ
		return nil, ErrBadRequest
	case 403: // исчерпан лимит запросов
		return nil, ErrForbidden
	case 404: // информация не найдена
		return nil, ErrNotFound
	default: // другая нехорошая ошибка
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &response, nil
}
