package task_definition_manager

import (
	"context"
	"github.com/dracory/taskadmin/shared"
	"net/http"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const actionModalTaskFilterShow = "modal_task_filter_show"

// UiInterface defines the task_definition_manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskDefinitionManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_definition_manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskDefinitionManager handles the task_definition_manager controller requests
func (u *ui) TaskDefinitionManager(w http.ResponseWriter, r *http.Request) {
	shared.WriteTag(w, u.ToTag(w, r))
}

func (c *ui) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	layout := c.Layout(r)
	data, errorMessage := c.prepareData(r)

	layout.SetTitle("Task Manager | Tasks")

	if errorMessage != "" {
		layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
		}
	}, 1000);`

	swalScript := `setTimeout(() => {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
		}
	}, 1000);`

	layout.SetBody(c.page(&data).ToHTML())
	layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(layout.Render(w, r))
}

func (controller *ui) page(data *taskDefinitionManagerControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.Store(), controller.Logger(), data.request)
	breadcrumbs := shared.Breadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Task Manager",
			URL:  "",
		},
	})

	buttonTaskCreate := hb.Button().
		Class("btn btn-primary btn-sm").
		Child(hb.I().Class("bi bi-plus-circle me-1")).
		HTML("New Task").
		HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_CREATE, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	return hb.Div().
		Class("container py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Task Manager", buttonTaskCreate)).
		Child(adminHeader).
		Child(controller.tableRecords(data))
}

func (controller *ui) tableRecords(data *taskDefinitionManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Name", taskstore.COLUMN_TITLE)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Alias", taskstore.COLUMN_ALIAS)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Reference", taskstore.COLUMN_ID)).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Status", taskstore.COLUMN_STATUS)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", taskstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Modified", taskstore.COLUMN_UPDATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(task taskstore.TaskDefinitionInterface, _ int) hb.TagInterface {
				taskName := task.GetTitle()
				taskAlias := task.GetAlias()

				buttonDelete := hb.Button().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete task").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_DELETE, map[string]string{
						"task_id": task.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonUpdate := hb.Button().
					Class("btn btn-sm btn-success").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit task").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_UPDATE, map[string]string{
						"task_id": task.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				createdAtDate := task.GetCreatedAtCarbon().Format("d M Y")
				createdAtTime := task.GetCreatedAtCarbon().ToTimeString()
				updatedAtDate := task.GetUpdatedAtCarbon().Format("d M Y")
				updatedAtTime := task.GetUpdatedAtCarbon().ToTimeString()

				status := hb.Span().
					Style(`font-weight: bold;`).
					// StyleIf(taskdTask.IsFailed(), `color:red;`).
					HTML(task.GetStatus())

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Text(taskName)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Alias: ").
							Text(taskAlias)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Ref: ").
							Text(task.GetID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().Text(createdAtDate)).
						Child(hb.Div().Text(createdAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Child(hb.Div().Text(updatedAtDate)).
						Child(hb.Div().Text(updatedAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Style("text-align: center;").
						Child(buttonUpdate).
						Child(buttonDelete),
				})
			})),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.recordCount), data.pageInt, data.perPage),
	})
}

func (controller *ui) sortableColumnLabel(data *taskDefinitionManagerControllerData, tableLabel, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == taskstore.ASC, taskstore.DESC).Else(taskstore.ASC)

	if !isSelected {
		direction = taskstore.ASC
	}

	link := shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_MANAGER, map[string]string{
		"controller":     shared.CONTROLLER_TASK_DEFINITION_MANAGER,
		"page":           "0",
		"by":             columnName,
		"sort":           direction,
		"date_from":      data.formCreatedFrom,
		"date_to":        data.formCreatedTo,
		"status":         data.formStatus,
		"filter_task_id": data.formTaskID,
		"task_id":        data.formTaskID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(shared.SortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *ui) tableFilter(data *taskDefinitionManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_MANAGER, map[string]string{
			"action":       actionModalTaskFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"task_id":      data.formTaskID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing tasks").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formTaskID != "" {
		description = append(description, hb.Span().Text("and task id: "+data.formTaskID).ToHTML())
	}

	// 	if data.formTaskID != "" {

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *ui) tablePagination(
	data *taskDefinitionManagerControllerData,
	count,
	page,
	perPage int,
) hb.TagInterface {
	url := shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_MANAGER, map[string]string{
		"status":       data.formStatus,
		"name":         data.formName,
		"created_from": data.formCreatedFrom,
		"created_to":   data.formCreatedTo,
		"by":           data.sortBy,
		"order":        data.sortOrder,
	})

	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (controller *ui) prepareData(r *http.Request) (data taskDefinitionManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.taskID = req.GetStringTrimmed(r, "task_id")

	data.page = req.GetStringTrimmed(r, "page")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(req.GetStringTrimmedOr(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = req.GetStringTrimmedOr(r, "sort", taskstore.DESC)
	data.sortBy = req.GetStringTrimmedOr(r, "by", taskstore.COLUMN_CREATED_AT)
	data.formCreatedFrom = req.GetStringTrimmed(r, "filter_created_from")
	data.formCreatedTo = req.GetStringTrimmed(r, "filter_created_to")
	data.formName = req.GetStringTrimmed(r, "filter_name")
	data.formTaskID = req.GetStringTrimmed(r, "filter_task_id")
	data.formStatus = req.GetStringTrimmed(r, "filter_status")

	data.recordList, data.recordCount, err = controller.fetchRecordList(&data)

	if err != nil {
		controller.Logger().Error("At ui > prepareData", "error", err.Error())
		return data, "error retrieving web tasks"
	}

	return data, ""
}

func (controller *ui) fetchRecordList(data *taskDefinitionManagerControllerData) (records []taskstore.TaskDefinitionInterface, recordCount int64, err error) {
	taskIDs := []string{}

	if data.formTaskID != "" {
		taskIDs = append(taskIDs, data.formTaskID)
	}

	query := taskstore.TaskDefinitionQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if len(taskIDs) > 0 {
		query = query.SetIDIn(taskIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	recordList, err := controller.Store().TaskDefinitionList(context.Background(), query)

	if err != nil {
		return records, 0, err
	}

	recordCount, err = controller.Store().TaskDefinitionCount(context.Background(), query)

	if err != nil {
		return []taskstore.TaskDefinitionInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type taskDefinitionManagerControllerData struct {
	request *http.Request
	action  string

	page      string
	pageInt   int
	perPage   int
	sortOrder string
	sortBy    string

	formStatus      string
	formName        string
	formCreatedFrom string
	formCreatedTo   string
	formTaskID      string

	recordList  []taskstore.TaskDefinitionInterface
	recordCount int64

	taskID string
}
