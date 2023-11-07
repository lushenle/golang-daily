package main

type OrderState int

const (
	CREATE OrderState = iota
	PAID
	DELIVERING
	RECEIVED
	DONE
	CANCELLING
	RETURNING
	CLOSED
)

type Order struct {
	state OrderState
}

func NewOrder() *Order {
	return &Order{
		state: CREATE,
	}
}

func (o *Order) canPay() bool {
	return o.state == CREATE
}

func (o *Order) canDeliver() bool {
	return o.state == PAID
}

func (o *Order) canCancel() bool {
	return o.state == CREATE || o.state == PAID
}

func (o *Order) canReceive() bool {
	return o.state == DELIVERING
}

func (o *Order) paymentService() bool {
	// 调用 RPC 接口完成支付
	return false
}

func (o *Order) pay() bool {
	if o.canPay() {
		if o.paymentService() {
			o.state = PAID
			return true
		}
		return false
	}

	return false
}

func (o *Order) cancel() bool {
	if o.canCancel() {
		o.state = CANCELLING
		// 取消订单，申请审批和清理数据
		o.state = CLOSED
		return true
	}

	return false
}
