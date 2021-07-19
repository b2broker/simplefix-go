package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
)

type UnderlyingInstrument struct {
	*fix.Component
}

func makeUnderlyingInstrument() *UnderlyingInstrument {
	return &UnderlyingInstrument{fix.NewComponent(
		fix.NewKeyValue(FieldUnderlyingSymbol, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSymbolSfx, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSecurityID, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSecurityIDSource, &fix.String{}),
		NewUnderlyingSecurityAltIDGrp().Group,
		fix.NewKeyValue(FieldUnderlyingProduct, &fix.Int{}),
		fix.NewKeyValue(FieldUnderlyingCFICode, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSecurityType, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSecuritySubType, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingMaturityMonthYear, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingMaturityDate, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingCouponPaymentDate, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingIssueDate, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingRepoCollateralSecurityType, &fix.Int{}),
		fix.NewKeyValue(FieldUnderlyingRepurchaseTerm, &fix.Int{}),
		fix.NewKeyValue(FieldUnderlyingRepurchaseRate, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingFactor, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingCreditRating, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingInstrRegistry, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingCountryOfIssue, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingStateOrProvinceOfIssue, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingLocaleOfIssue, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingRedemptionDate, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingStrikePrice, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingStrikeCurrency, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingOptAttribute, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingContractMultiplier, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingCouponRate, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingSecurityExchange, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingIssuer, &fix.String{}),
		fix.NewKeyValue(FieldEncodedUnderlyingIssuerLen, &fix.Int{}),
		fix.NewKeyValue(FieldEncodedUnderlyingIssuer, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingSecurityDesc, &fix.String{}),
		fix.NewKeyValue(FieldEncodedUnderlyingSecurityDescLen, &fix.Int{}),
		fix.NewKeyValue(FieldEncodedUnderlyingSecurityDesc, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingCPProgram, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingCPRegType, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingCurrency, &fix.String{}),
		fix.NewKeyValue(FieldUnderlyingQty, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingPx, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingDirtyPrice, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingEndPrice, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingStartValue, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingCurrentValue, &fix.Float{}),
		fix.NewKeyValue(FieldUnderlyingEndValue, &fix.Float{}),
		makeUnderlyingStipulations().Component,
	)}
}

