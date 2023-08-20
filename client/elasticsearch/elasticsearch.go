package elasticsearch

import (
	"fmt"
)

type ElasticSearch struct {
	client *elasticsearch.Client
	index  string
	alias  string
}

func New(addresses []string) (*ElasticSearch, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticSearch{
		client: client,
	}, nil
}

func (e *ElasticSearch) CreateIndex(index string) error {
	e.index = index
	e.alias = index + "_alias"

	res, err := e.client.Indices.Exists([]string{e.index})
	if err != nil {
		return fmt.Errorf("cannot check index existence: %w", err)
	}
	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return fmt.Errorf("error in index existence response: %s", res.String())
	}

	res, err = e.client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error in index creation response: %s", res.String())
	}

	res, err = e.client.Indices.PutAlias([]string{e.index}, e.alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error in index alias creation response: %s", res.String())
	}

	return nil
}

type document struct {
	Source interface{} `json:"_source"`
}
