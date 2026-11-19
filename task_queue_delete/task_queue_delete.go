package task_queue_delete

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

const fieldQueueID = shared.FIELD_QUEUE_ID

// UiInterface defines the task_queue_delete controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskQueueDelete(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_queue_delete controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskQueueDelete handles the task_queue_delete controller requests
func (u *ui) TaskQueueDelete(w http.ResponseWriter, r *http.Request) {
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

func (c *ui) formSubmitted(data taskQueueDeleteControllerData) hb.TagInterface {
	if err := c.Store().TaskQueueSoftDelete(context.Background(), data.queue); err != nil {
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
			Text:              "Queue successfully deleted.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *ui) modal(data taskQueueDeleteControllerData) *hb.Tag {
	fieldDanger := form.NewField(form.FieldOptions{
		Type: form.FORM_FIELD_TYPE_RAW,
		Value: hb.Wrap().
			Child(hb.Paragraph().
				Child(hb.Text(`You are about to permanently delete this queued task:`))).
			Child(hb.Paragraph().
				Style("font-weight: bold;").
				Child(hb.Text(`Ref. "` + data.queue.GetID() + `"`))).
			Child(hb.Paragraph().
				Child(hb.Text(`Are you sure you want to proceed?`))).
			Child(hb.Paragraph().
				Class("text-danger").
				Child(hb.Text(`Warning: This action cannot be undone.`))).
			ToHTML(),
		Required: true,
	})

	fieldQueueID := form.NewField(form.FieldOptions{
		Label:    "Queue ID",
		Name:     fieldQueueID,
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.queueID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormQueueDelete",
		Fields: []form.FieldInterface{
			fieldDanger,
			fieldQueueID,
		},
	})

	modalCloseScript := `document.getElementById('ModalQueueDelete').remove();document.getElementById('ModalBackdrop').remove();`
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
		HxInclude(`#ModalQueueDelete`).
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_DELETE, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalQueueDelete").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Delete Queue").
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

func (c *ui) prepareData(r *http.Request) (data taskQueueDeleteControllerData, err error) {
	data.request = r

	data.queueID = req.GetStringTrimmed(r, fieldQueueID)

	if data.queueID == "" {
		return data, errors.New("queue_id is required")
	}

	data.queue, err = c.Store().TaskQueueFindByID(context.Background(), data.queueID)

	if err != nil {
		return data, err
	}

	if data.queue == nil {
		return data, errors.New("queue not found")
	}

	return data, nil
}

type taskQueueDeleteControllerData struct {
	request *http.Request
	queueID string
	queue   taskstore.TaskQueueInterface
}