func NewUnderlyingInstrument() *UnderlyingInstrument {
	return makeUnderlyingInstrument()
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSymbol() string {
	kv := underlyingInstrument.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSymbol(underlyingSymbol string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSymbol)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSymbolSfx() string {
	kv := underlyingInstrument.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSymbolSfx(underlyingSymbolSfx string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSymbolSfx)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityID() string {
	kv := underlyingInstrument.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityID(underlyingSecurityID string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecurityID)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityIDSource() string {
	kv := underlyingInstrument.Get(3)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityIDSource(underlyingSecurityIDSource string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(3).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecurityIDSource)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityAltIDGrp() *UnderlyingSecurityAltIDGrp {
	group := underlyingInstrument.Get(4).(*fix.Group)

	return &UnderlyingSecurityAltIDGrp{group}
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityAltIDGrp(noUnderlyingSecurityAltID *UnderlyingSecurityAltIDGrp) *UnderlyingInstrument {
	underlyingInstrument.Set(4, noUnderlyingSecurityAltID.Group)

	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingProduct() int {
	kv := underlyingInstrument.Get(5)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingProduct(underlyingProduct int) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(5).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingProduct)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCFICode() string {
	kv := underlyingInstrument.Get(6)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCFICode(underlyingCFICode string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(6).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCFICode)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityType() string {
	kv := underlyingInstrument.Get(7)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityType(underlyingSecurityType string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(7).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecurityType)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecuritySubType() string {
	kv := underlyingInstrument.Get(8)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecuritySubType(underlyingSecuritySubType string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(8).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecuritySubType)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingMaturityMonthYear() string {
	kv := underlyingInstrument.Get(9)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingMaturityMonthYear(underlyingMaturityMonthYear string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(9).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingMaturityMonthYear)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingMaturityDate() string {
	kv := underlyingInstrument.Get(10)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingMaturityDate(underlyingMaturityDate string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(10).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingMaturityDate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCouponPaymentDate() string {
	kv := underlyingInstrument.Get(11)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCouponPaymentDate(underlyingCouponPaymentDate string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(11).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCouponPaymentDate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingIssueDate() string {
	kv := underlyingInstrument.Get(12)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingIssueDate(underlyingIssueDate string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(12).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingIssueDate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingRepoCollateralSecurityType() int {
	kv := underlyingInstrument.Get(13)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingRepoCollateralSecurityType(underlyingRepoCollateralSecurityType int) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(13).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingRepoCollateralSecurityType)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingRepurchaseTerm() int {
	kv := underlyingInstrument.Get(14)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingRepurchaseTerm(underlyingRepurchaseTerm int) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(14).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingRepurchaseTerm)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingRepurchaseRate() float64 {
	kv := underlyingInstrument.Get(15)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingRepurchaseRate(underlyingRepurchaseRate float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(15).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingRepurchaseRate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingFactor() float64 {
	kv := underlyingInstrument.Get(16)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingFactor(underlyingFactor float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(16).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingFactor)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCreditRating() string {
	kv := underlyingInstrument.Get(17)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCreditRating(underlyingCreditRating string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(17).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCreditRating)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingInstrRegistry() string {
	kv := underlyingInstrument.Get(18)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingInstrRegistry(underlyingInstrRegistry string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(18).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingInstrRegistry)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCountryOfIssue() string {
	kv := underlyingInstrument.Get(19)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCountryOfIssue(underlyingCountryOfIssue string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(19).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCountryOfIssue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingStateOrProvinceOfIssue() string {
	kv := underlyingInstrument.Get(20)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingStateOrProvinceOfIssue(underlyingStateOrProvinceOfIssue string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(20).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingStateOrProvinceOfIssue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingLocaleOfIssue() string {
	kv := underlyingInstrument.Get(21)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingLocaleOfIssue(underlyingLocaleOfIssue string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(21).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingLocaleOfIssue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingRedemptionDate() string {
	kv := underlyingInstrument.Get(22)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingRedemptionDate(underlyingRedemptionDate string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(22).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingRedemptionDate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingStrikePrice() float64 {
	kv := underlyingInstrument.Get(23)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingStrikePrice(underlyingStrikePrice float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(23).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingStrikePrice)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingStrikeCurrency() string {
	kv := underlyingInstrument.Get(24)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingStrikeCurrency(underlyingStrikeCurrency string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(24).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingStrikeCurrency)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingOptAttribute() string {
	kv := underlyingInstrument.Get(25)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingOptAttribute(underlyingOptAttribute string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(25).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingOptAttribute)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingContractMultiplier() float64 {
	kv := underlyingInstrument.Get(26)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingContractMultiplier(underlyingContractMultiplier float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(26).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingContractMultiplier)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCouponRate() float64 {
	kv := underlyingInstrument.Get(27)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCouponRate(underlyingCouponRate float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(27).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCouponRate)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityExchange() string {
	kv := underlyingInstrument.Get(28)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityExchange(underlyingSecurityExchange string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(28).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecurityExchange)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingIssuer() string {
	kv := underlyingInstrument.Get(29)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingIssuer(underlyingIssuer string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(29).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingIssuer)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) EncodedUnderlyingIssuerLen() int {
	kv := underlyingInstrument.Get(30)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (underlyingInstrument *UnderlyingInstrument) SetEncodedUnderlyingIssuerLen(encodedUnderlyingIssuerLen int) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(30).(*fix.KeyValue)
	_ = kv.Load().Set(encodedUnderlyingIssuerLen)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) EncodedUnderlyingIssuer() string {
	kv := underlyingInstrument.Get(31)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetEncodedUnderlyingIssuer(encodedUnderlyingIssuer string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(31).(*fix.KeyValue)
	_ = kv.Load().Set(encodedUnderlyingIssuer)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingSecurityDesc() string {
	kv := underlyingInstrument.Get(32)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingSecurityDesc(underlyingSecurityDesc string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(32).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingSecurityDesc)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) EncodedUnderlyingSecurityDescLen() int {
	kv := underlyingInstrument.Get(33)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (underlyingInstrument *UnderlyingInstrument) SetEncodedUnderlyingSecurityDescLen(encodedUnderlyingSecurityDescLen int) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(33).(*fix.KeyValue)
	_ = kv.Load().Set(encodedUnderlyingSecurityDescLen)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) EncodedUnderlyingSecurityDesc() string {
	kv := underlyingInstrument.Get(34)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetEncodedUnderlyingSecurityDesc(encodedUnderlyingSecurityDesc string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(34).(*fix.KeyValue)
	_ = kv.Load().Set(encodedUnderlyingSecurityDesc)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCPProgram() string {
	kv := underlyingInstrument.Get(35)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCPProgram(underlyingCPProgram string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(35).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCPProgram)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCPRegType() string {
	kv := underlyingInstrument.Get(36)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCPRegType(underlyingCPRegType string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(36).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCPRegType)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCurrency() string {
	kv := underlyingInstrument.Get(37)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCurrency(underlyingCurrency string) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(37).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCurrency)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingQty() float64 {
	kv := underlyingInstrument.Get(38)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingQty(underlyingQty float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(38).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingQty)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingPx() float64 {
	kv := underlyingInstrument.Get(39)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingPx(underlyingPx float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(39).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingPx)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingDirtyPrice() float64 {
	kv := underlyingInstrument.Get(40)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingDirtyPrice(underlyingDirtyPrice float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(40).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingDirtyPrice)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingEndPrice() float64 {
	kv := underlyingInstrument.Get(41)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingEndPrice(underlyingEndPrice float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(41).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingEndPrice)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingStartValue() float64 {
	kv := underlyingInstrument.Get(42)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingStartValue(underlyingStartValue float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(42).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingStartValue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingCurrentValue() float64 {
	kv := underlyingInstrument.Get(43)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingCurrentValue(underlyingCurrentValue float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(43).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingCurrentValue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingEndValue() float64 {
	kv := underlyingInstrument.Get(44)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(float64)
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingEndValue(underlyingEndValue float64) *UnderlyingInstrument {
	kv := underlyingInstrument.Get(44).(*fix.KeyValue)
	_ = kv.Load().Set(underlyingEndValue)
	return underlyingInstrument
}

func (underlyingInstrument *UnderlyingInstrument) UnderlyingStipulations() *UnderlyingStipulations {
	component := underlyingInstrument.Get(45).(*fix.Component)

	return &UnderlyingStipulations{component}
}

func (underlyingInstrument *UnderlyingInstrument) SetUnderlyingStipulations(underlyingStipulations *UnderlyingStipulations) *UnderlyingInstrument {
	underlyingInstrument.Set(45, underlyingStipulations.Component)

	return underlyingInstrument
}
