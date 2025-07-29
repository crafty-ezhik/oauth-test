package handlers

import (
	"encoding/json"
	"github.com/crafty-ezhik/oauth-test/utils"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type AuthHandler interface {
	GetGoogleAuthRedirectURI() http.HandlerFunc
	GoogleCode() http.HandlerFunc
}

type AuthHandlerImpl struct{}

func NewAuthHandler() AuthHandler {
	return &AuthHandlerImpl{}
}

func (h *AuthHandlerImpl) GetGoogleAuthRedirectURI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Запрос на авторизацию через Google")
		url, err := utils.GenerateGoogleAuthURL()
		if err != nil {
			http.Error(w, err.Error(), http.StatusFound)
			return
		}
		slog.Info("Делаем перенаправление на страницу авторизации Google")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (h *AuthHandlerImpl) GoogleCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Пришел callback")
		slog.Info("Читаем body")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var reqBody RequestBody
		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Info("Начинаем процедуру обмена кода на пару ключей")
		// ОБмениваем полученный код на пару токенов
		googleTokeUrl := "https://oauth2.googleapis.com/token"

		slog.Info("Устанавливаем необходимые параметры запроса")
		reqData := url.Values{}
		reqData.Set("client_id", os.Getenv("OAUTH_GOOGLE_CLIENT_ID"))
		reqData.Set("client_secret", os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"))
		reqData.Set("grant_type", "authorization_code")
		reqData.Set("code", reqBody.Code)
		reqData.Set("redirect_uri", os.Getenv("OAUTH_GOOGLE_REDIRECT_URI"))

		slog.Info("Делаем запрос на получение токенов")
		// Делаем запрос на получение
		resp, err := http.PostForm(googleTokeUrl, reqData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Читаем тело ответа")
		// Читаем тело ответа
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var jsonRespData ResponseBody
		err = json.Unmarshal(respBody, &jsonRespData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Парсим JWT")
		parsedData, _, err := new(jwt.Parser).ParseUnverified(jsonRespData.IdToken, jwt.MapClaims{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Формируем данные о пользователе")
		w.Header().Set("Content-Type", "application/json")
		answer := map[string]interface{}{
			"user": map[string]interface{}{
				"email":   parsedData.Claims.(jwt.MapClaims)["email"],
				"picture": parsedData.Claims.(jwt.MapClaims)["picture"],
				"name":    parsedData.Claims.(jwt.MapClaims)["name"],
			},
		}

		slog.Info("Получаем данные о файлах на GoogleDisk")
		diskUrl := "https://www.googleapis.com/drive/v3/files"
		req, err := http.NewRequest("GET", diskUrl, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+jsonRespData.AccessToken)

		fileResp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fileResp.Body.Close()

		fileData, err := io.ReadAll(fileResp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonDataFile DiskPayload
		err = json.Unmarshal(fileData, &jsonDataFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Создадим слайс с именами файлов, чтобы потом вернуть клиенту
		files := make([]string, 0, len(jsonDataFile.Files))
		for _, v := range jsonDataFile.Files {
			files = append(files, v.Name)
		}

		// Добавим к ответу
		answer["files"] = files

		slog.Info("Формируем JSON и отправляем клиенту")
		answerBytes, err := json.Marshal(answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(answerBytes)
	}
}
