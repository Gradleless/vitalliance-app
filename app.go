package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token                 string `json:"token"`
		IsEmployee            bool   `json:"isEmployee"`
		PkAdv                 int    `json:"pkAdv"`
		PkCrf                 int    `json:"pkCrf"`
		LastName              string `json:"lastName"`
		FirstName             string `json:"firstName"`
		AccountCreationStatus int    `json:"accountCreationStatus"`
		ApplicationStatus     int    `json:"applicationStatus"`
		Email                 string `json:"email"`
		IsAdvReferent         bool   `json:"isAdvReferent"`
		IsMandataire          bool   `json:"isMandataire"`
		IsCurrentlyMandataire bool   `json:"isCurrentlyMandataire"`
		IsPrestataire         bool   `json:"isPrestataire"`
	} `json:"data"`
}

type PointageData struct {
	IDMission     string  `json:"idMission"`
	PKHorodatage  int     `json:"pkHorodatage"`
	Horodatage    int64   `json:"horodatage"`
	PKClient      int64   `json:"pkClient"`
	Origine       string  `json:"origine"`
	Infos         string  `json:"infos"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	ClientName    string  `json:"clientName"`
	TypePointage  int     `json:"typePointage"`
	PKAdv         int64   `json:"pkAdv"`
	AdvName       string  `json:"advName"`
	IsVitaphone   bool    `json:"isVitaphone"`
	UIDHorodatage string  `json:"uidHorodatage"`
}

type PointageTime struct {
	StartPointageDate string `json:"startPointageDate"`
	Date              string `json:"date"`
	Hours             int    `json:"hours"`
	Minutes           int    `json:"minutes"`
	ClientName        string `json:"clientName"`
	Infos             string `json:"infos"`
	EndPointageDate   string `json:"endPointageDate"`
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Login(username string, password string) (*LoginResponse, error) {
	url := "https://adv-portal-api.vitacloud.fr/api/adv/auth"
	method := "POST"

	payload := []byte(`{"version":"0.0.27","os":"Firefox","username":"` + username + `","password":"` + password + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("User-Agent", "Firefox")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("username:password")))
	req.Header.Add("Origin", "https://intervenant.vitalliance.fr/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://intervenant.vitalliance.fr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Priority", "u=1")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &loginResponse, nil
}

func (a *App) GetPointage(bearer string) (*[]PointageData, error) {

	endMonth := strconv.FormatInt(time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.Now().Location()).Unix()*1000, 10)
	url := "https://adv-portal-api.vitacloud.fr/api/adv/telegestion/pointage/history/1712448000000/" + endMonth

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Firefox")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Authorization", "Bearer "+bearer)
	req.Header.Add("Origin", "https://intervenant.vitalliance.fr/")
	req.Header.Add("DNT", "1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://intervenant.vitalliance.fr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var response struct {
		Data []PointageData `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pointages := append([]PointageData{}, response.Data...)
	return &pointages, nil
}

func OrderByMonth(pointages []PointageData) map[string][]PointageData {
	orderByMonth := make(map[string][]PointageData)

	for _, pointage := range pointages {

		horodatageTime := time.Unix(pointage.Horodatage/1000, 0)

		month := horodatageTime.Format("January")
		year := horodatageTime.Format("2006")

		month = convertToFrenchMonth(month)

		key := month + " " + year

		orderByMonth[key] = append(orderByMonth[key], pointage)
	}

	return orderByMonth
}

// Dites-moi qu'il y a un autre moyen de faire ça, je vous en prie, je vous en supplie, je vous en conjure
func convertToFrenchMonth(month string) string {
	switch month {
	case "January":
		return "Janvier"
	case "February":
		return "Février"
	case "March":
		return "Mars"
	case "April":
		return "Avril"
	case "May":
		return "Mai"
	case "June":
		return "Juin"
	case "July":
		return "Juillet"
	case "August":
		return "Août"
	case "September":
		return "Septembre"
	case "October":
		return "Octobre"
	case "November":
		return "Novembre"
	case "December":
		return "Décembre"
	default:
		return month
	}
}

func (a *App) CalculateTimeHorodatage(pointages []PointageData) (map[string][]PointageTime, error) {
	pointagestime := make(map[string][]PointageTime)

	var pointagesOrdered = OrderByMonth(pointages)

	for month, pointages := range pointagesOrdered {

		for i := 0; i < len(pointages)-1; i += 2 {
			start := pointages[i]
			end := pointages[i+1]
			startTime := time.Unix(start.Horodatage/1000, 0)
			endTime := time.Unix(end.Horodatage/1000, 0)

			if pointages[i].ClientName != pointages[i+1].ClientName {
				fmt.Println("Error: Start and end pointages are not for the same client")
			}

			between := startTime.Sub(endTime)
			hours := int(between.Hours())
			minutes := int(between.Minutes()) - hours*60
			date := startTime.Format("2006-01-02")

			pointagestime[month] = append(pointagestime[month], PointageTime{
				StartPointageDate: startTime.Format("15:04"),
				Date:              date,
				Hours:             hours,
				Minutes:           minutes,
				ClientName:        start.ClientName,
				Infos:             start.Infos,
				EndPointageDate:   endTime.Format("15:04"),
			})
		}
	}
	return pointagestime, nil

}
