package task_definition_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
)

func (u *ui) handleLoadDefinitions(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page        int         `json:"page"`
		PerPage     int         `json:"per_page"`
		SortBy      string      `json:"sort_by"`
		Sort        string      `json:"sort"`
		Status      string      `json:"status"`       // legacy single filter
		CreatedFrom string      `json:"created_from"` // legacy single filter
		CreatedTo   string      `json:"created_to"`   // legacy single filter
		Conditions  []condition `json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.Page < 0 {
		reqBody.Page = 0
	}
	if reqBody.PerPage <= 0 {
		reqBody.PerPage = 10
	}
	if reqBody.SortBy == "" {
		reqBody.SortBy = taskstore.COLUMN_CREATED_AT
	}
	if reqBody.Sort == "" {
		reqBody.Sort = taskstore.DESC
	}

	query := taskstore.TaskDefinitionQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortOrder(reqBody.Sort)

	if reqBody.Status != "" {
		query.SetStatus(reqBody.Status)
	}
	if reqBody.CreatedFrom != "" {
		query.SetCreatedAtGte(reqBody.CreatedFrom + " 00:00:00")
	}
	if reqBody.CreatedTo != "" {
		query.SetCreatedAtLte(reqBody.CreatedTo + " 23:59:59")
	}

	applyConditions(query, reqBody.Conditions)

	definitions, err := store.TaskDefinitionList(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to load task definitions", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to load task definitions").ToString()))
		return ""
	}

	total, err := store.TaskDefinitionCount(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to count task definitions", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to count task definitions").ToString()))
		return ""
	}

	definitionList := []map[string]any{}
	for _, definition := range definitions {
		definitionList = append(definitionList, map[string]any{
			"id":         definition.GetID(),
			"title":      definition.GetTitle(),
			"alias":      definition.GetAlias(),
			"status":     definition.GetStatus(),
			"created_at": definition.GetCreatedAt(),
			"updated_at": definition.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Task definitions loaded successfully", map[string]any{
		"definitions": definitionList,
		"total":       total,
		"page":        reqBody.Page,
		"per_page":    reqBody.PerPage,
	}).ToString()))
	return ""
}
