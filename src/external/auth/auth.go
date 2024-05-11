package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hcsouza/fiap-tech-fast-food/src/external/api/infra/config"
)

type AuthorizationToken struct {
	AccessToken string            `json:"AccessToken"`
	Headers     map[string]string `json:"headers"`
	StatusCode  int               `json:"statusCode"`
}

func GetAuthorizationToken(cpf string) (*AuthorizationToken, error) {
	apiURL := fmt.Sprintf("%s?cpf=%s", config.GetApiCfg().AuthorizationBaseUrl, cpf)
	fmt.Println(apiURL)
	response, err := http.Get(apiURL)

	if err != nil {
		return &AuthorizationToken{}, fmt.Errorf("erro ao fazer a solicitação HTTP")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &AuthorizationToken{}, fmt.Errorf("erro ao ler a resposta")
	}

	var apiResponse AuthorizationToken
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return &AuthorizationToken{}, fmt.Errorf("erro ao analisar o JSON")
	}
	return &apiResponse, nil
}
