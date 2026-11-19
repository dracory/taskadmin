package task_definition_update

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
)

// UiInterface defines the task_definition_update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskDefinitionUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_definition_update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskDefinitionUpdate handles the task_definition_update controller requests
func (u *ui) TaskDefinitionUpdate(w http.ResponseWriter, r *http.Request) {
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
		return c.formSubmitted(&data)
	}

	return c.modal(&data)
}

func (c *ui) formSubmitted(data *taskDefinitionUpdateControllerData) hb.TagInterface {
	if data.formTitle == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Title is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if data.formAlias == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Alias is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if data.formStatus == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Status is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	data.task.
		SetTitle(data.formTitle).
		SetAlias(data.formAlias).
		SetStatus(data.formStatus).
		SetDescription(data.formDescription)

	err := c.Store().TaskDefinitionUpdate(context.Background(), data.task)

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

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Success",
			Text:              "Task successfully updated.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *ui) modal(data *taskDefinitionUpdateControllerData) *hb.Tag {
	fieldTitleVal := form.NewField(form.FieldOptions{
		Label:    "Title",
		Name:     shared.FIELD_TITLE,
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.formTitle,
		Help:     "The title of the task as displayed in the dashboard.",
		Required: true,
	})

	fieldAliasVal := form.NewField(form.FieldOptions{
		Label:    "Alias / Command Name",
		Name:     shared.FIELD_ALIAS,
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.formAlias,
		Help:     "The alias / the command name of the task. Should be unique.",
		Required: true,
	})

	fieldStatusVal := form.NewField(form.FieldOptions{
		Label:    "Status",
		Name:     shared.FIELD_STATUS,
		Type:     form.FORM_FIELD_TYPE_SELECT,
		Value:    data.formStatus,
		Help:     "The status of the task.",
		Required: true,
		Options: []form.FieldOption{
			{
				Value: "-- select status --",
				Key:   "",
			},
			{
				Value: "Active",
				Key:   taskstore.TaskDefinitionStatusActive,
			},
			{
				Value: "Inactive",
				Key:   taskstore.TaskDefinitionStatusCanceled,
			},
		},
	})

	fieldDescriptionVal := form.NewField(form.FieldOptions{
		Label:    "Description",
		Name:     shared.FIELD_DESCRIPTION,
		Type:     form.FORM_FIELD_TYPE_TEXTAREA,
		Value:    data.formDescription,
		Help:     "The description of the task.",
		Required: true,
	})

	fieldTaskIDVal := form.NewField(form.FieldOptions{
		Label:    "Task ID",
		Name:     shared.FIELD_TASK_ID,
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.taskID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormTaskUpdate",
		Fields: []form.FieldInterface{
			fieldTitleVal,
			fieldAliasVal,
			fieldStatusVal,
			fieldDescriptionVal,
			fieldTaskIDVal,
		},
	})

	modalCloseScript := `document.getElementById('ModalTaskUpdate').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonUpdate := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Save").
		Class("btn btn-success float-end").
		HxInclude(`#ModalTaskUpdate`).
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_DEFINITION_UPDATE, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalTaskUpdate").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Edit Task").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonUpdate),
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

func (c *ui) prepareData(r *http.Request) (data taskDefinitionUpdateControllerData, err error) {
	data.request = r

	data.taskID = req.GetStringTrimmed(r, shared.FIELD_TASK_ID)

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

	if r.Method == http.MethodGet {
		data.formAlias = data.task.GetAlias()
		data.formDescription = data.task.GetDescription()
		data.formStatus = data.task.GetStatus()
		data.formTitle = data.task.GetTitle()
	}

	if r.Method == http.MethodPost {
		data.formAlias = req.GetStringTrimmed(r, shared.FIELD_ALIAS)
		data.formDescription = req.GetStringTrimmed(r, shared.FIELD_DESCRIPTION)
		data.formStatus = req.GetStringTrimmed(r, shared.FIELD_STATUS)
		data.formTitle = req.GetStringTrimmed(r, shared.FIELD_TITLE)
	}

	return data, nil
}

type taskDefinitionUpdateControllerData struct {
	request *http.Request
	taskID  string
	task    taskstore.TaskDefinitionInterface

	formAlias       string
	formDescription string
	formStatus      string
	formTitle       string
}
