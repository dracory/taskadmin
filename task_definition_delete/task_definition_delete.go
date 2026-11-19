package task_definition_delete

import (
	"context"
	"errors"
	"github.com/dracory/taskadmin/shared"
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
	"github.com/spf13/cast"
)

const fieldTaskID = shared.FIELD_TASK_ID

// UiInterface defines the task_definition_delete controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskDefinitionDelete(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_definition_delete controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskDefinitionDelete handles the task_definition_delete controller requests
func (u *ui) TaskDefinitionDelete(w http.ResponseWriter, r *http.Request) {
	shared.WriteTag(w, u.ToTag(w, r))
}

func (c *ui) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modal(data)
}

func (c *ui) formSubmitted(data taskDefinitionDeleteControllerData) hb.TagInterface {
	for _, queuedTask := range data.relatedQueuesToDelete {
		if err := c.Store().TaskQueueSoftDeleteByID(context.Background(), queuedTask.GetID()); err != nil {
			return hb.Swal(hb.SwalOptions{
				Icon:              "error",
				Title:             "Error",
				Text:              err.Error(),
				Position:          "top-right",
				ShowCancelButton:  false,
				ShowConfirmButton: false,
			})
		}
	}

	if err := c.Store().TaskDefinitionSoftDelete(context.Background(), data.task); err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Success",
			Text:              "Task successfully deleted.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *ui) modal(data taskDefinitionDeleteControllerData) *hb.Tag {
	fieldDanger := form.NewField(form.FieldOptions{
		Type: form.FORM_FIELD_TYPE_RAW,
		Value: hb.Wrap().
			Child(hb.Paragraph().
				Child(hb.Text(`You are about to permanently delete this task definition:`))).
			Child(hb.Paragraph().
				Style("font-weight: bold;").
				Child(hb.Text(`"` + data.task.GetTitle() + `"`))).
			Child(hb.Paragraph().
				Child(hb.Text(`This will also delete all the tasks created from this definition, including ` + cast.ToString(data.relatedQueuedQueuesCount) + ` pending tasks.`))).
			Child(hb.Paragraph().
				Child(hb.Text(`Are you sure you want to proceed?`))).
			Child(hb.Paragraph().
				Class("text-danger").
				Child(hb.Text(`Warning: This action cannot be undone.`))).
			ToHTML(),
		Required: true,
	})

	fieldTaskID := form.NewField(form.FieldOptions{
		Label:    "Task ID",
		Name:     fieldTaskID,
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.taskID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormTaskDelete",
		Fields: []form.FieldInterface{
			fieldDanger,
			fieldTaskID,
		},
	})

	modalCloseScript := `document.getElementById('ModalTaskDelete').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonDelete := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Delete").
		Class("btn btn-danger float-end").
		HxInclude(`#ModalTaskDelete`).
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_DELETE, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalTaskDelete").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Delete Task").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonDelete),
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

func (c *ui) prepareData(r *http.Request) (data taskDefinitionDeleteControllerData, err error) {
	data.request = r

	data.taskID = req.GetStringTrimmed(r, fieldTaskID)

	if data.taskID == "" {
		return data, errors.New("task_id is required")
	}

	data.task, err = c.Store().TaskDefinitionFindByID(context.Background(), data.taskID)

	if err != nil {
		return data, err
	}

	if data.task == nil {
		return data, errors.New("task not found")
	}

	if r.Method == http.MethodPost {
		data.relatedQueuesToDelete, err = c.Store().TaskQueueList(context.Background(), taskstore.TaskQueueQuery().
			SetTaskID(data.task.GetID()).
			SetColumns([]string{taskstore.COLUMN_ID}))

		if err != nil {
			return data, err
		}
	}

	if r.Method == http.MethodGet {
		data.relatedQueuedQueuesCount, err = c.Store().TaskQueueCount(context.Background(), taskstore.TaskQueueQuery().
			SetTaskID(data.task.GetID()).
			SetColumns([]string{taskstore.COLUMN_ID}))

		if err != nil {
			return data, err
		}
	}

	return data, nil
}

type taskDefinitionDeleteControllerData struct {
	request                  *http.Request
	taskID                   string
	task                     taskstore.TaskDefinitionInterface
	relatedQueuesToDelete    []taskstore.TaskQueueInterface
	relatedQueuedQueuesCount int64
}
