package task_queue_requeue

import (
	"context"
	"encoding/json"
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

const fieldQueueID = shared.FIELD_QUEUE_ID

// UiInterface defines the task_queue_requeue controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskQueueRequeue(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task_queue_requeue controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskQueueRequeue handles the task_queue_requeue controller requests
func (u *ui) TaskQueueRequeue(w http.ResponseWriter, r *http.Request) {
	shared.WriteTag(w, u.ToTag(w, r))
}

func (c *ui) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modal(data)
}

// TODO add queue name, currently set to degault queue
func (c *ui) formSubmitted(data taskQueueRequeueControllerData) hb.TagInterface {
	if data.formParameters == "" {
		data.formParameters = "{}"
	}

	if !shared.IsJSON(data.formParameters) {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task Parameters is not valid JSON", Position: "top-right"})
	}

	task, err := c.Store().TaskDefinitionFindByID(context.Background(), data.queue.GetTaskID())

	if err != nil {
		c.Logger().Error("At ui > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	if task == nil {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task not found", Position: "top-right"})
	}

	taskParametersAny := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.formParameters), &taskParametersAny); err != nil {
		c.Logger().Error("At ui > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	taskParametersMap := cast.ToStringMap(taskParametersAny)

	_, err = c.Store().TaskDefinitionEnqueueByAlias(context.Background(), taskstore.DefaultQueueName, task.GetAlias(), taskParametersMap)

	if err != nil {
		c.Logger().Error("At ui > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queue successfully created.", Position: "top-right"})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *ui) modal(data taskQueueRequeueControllerData) *hb.Tag {
	modalID := `ModalQueueRequeue`
	formID := modalID + `Form`

	divInfo := hb.Div().
		Class("alert alert-info").
		Text(`A new task will be created with the following parameters. You may  edit the parameters if necessary`)

	fieldInfo := form.NewField(form.FieldOptions{
		Label:    "Queue",
		Type:     form.FORM_FIELD_TYPE_RAW,
		Value:    divInfo.ToHTML(),
		Required: true,
	})

	fieldParam := form.NewField(form.FieldOptions{
		Label:    "Parameters",
		Name:     shared.FIELD_PARAMETERS,
		Type:     form.FORM_FIELD_TYPE_TEXTAREA,
		Value:    data.formParameters,
		Help:     "The parameters of the queued task. Must be valid JSON.",
		Required: true,
	})

	fieldParamSize := form.NewField(form.FieldOptions{
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: hb.Style(`#` + formID + ` textarea[name="` + shared.FIELD_PARAMETERS + `"] { height: 200px; }`).ToHTML(),
	})

	fieldQueueID := form.NewField(form.FieldOptions{
		Label:    "Queue ID",
		Name:     fieldQueueID,
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.queueID,
		Required: true,
	})

	formCreate := form.NewForm(form.FormOptions{
		ID: formID,
		Fields: []form.FieldInterface{
			fieldQueueID,
			fieldInfo,
			fieldParamSize,
			fieldParam,
		},
	})

	modalCloseScript := `document.getElementById('` + modalID + `').remove();document.getElementById('ModalBackdrop').remove();`

	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonRequeue := hb.Button().
		Child(hb.I().Class("bi bi-database-add me-2")).
		HTML("Add to queue").
		Class("btn btn-success float-end").
		HxInclude(`#` + modalID).
		HxPost(shared.URL(data.request, shared.CONTROLLER_TASK_QUEUE_REQUEUE, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Requeue as New Task to Queue").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formCreate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonRequeue),
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

func (c *ui) prepareData(r *http.Request) (data taskQueueRequeueControllerData, err error) {
	data.request = r
	data.queueID = req.GetStringTrimmed(r, fieldQueueID)

	if data.queueID == "" {
		return data, errors.New("queue_id is required")
	}

	if data.queue, err = c.Store().TaskQueueFindByID(context.Background(), data.queueID); err != nil {
		return data, err
	}

	if data.queue == nil {
		return data, errors.New("queue not found")
	}

	if r.Method == http.MethodGet {
		data.formParameters = data.queue.GetParameters()
	}

	if r.Method == http.MethodPost {
		data.formParameters = req.GetStringTrimmed(r, shared.FIELD_PARAMETERS)
	}

	return data, nil
}

type taskQueueRequeueControllerData struct {
	request *http.Request
	queueID string
	queue   taskstore.TaskQueueInterface

	formParameters string
}
