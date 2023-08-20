package elasticsearch

import (
	"bytes"
	"context"
	"elastic-project/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"time"
)

type UserInfoStorage struct {
	elastic ElasticSearch
	timeout time.Duration
}

type UserInfoStorer interface {
	Insert(ctx context.Context, userInfo UserInfo) error
	Update(ctx context.Context, userInfo UserInfo) error
	Delete(ctx context.Context, id string) error
	FindOne(ctx context.Context, id string) (UserInfo, error)
	FindByKeyAndValue(queryType string, key string, value string) ([]UserInfo, error)
	FindByQuery(jsonString string) ([]UserInfo, error)
}

type UserInfo struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Job        string     `json:"job"`
	ChildNames []string   `json:"childNames"`
	Comment    string     `json:"comment"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

func NewUserInfoStorage(elastic ElasticSearch) UserInfoStorer {
	return &UserInfoStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}
}

func (p UserInfoStorage) Insert(ctx context.Context, userInfo UserInfo) error {
	bdy, err := json.Marshal(userInfo)
	if err != nil {
		return fmt.Errorf("insert: marshall: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: userInfo.ID,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return model.ErrConflict
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (p UserInfoStorage) Update(ctx context.Context, userInfo UserInfo) error {
	bdy, err := json.Marshal(userInfo)
	if err != nil {
		return fmt.Errorf("update: marshall: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      p.elastic.alias,
		DocumentID: userInfo.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("update: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (p UserInfoStorage) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("delete: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}

	return nil
}

func (p UserInfoStorage) FindOne(ctx context.Context, id string) (UserInfo, error) {
	req := esapi.GetRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return UserInfo{}, fmt.Errorf("find one: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return UserInfo{}, model.ErrNotFound
	}

	if res.IsError() {
		return UserInfo{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var (
		userInfo UserInfo
		body     document
	)
	body.Source = &userInfo

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return UserInfo{}, fmt.Errorf("find one: decode: %w", err)
	}

	return userInfo, nil
}

func (p UserInfoStorage) FindByKeyAndValue(queryType string, key string, value string) ([]UserInfo, error) {
	return p.Search(queryType, key, value)
}

func (p UserInfoStorage) Search(queryType string, key string, value string) ([]UserInfo, error) {
	var userInfoList []UserInfo
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			queryType: map[string]interface{}{
				key: value,
			},
		},
	}
	err := json.NewEncoder(&buffer).Encode(query)
	if err != nil {
		return []UserInfo{}, err
	}
	es := p.elastic.client
	response, err := es.Search(es.Search.WithIndex("user_alias"), es.Search.WithBody(&buffer))
	if err != nil {
		return []UserInfo{}, err
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return []UserInfo{}, err
	}
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
		createdAt, _ := time.Parse(time.RFC3339Nano, craft["created_at"].(string))
		userInfo := UserInfo{ID: craft["id"].(string),
			Name:       craft["name"].(string),
			Job:        craft["job"].(string),
			ChildNames: convertToStringArray(craft["childNames"].([]interface{})),
			Comment:    craft["comment"].(string),
			CreatedAt:  &createdAt,
		}
		userInfoList = append(userInfoList, userInfo)
	}
	return userInfoList, nil
}

func (p UserInfoStorage) FindByQuery(jsonString string) ([]UserInfo, error) {
	var userInfoList []UserInfo
	var buffer bytes.Buffer
	err := json.Indent(&buffer, []byte(jsonString), "", "  ")
	if err != nil {
		return []UserInfo{}, err
	}
	es := p.elastic.client
	response, err := es.Search(es.Search.WithIndex("user_alias"), es.Search.WithBody(&buffer))
	if err != nil {
		return []UserInfo{}, err
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return []UserInfo{}, err
	}
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
		createdAt, _ := time.Parse(time.RFC3339Nano, craft["created_at"].(string))
		userInfo := UserInfo{ID: craft["id"].(string),
			Name:       craft["name"].(string),
			Job:        craft["job"].(string),
			ChildNames: convertToStringArray(craft["childNames"].([]interface{})),
			Comment:    craft["comment"].(string),
			CreatedAt:  &createdAt,
		}
		userInfoList = append(userInfoList, userInfo)
	}
	return userInfoList, nil
}

func convertToStringArray(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		if s, ok := v.(string); ok {
			result[i] = s
		}
	}
	return result
}
