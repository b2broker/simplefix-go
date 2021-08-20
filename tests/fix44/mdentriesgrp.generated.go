package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type MDEntriesGrp struct {
	*fix.Group
}

func NewMDEntriesGrp() *MDEntriesGrp {
	return &MDEntriesGrp{
		fix.NewGroup(FieldNoMDEntries,
			fix.NewKeyValue(FieldMDUpdateAction, &fix.String{}),
			fix.NewKeyValue(FieldDeleteReason, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryType, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryID, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryRefID, &fix.String{}),
			makeInstrument().Component,
			NewUnderlyingsGrp().Group,
			NewLegsGrp().Group,
			fix.NewKeyValue(FieldFinancialStatus, &fix.String{}),
			fix.NewKeyValue(FieldCorporateAction, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryPx, &fix.Float{}),
			fix.NewKeyValue(FieldCurrency, &fix.String{}),
			fix.NewKeyValue(FieldMDEntrySize, &fix.Float{}),
			fix.NewKeyValue(FieldMDEntryDate, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryTime, &fix.String{}),
			fix.NewKeyValue(FieldTickDirection, &fix.String{}),
			fix.NewKeyValue(FieldMDMkt, &fix.String{}),
			fix.NewKeyValue(FieldTradingSessionID, &fix.String{}),
			fix.NewKeyValue(FieldTradingSessionSubID, &fix.String{}),
			fix.NewKeyValue(FieldQuoteCondition, &fix.String{}),
			fix.NewKeyValue(FieldTradeCondition, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryOriginator, &fix.String{}),
			fix.NewKeyValue(FieldLocationID, &fix.String{}),
			fix.NewKeyValue(FieldDeskID, &fix.String{}),
			fix.NewKeyValue(FieldOpenCloseSettlFlag, &fix.String{}),
			fix.NewKeyValue(FieldTimeInForce, &fix.String{}),
			fix.NewKeyValue(FieldExpireDate, &fix.String{}),
			fix.NewKeyValue(FieldExpireTime, &fix.String{}),
			fix.NewKeyValue(FieldMinQty, &fix.Float{}),
			fix.NewKeyValue(FieldExecInst, &fix.String{}),
			fix.NewKeyValue(FieldSellerDays, &fix.Int{}),
			fix.NewKeyValue(FieldOrderID, &fix.String{}),
			fix.NewKeyValue(FieldQuoteEntryID, &fix.String{}),
			fix.NewKeyValue(FieldMDEntryBuyer, &fix.String{}),
			fix.NewKeyValue(FieldMDEntrySeller, &fix.String{}),
			fix.NewKeyValue(FieldNumberOfOrders, &fix.Int{}),
			fix.NewKeyValue(FieldMDEntryPositionNo, &fix.Int{}),
			fix.NewKeyValue(FieldScope, &fix.String{}),
			fix.NewKeyValue(FieldPriceDelta, &fix.Float{}),
			fix.NewKeyValue(FieldNetChgPrevDay, &fix.Float{}),
			fix.NewKeyValue(FieldText, &fix.String{}),
			fix.NewKeyValue(FieldEncodedTextLen, &fix.Int{}),
			fix.NewKeyValue(FieldEncodedText, &fix.String{}),
		),
	}
}

func (group *MDEntriesGrp) AddEntry(entry *MDEntriesEntry) *MDEntriesGrp {
	group.Group.AddEntry(entry.Items())

	return group
}

func (group *MDEntriesGrp) Entries() []*MDEntriesEntry {
	items := make([]*MDEntriesEntry, len(group.Group.Entries()))

	for i, item := range group.Group.Entries() {
		items[i] = &MDEntriesEntry{fix.NewComponent(item...)}
	}

	return items
}

type MDEntriesEntry struct {
	*fix.Component
}

