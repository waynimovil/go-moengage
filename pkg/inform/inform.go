package inform

import (
	"context"
	"github.com/waynimovil/go-moengage/internal"
	"github.com/waynimovil/go-moengage/pkg/models"
)

type Inform interface {
	SendAlert(ctx context.Context, req models.Alert) (models.AlertSuccessResponse, models.ResponseDetails, error)
}

const (
	sendAlertPath = "v1/send"
)

/*
	type Data interface {
		SendCustomer(ctx context.Context, req models.Customer)
	}
*/
type Channel struct {
	ReqHandler internal.HTTPHandler
}

func (wap *Channel) SendAlert(
	ctx context.Context,
	msg models.Alert,
) (msgResp models.AlertSuccessResponse, respDetails models.ResponseDetails, err error) {
	respDetails, err = wap.ReqHandler.PostJSONReq(ctx, &msg, &msgResp, sendAlertPath)
	return msgResp, respDetails, err
}
