package handlers

import (
	"LogC/internal/models"
	"LogC/internal/services"
	"LogC/internal/utils"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func AddLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Define the data to pass to the template
		data := struct {
		}{}

		// Parse and execute the template
		utils.OpenPage("add_log.html", data, w)
	} else if r.Method == http.MethodPost {
		// Handle POST request
		amount_str := r.FormValue("item-amount")
		amount, err := strconv.Atoi(amount_str)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		var logs []models.LogItem
		for i := 0; i < amount; i++ {
			// Handle each log item
			cur_content := ""
			cur_type_str := r.FormValue("type" + strconv.Itoa(i))
			if cur_type_str == "" {
				continue
			}
			cur_type := models.Text
			if cur_type_str == models.LogItemType.String(models.Text) {
				cur_content = r.FormValue("item" + strconv.Itoa(i))
				cur_type = models.Text
			} else if cur_type_str == models.LogItemType.String(models.Image) {
				file, _, err := r.FormFile("item" + strconv.Itoa(i))
				if err != nil {
					http.Error(w, "Invalid file", http.StatusBadRequest)
					fmt.Println(err)
					return
				}
				defer file.Close()

				var cur_data models.LogData
				cur_data.Data, err = io.ReadAll(file)
				if err != nil {
					http.Error(w, "Invalid file", http.StatusBadRequest)
					fmt.Println(err)
					return
				}
				added_id, err := services.SaveImageToDB(cur_data)
				if err != nil {
					http.Error(w, "Invalid file", http.StatusBadRequest)
					fmt.Println(err)
					return
				}
				cur_content = strconv.Itoa(added_id)

				cur_type = models.Image
			} else if cur_type_str == models.LogItemType.String(models.Title) {
				cur_content = r.FormValue("item" + strconv.Itoa(i))
				cur_type = models.Title
			}
			cur_log := models.LogItem{Content: cur_content, Type: cur_type, Order: i}
			logs = append(logs, cur_log)
		}

		title := r.FormValue("title")
		date := time.Now()
		log := models.Log{Title: title, Date: date, Items: logs}

		//services.SaveLogsToFile(log)
		services.SaveLogsToDB(log)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
	}
}
