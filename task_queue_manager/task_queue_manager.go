package task_queue_manager

import (
	"context"
	"github.com/dracory/taskadmin/shared"
	"net/http"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/req"

	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const fieldQueueID = shared.FIELD_QUEUE_ID

// UiInterface defines the task_queue_manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskQueueManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_queue_manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskQueueManager handles the task_queue_manager controller requests
func (u *ui) TaskQueueManager(w http.ResponseWriter, r *http.Request) {
	shared.WriteTag(w, u.ToTag(w, r))
}

func (c *ui) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	layout := c.Layout(r)
	data, errorMessage := c.prepareData(r)

	layout.SetTitle("Queue Manager | Tasks")

	if errorMessage != "" {
		layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(layout.Render(w, r))
	}

	if data.action == shared.ACTION_MODAL_QUEUE_FILTER_SHOW {
		return c.onModalRecordFilterShow(&data)
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

func (*ui) onModalRecordFilterShow(data *taskQueueManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	fieldQueueID := form.NewField(form.FieldOptions{
		Label: "Queue ID",
		Name:  shared.FIELD_FILTER_QUEUE_ID,
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formQueueID,
		Help:  `Find queue by reference number (ID).`,
	})

	filterForm := form.NewForm(form.FormOptions{
		ID:        "FormFilters",
		Method:    http.MethodGet,
		ActionURL: shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_MANAGER, map[string]string{}),
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  shared.FIELD_FILTER_STATUS,
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the queue.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Canceled",
						Key:   taskstore.TaskQueueStatusCanceled,
					},
					{
						Value: "Deleted",
						Key:   taskstore.TaskQueueStatusDeleted,
					},
					{
						Value: "Failed",
						Key:   taskstore.TaskQueueStatusFailed,
					},
					{
						Value: "Paused",
						Key:   taskstore.TaskQueueStatusPaused,
					},
					{
						Value: "Queued",
						Key:   taskstore.TaskQueueStatusQueued,
					},
					{
						Value: "Running",
						Key:   taskstore.TaskQueueStatusRunning,
					},
					{
						Value: "Success",
						Key:   taskstore.TaskQueueStatusSuccess,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label: "Name",
				Name:  shared.FIELD_FILTER_NAME,
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formName,
				Help:  `Filter by name.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created From",
				Name:  shared.FIELD_FILTER_CREATED_FROM,
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedFrom,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created To",
				Name:  shared.FIELD_FILTER_CREATED_TO,
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedTo,
				Help:  `Filter by creation date.`,
			}),
			fieldQueueID,
			form.NewField(form.FieldOptions{
				Label: "Path",
				Name:  "path",
				Type:  form.FORM_FIELD_TYPE_HIDDEN,
				Value: shared.CONTROLLER_TASK_QUEUE_MANAGER,
				Help:  `Path to this page.`,
			}),
		},
	}).Build()

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}

func (controller *ui) page(data *taskQueueManagerControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.Store(), controller.Logger(), data.request)
	breadcrumbs := shared.Breadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Queue Manager",
			URL:  "",
		},
	})

	buttonQueueCreate := hb.Button().
		Class("btn btn-primary btn-sm").
		Child(hb.I().Class("bi bi-plus-circle me-1")).
		HTML("New Task to Queue").
		HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_CREATE, map[string]string{})).
		HxTarget("body").
		HxSwap("beforeend")

	return hb.Div().
		Class("container py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Queue Manager", buttonQueueCreate)).
		Child(adminHeader).
		Child(controller.tableRecords(data))
}

func (controller *ui) tableRecords(data *taskQueueManagerControllerData) hb.TagInterface {
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
						Child(controller.sortableColumnLabel(data, "Started", "started_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Finished", "completed_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Duration", "duration")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", taskstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(queuedTask taskstore.TaskQueueInterface, _ int) hb.TagInterface {
				task, taskExists := lo.Find(data.taskList, func(t taskstore.TaskDefinitionInterface) bool {
					return t.GetID() == queuedTask.GetTaskID()
				})

				taskName := lo.IfF(taskExists, func() string { return task.GetTitle() }).Else("Unknown")
				taskAlias := lo.IfF(taskExists, func() string { return task.GetAlias() }).Else("Unknown")

				buttonDelete := hb.Button().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete task from queue").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_DELETE, map[string]string{
						"queue_id": queuedTask.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonParameters := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-list-stars")).
					Title("See queued task parameters").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_PARAMETERS, map[string]string{
						"queue_id": queuedTask.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonDetails := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-info-circle-fill")).
					Title("See the details of the job run").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_DETAILS, map[string]string{
						"queue_id": queuedTask.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonAddToQueue := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-database-add")).
					Title("Add as new task to the queue").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_REQUEUE, map[string]string{
						"queue_id": queuedTask.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonRestart := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-arrow-clockwise")).
					Title("Restart this job").
					HxGet(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_TASK_RESTART, map[string]string{
						"queue_id": queuedTask.GetID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				startedAtDate := lo.IfF(!queuedTask.GetStartedAt().IsZero(), func() string {
					return queuedTask.GetStartedAtCarbon().Format("d M Y")
				}).Else("-")
				startedAtTime := lo.IfF(!queuedTask.GetStartedAt().IsZero(), func() string {
					return queuedTask.GetStartedAtCarbon().ToTimeString()
				}).Else("-")
				completeddAtDate := lo.IfF(!queuedTask.GetCompletedAt().IsZero(), func() string {
					return queuedTask.GetCompletedAtCarbon().Format("d M Y")
				}).Else("-")
				completeddAtTime := lo.IfF(!queuedTask.GetCompletedAt().IsZero(), func() string {
					return queuedTask.GetCompletedAtCarbon().ToTimeString()
				}).Else("-")

				elapsedTime := lo.IfF(!queuedTask.GetStartedAt().IsZero() && !queuedTask.GetCompletedAt().IsZero(), func() string {
					diffSeconds := queuedTask.GetCompletedAtCarbon().DiffAbsInSeconds(queuedTask.GetStartedAtCarbon())
					return cast.ToString(diffSeconds) + "s"
				}).Else("-")

				createdAtDate := queuedTask.GetCreatedAtCarbon().Format("d M Y")
				createdAtTime := queuedTask.GetCreatedAtCarbon().ToTimeString()

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(queuedTask.IsSuccess(), `color:green;`).
					StyleIf(queuedTask.IsRunning(), `color:silver;`).
					StyleIf(queuedTask.IsQueued(), `color:blue;`).
					StyleIf(queuedTask.IsFailed(), `color:red;`).
					HTML(queuedTask.GetStatus())

				return hb.TR().
					// Name, Alias, Ref
					Child(hb.TD().
						Child(hb.Div().Text(taskName)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Alias: ").
							Text(taskAlias)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Ref: ").
							Text(queuedTask.GetID()))).

					// Status
					Child(hb.TD().
						Child(status)).

					// Started At
					Child(hb.TD().
						Child(hb.Div().Text(startedAtDate)).
						Child(hb.Div().Text(startedAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Completed At
					Child(hb.TD().
						Child(hb.Div().Text(completeddAtDate)).
						Child(hb.Div().Text(completeddAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Duration
					Child(hb.TD().
						Child(hb.Div().Text(elapsedTime)).
						Style("white-space: nowrap;")).

					// Created At
					Child(hb.TD().
						Child(hb.Div().Text(createdAtDate)).
						Child(hb.Div().Text(createdAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Actions
					Child(hb.TD().
						Style("text-align: center;").
						Child(buttonParameters).
						Child(buttonDetails).
						Child(buttonAddToQueue).
						Child(buttonRestart).
						Child(buttonDelete))
			})),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.recordCount), data.pageInt, data.perPage),
	})
}

func (controller *ui) sortableColumnLabel(data *taskQueueManagerControllerData, tableLabel, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == taskstore.ASC, taskstore.DESC).Else(taskstore.ASC)

	if !isSelected {
		direction = taskstore.ASC
	}

	link := shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_MANAGER, map[string]string{
		"controller":     shared.CONTROLLER_TASK_QUEUE_MANAGER,
		"page":           "0",
		"by":             columnName,
		"sort":           direction,
		"date_from":      data.formCreatedFrom,
		"date_to":        data.formCreatedTo,
		"status":         data.formStatus,
		"filter_task_id": data.formTaskID,
		"queue_id":       data.formQueueID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(shared.SortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *ui) tableFilter(data *taskQueueManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_MANAGER, map[string]string{
			"action":       shared.ACTION_MODAL_QUEUE_FILTER_SHOW,
			"name":         data.formName,
			"status":       data.formStatus,
			"queue_id":     data.formQueueID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing queues").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formQueueID != "" {
		description = append(description, hb.Span().Text("and queue id: "+data.formQueueID).ToHTML())
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
	data *taskQueueManagerControllerData,
	count,
	page,
	perPage int,
) hb.TagInterface {
	url := shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_MANAGER, map[string]string{
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

func (controller *ui) prepareData(r *http.Request) (data taskQueueManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.queueID = req.GetStringTrimmed(r, "queue_id")

	data.page = req.GetStringTrimmed(r, "page")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(req.GetStringTrimmedOr(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = req.GetStringTrimmed(r, "sort")
	data.sortBy = req.GetStringTrimmed(r, "by")

	data.formCreatedFrom = req.GetStringTrimmed(r, shared.FIELD_FILTER_CREATED_FROM)
	data.formCreatedTo = req.GetStringTrimmed(r, shared.FIELD_FILTER_CREATED_TO)
	data.formName = req.GetStringTrimmed(r, shared.FIELD_FILTER_NAME)
	data.formQueueID = req.GetStringTrimmed(r, shared.FIELD_FILTER_QUEUE_ID)
	data.formStatus = req.GetStringTrimmed(r, shared.FIELD_FILTER_STATUS)

	data.recordList, data.recordCount, err = controller.fetchRecordList(&data)

	if err != nil {
		controller.Logger().Error("At ui > prepareData", "error", err.Error())
		return data, "error retrieving web queues"
	}

	data.taskList, err = controller.Store().TaskDefinitionList(context.Background(), taskstore.TaskDefinitionQuery().
		SetOrderBy(taskstore.COLUMN_ALIAS).
		SetSortOrder(taskstore.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		controller.Logger().Error("At ui > prepareData", "error", err.Error())
		return data, "error retrieving tasks"
	}

	return data, ""
}

func (controller *ui) fetchRecordList(data *taskQueueManagerControllerData) (records []taskstore.TaskQueueInterface, recordCount int64, err error) {
	queueIDs := []string{}

	if data.formQueueID != "" {
		queueIDs = append(queueIDs, data.formQueueID)
	}

	query := taskstore.TaskQueueQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if len(queueIDs) > 0 {
		query = query.SetIDIn(queueIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	recordList, err := controller.Store().TaskQueueList(context.Background(), query)

	if err != nil {
		return records, 0, err
	}

	recordCount, err = controller.Store().TaskQueueCount(context.Background(), query)

	if err != nil {
		return []taskstore.TaskQueueInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type taskQueueManagerControllerData struct {
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
	formQueueID     string
	formTaskID      string

	recordList  []taskstore.TaskQueueInterface
	recordCount int64

	queueID  string
	taskList []taskstore.TaskDefinitionInterface
}
