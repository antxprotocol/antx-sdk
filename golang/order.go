package sdk

import (
	ordertypes "github.com/antxprotocol/antx-proto/gen/go/antx/chain/order"
	"github.com/antxprotocol/antx-sdk-golang/constants"
	"github.com/antxprotocol/antx-sdk-golang/types"
)

// CreateOrder creates an order
func (c *AntxClient) CreateOrder(order *types.CreateOrderParam) (string, error) {
	msg := ordertypes.MsgCreateOrder{
		AgentAddress:      c.GetAgentAddress(),
		SubaccountId:      order.SubaccountId,
		ExchangeId:        order.ExchangeId,
		MarginMode:        order.MarginMode,
		Leverage:          order.Leverage,
		IsBuy:             order.IsBuy,
		PriceScale:        order.PriceScale,
		PriceValue:        order.PriceValue,
		SizeScale:         order.SizeScale,
		SizeValue:         order.SizeValue,
		ClientOrderId:     order.ClientOrderId,
		TimeInForce:       order.TimeInForce,
		ReduceOnly:        order.ReduceOnly,
		ExpireTime:        order.ExpireTime,
		IsMarket:          order.IsMarket,
		IsPositionTp:      order.IsPositionTp,
		IsPositionSl:      order.IsPositionSl,
		TriggerType:       order.TriggerType,
		TriggerPriceType:  order.TriggerPriceType,
		TriggerPriceValue: order.TriggerPriceValue,
		IsSetOpenTp:       order.IsSetOpenTp,
		OpenTpParam:       &order.OpenTpParam,
		IsSetOpenSl:       order.IsSetOpenSl,
		OpenSlParam:       &order.OpenSlParam,
	}

	txHash, err := c.signAndSendTx(constants.MsgCreateOrderTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateOrderBatch creates orders in batch
func (c *AntxClient) CreateOrderBatch(orders *types.CreateOrderBatchParam) (string, error) {
	batchList := make([]*ordertypes.CreateOrderParam, 0, len(orders.CreateOrderParam))
	for _, order := range orders.CreateOrderParam {
		batchList = append(batchList, &ordertypes.CreateOrderParam{
			IsBuy:             order.IsBuy,
			PriceScale:        order.PriceScale,
			PriceValue:        order.PriceValue,
			SizeScale:         order.SizeScale,
			SizeValue:         order.SizeValue,
			ClientOrderId:     order.ClientOrderId,
			TimeInForce:       order.TimeInForce,
			ReduceOnly:        order.ReduceOnly,
			ExpireTime:        order.ExpireTime,
			IsMarket:          order.IsMarket,
			IsPositionTp:      order.IsPositionTp,
			IsPositionSl:      order.IsPositionSl,
			TriggerType:       order.TriggerType,
			TriggerPriceType:  order.TriggerPriceType,
			TriggerPriceValue: order.TriggerPriceValue,
			IsSetOpenTp:       order.IsSetOpenTp,
			OpenTpParam:       &order.OpenTpParam,
			IsSetOpenSl:       order.IsSetOpenSl,
			OpenSlParam:       &order.OpenSlParam,
		})
	}

	msg := ordertypes.MsgCreateOrderBatch{
		AgentAddress:     orders.AgentAddress,
		SubaccountId:     orders.SubaccountId,
		ExchangeId:       orders.ExchangeId,
		MarginMode:       orders.MarginMode,
		Leverage:         orders.Leverage,
		CreateOrderParam: batchList,
	}

	txHash, err := c.signAndSendTx(constants.MsgCreateOrderBatchTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CancelOrder cancels an order
func (c *AntxClient) CancelOrder(order *types.CancelOrderParam) (string, error) {
	msg := ordertypes.MsgCancelOrder{
		AgentAddress: c.GetAgentAddress(),
		SubaccountId: order.SubaccountId,
		OrderId:      order.OrderIdList,
	}

	txHash, err := c.signAndSendTx(constants.MsgCancelOrderTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CancelOrderByClientId cancels an order by client ID
func (c *AntxClient) CancelOrderByClientId(order *types.CancelOrderByClientIdParam) (string, error) {
	msg := ordertypes.MsgCancelOrderByClientId{
		AgentAddress:  c.GetAgentAddress(),
		SubaccountId:  order.SubaccountId,
		ClientOrderId: order.ClientOrderIdList,
	}

	txHash, err := c.signAndSendTx(constants.MsgCancelOrderByClientIdTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CancelAllOrder cancels all orders
func (c *AntxClient) CancelAllOrder(order *types.CancelAllOrderParam) (string, error) {
	msg := ordertypes.MsgCancelAllOrder{
		AgentAddress:     c.GetAgentAddress(),
		SubaccountId:     order.SubaccountId,
		FilterExchangeId: order.FilterExchangeIdList,
	}

	txHash, err := c.signAndSendTx(constants.MsgCancelAllOrderTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CloseAllPosition closes all positions
func (c *AntxClient) CloseAllPosition(order *types.CloseAllPositionParam) (string, error) {
	msg := ordertypes.MsgCloseAllPosition{
		AgentAddress:     c.GetAgentAddress(),
		SubaccountId:     order.SubaccountId,
		FilterExchangeId: order.FilterExchangeIdList,
	}

	txHash, err := c.signAndSendTx(constants.MsgCloseAllPositionTypeURL, &msg, true)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