func makeMDEntriesEntry() *MDEntriesEntry {
	return &MDEntriesEntry{fix.NewComponent(
		fix.NewKeyValue(FieldMDUpdateAction, &fix.String{}),
		fix.NewKeyValue(FieldDeleteReason, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryType, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryID, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryRefID, &fix.String{}),
		makeInstrument().Component,
		NewUnderlyingsGrp().Group,
		NewLegsGrp().Group,
		fix.NewKeyValue(FieldFinancialStatus, &fix.String{}),
		fix.NewKeyValue(FieldCorporateAction, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryPx, &fix.Float{}),
		fix.NewKeyValue(FieldCurrency, &fix.String{}),
		fix.NewKeyValue(FieldMDEntrySize, &fix.Float{}),
		fix.NewKeyValue(FieldMDEntryDate, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryTime, &fix.String{}),
		fix.NewKeyValue(FieldTickDirection, &fix.String{}),
		fix.NewKeyValue(FieldMDMkt, &fix.String{}),
		fix.NewKeyValue(FieldTradingSessionID, &fix.String{}),
		fix.NewKeyValue(FieldTradingSessionSubID, &fix.String{}),
		fix.NewKeyValue(FieldQuoteCondition, &fix.String{}),
		fix.NewKeyValue(FieldTradeCondition, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryOriginator, &fix.String{}),
		fix.NewKeyValue(FieldLocationID, &fix.String{}),
		fix.NewKeyValue(FieldDeskID, &fix.String{}),
		fix.NewKeyValue(FieldOpenCloseSettlFlag, &fix.String{}),
		fix.NewKeyValue(FieldTimeInForce, &fix.String{}),
		fix.NewKeyValue(FieldExpireDate, &fix.String{}),
		fix.NewKeyValue(FieldExpireTime, &fix.String{}),
		fix.NewKeyValue(FieldMinQty, &fix.Float{}),
		fix.NewKeyValue(FieldExecInst, &fix.String{}),
		fix.NewKeyValue(FieldSellerDays, &fix.Int{}),
		fix.NewKeyValue(FieldOrderID, &fix.String{}),
		fix.NewKeyValue(FieldQuoteEntryID, &fix.String{}),
		fix.NewKeyValue(FieldMDEntryBuyer, &fix.String{}),
		fix.NewKeyValue(FieldMDEntrySeller, &fix.String{}),
		fix.NewKeyValue(FieldNumberOfOrders, &fix.Int{}),
		fix.NewKeyValue(FieldMDEntryPositionNo, &fix.Int{}),
		fix.NewKeyValue(FieldScope, &fix.String{}),
		fix.NewKeyValue(FieldPriceDelta, &fix.Float{}),
		fix.NewKeyValue(FieldNetChgPrevDay, &fix.Float{}),
		fix.NewKeyValue(FieldText, &fix.String{}),
		fix.NewKeyValue(FieldEncodedTextLen, &fix.Int{}),
		fix.NewKeyValue(FieldEncodedText, &fix.String{}),
	)}
}

func NewMDEntriesEntry() *MDEntriesEntry {
	return makeMDEntriesEntry()
}

