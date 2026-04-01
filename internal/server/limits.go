package server
type Tier string
const(TierFree Tier="free";TierPro Tier="pro")
type Limits struct{Tier Tier;Description string}
func LimitsFor(tier string)Limits{if tier=="pro"{return Limits{Tier:TierPro,Description:"Unlimited templates and renders"}}
return Limits{Tier:TierFree,Description:"3 templates, 100 renders/day"}}
func(l Limits)IsPro()bool{return l.Tier==TierPro}