func (mDEntriesEntry *MDEntriesEntry) MDUpdateAction() string {
	kv := mDEntriesEntry.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDUpdateAction(mDUpdateAction string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(mDUpdateAction)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) DeleteReason() string {
	kv := mDEntriesEntry.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetDeleteReason(deleteReason string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(deleteReason)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryType() string {
	kv := mDEntriesEntry.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryType(mDEntryType string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryType)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryID() string {
	kv := mDEntriesEntry.Get(3)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryID(mDEntryID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(3).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryRefID() string {
	kv := mDEntriesEntry.Get(4)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryRefID(mDEntryRefID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(4).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryRefID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) Instrument() *Instrument {
	component := mDEntriesEntry.Get(5).(*fix.Component)

	return &Instrument{component}
}

func (mDEntriesEntry *MDEntriesEntry) SetInstrument(instrument *Instrument) *MDEntriesEntry {
	mDEntriesEntry.Set(5, instrument.Component)

	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) UnderlyingsGrp() *UnderlyingsGrp {
	group := mDEntriesEntry.Get(6).(*fix.Group)

	return &UnderlyingsGrp{group}
}

func (mDEntriesEntry *MDEntriesEntry) SetUnderlyingsGrp(noUnderlyings *UnderlyingsGrp) *MDEntriesEntry {
	mDEntriesEntry.Set(6, noUnderlyings.Group)

	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) LegsGrp() *LegsGrp {
	group := mDEntriesEntry.Get(7).(*fix.Group)

	return &LegsGrp{group}
}

func (mDEntriesEntry *MDEntriesEntry) SetLegsGrp(noLegs *LegsGrp) *MDEntriesEntry {
	mDEntriesEntry.Set(7, noLegs.Group)

	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) FinancialStatus() string {
	kv := mDEntriesEntry.Get(8)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetFinancialStatus(financialStatus string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(8).(*fix.KeyValue)
	_ = kv.Load().Set(financialStatus)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) CorporateAction() string {
	kv := mDEntriesEntry.Get(9)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetCorporateAction(corporateAction string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(9).(*fix.KeyValue)
	_ = kv.Load().Set(corporateAction)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryPx() float64 {
	kv := mDEntriesEntry.Get(10)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryPx(mDEntryPx float64) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(10).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryPx)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) Currency() string {
	kv := mDEntriesEntry.Get(11)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetCurrency(currency string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(11).(*fix.KeyValue)
	_ = kv.Load().Set(currency)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntrySize() float64 {
	kv := mDEntriesEntry.Get(12)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntrySize(mDEntrySize float64) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(12).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntrySize)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryDate() string {
	kv := mDEntriesEntry.Get(13)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryDate(mDEntryDate string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(13).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryDate)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryTime() string {
	kv := mDEntriesEntry.Get(14)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryTime(mDEntryTime string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(14).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryTime)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) TickDirection() string {
	kv := mDEntriesEntry.Get(15)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetTickDirection(tickDirection string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(15).(*fix.KeyValue)
	_ = kv.Load().Set(tickDirection)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDMkt() string {
	kv := mDEntriesEntry.Get(16)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDMkt(mDMkt string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(16).(*fix.KeyValue)
	_ = kv.Load().Set(mDMkt)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) TradingSessionID() string {
	kv := mDEntriesEntry.Get(17)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetTradingSessionID(tradingSessionID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(17).(*fix.KeyValue)
	_ = kv.Load().Set(tradingSessionID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) TradingSessionSubID() string {
	kv := mDEntriesEntry.Get(18)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetTradingSessionSubID(tradingSessionSubID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(18).(*fix.KeyValue)
	_ = kv.Load().Set(tradingSessionSubID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) QuoteCondition() string {
	kv := mDEntriesEntry.Get(19)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetQuoteCondition(quoteCondition string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(19).(*fix.KeyValue)
	_ = kv.Load().Set(quoteCondition)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) TradeCondition() string {
	kv := mDEntriesEntry.Get(20)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetTradeCondition(tradeCondition string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(20).(*fix.KeyValue)
	_ = kv.Load().Set(tradeCondition)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryOriginator() string {
	kv := mDEntriesEntry.Get(21)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryOriginator(mDEntryOriginator string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(21).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryOriginator)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) LocationID() string {
	kv := mDEntriesEntry.Get(22)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetLocationID(locationID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(22).(*fix.KeyValue)
	_ = kv.Load().Set(locationID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) DeskID() string {
	kv := mDEntriesEntry.Get(23)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetDeskID(deskID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(23).(*fix.KeyValue)
	_ = kv.Load().Set(deskID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) OpenCloseSettlFlag() string {
	kv := mDEntriesEntry.Get(24)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetOpenCloseSettlFlag(openCloseSettlFlag string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(24).(*fix.KeyValue)
	_ = kv.Load().Set(openCloseSettlFlag)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) TimeInForce() string {
	kv := mDEntriesEntry.Get(25)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetTimeInForce(timeInForce string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(25).(*fix.KeyValue)
	_ = kv.Load().Set(timeInForce)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) ExpireDate() string {
	kv := mDEntriesEntry.Get(26)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetExpireDate(expireDate string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(26).(*fix.KeyValue)
	_ = kv.Load().Set(expireDate)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) ExpireTime() string {
	kv := mDEntriesEntry.Get(27)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetExpireTime(expireTime string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(27).(*fix.KeyValue)
	_ = kv.Load().Set(expireTime)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MinQty() float64 {
	kv := mDEntriesEntry.Get(28)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (mDEntriesEntry *MDEntriesEntry) SetMinQty(minQty float64) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(28).(*fix.KeyValue)
	_ = kv.Load().Set(minQty)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) ExecInst() string {
	kv := mDEntriesEntry.Get(29)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetExecInst(execInst string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(29).(*fix.KeyValue)
	_ = kv.Load().Set(execInst)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) SellerDays() int {
	kv := mDEntriesEntry.Get(30)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (mDEntriesEntry *MDEntriesEntry) SetSellerDays(sellerDays int) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(30).(*fix.KeyValue)
	_ = kv.Load().Set(sellerDays)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) OrderID() string {
	kv := mDEntriesEntry.Get(31)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetOrderID(orderID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(31).(*fix.KeyValue)
	_ = kv.Load().Set(orderID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) QuoteEntryID() string {
	kv := mDEntriesEntry.Get(32)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetQuoteEntryID(quoteEntryID string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(32).(*fix.KeyValue)
	_ = kv.Load().Set(quoteEntryID)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryBuyer() string {
	kv := mDEntriesEntry.Get(33)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryBuyer(mDEntryBuyer string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(33).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryBuyer)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntrySeller() string {
	kv := mDEntriesEntry.Get(34)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntrySeller(mDEntrySeller string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(34).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntrySeller)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) NumberOfOrders() int {
	kv := mDEntriesEntry.Get(35)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (mDEntriesEntry *MDEntriesEntry) SetNumberOfOrders(numberOfOrders int) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(35).(*fix.KeyValue)
	_ = kv.Load().Set(numberOfOrders)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) MDEntryPositionNo() int {
	kv := mDEntriesEntry.Get(36)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (mDEntriesEntry *MDEntriesEntry) SetMDEntryPositionNo(mDEntryPositionNo int) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(36).(*fix.KeyValue)
	_ = kv.Load().Set(mDEntryPositionNo)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) Scope() string {
	kv := mDEntriesEntry.Get(37)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetScope(scope string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(37).(*fix.KeyValue)
	_ = kv.Load().Set(scope)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) PriceDelta() float64 {
	kv := mDEntriesEntry.Get(38)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (mDEntriesEntry *MDEntriesEntry) SetPriceDelta(priceDelta float64) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(38).(*fix.KeyValue)
	_ = kv.Load().Set(priceDelta)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) NetChgPrevDay() float64 {
	kv := mDEntriesEntry.Get(39)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (mDEntriesEntry *MDEntriesEntry) SetNetChgPrevDay(netChgPrevDay float64) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(39).(*fix.KeyValue)
	_ = kv.Load().Set(netChgPrevDay)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) Text() string {
	kv := mDEntriesEntry.Get(40)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetText(text string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(40).(*fix.KeyValue)
	_ = kv.Load().Set(text)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) EncodedTextLen() int {
	kv := mDEntriesEntry.Get(41)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (mDEntriesEntry *MDEntriesEntry) SetEncodedTextLen(encodedTextLen int) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(41).(*fix.KeyValue)
	_ = kv.Load().Set(encodedTextLen)
	return mDEntriesEntry
}

func (mDEntriesEntry *MDEntriesEntry) EncodedText() string {
	kv := mDEntriesEntry.Get(42)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (mDEntriesEntry *MDEntriesEntry) SetEncodedText(encodedText string) *MDEntriesEntry {
	kv := mDEntriesEntry.Get(42).(*fix.KeyValue)
	_ = kv.Load().Set(encodedText)
	return mDEntriesEntry
}
