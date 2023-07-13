
// Code auto generated , DO NOT EDIT.

import type * as lb from "./lb"

export enum ErrCode {
	Success=0,
	ErrAlreadyRegister=120001,
	ErrPasswordMistake=120002,
	ErrPlayerNotFound=120003,
	ErrRoomNotFound=120004,
	ErrGameNotFound=120005,
	ErrGamePlayerNotFound=120006,
	ErrPlayCardNotFound=120007,
}

export enum GameStateChange {
	GameStateChangeNil=0,
	GameStateChangeWantLandlord=1,
	GameStateChangeGaming=2,
	GameStateChangeGameOver=3,
}

export enum Gender {
	GenderNil=0,
	GenderMale=1,
	GenderFemale=2,
}

export enum CardType {
	CardTypeNil=0,
	CardTypePassCards=1,
	CardTypeNoCards=2,
	CardTypeErrorCards=3,
	CardTypeSingleCard=4,
	CardTypeDoubleCard=5,
	CardTypeThreeCard=6,
	CardTypeThreeOneCard=7,
	CardTypeThreeTwoCard=8,
	CardTypeBombTwoCard=9,
	CardTypeStraight=10,
	CardTypeConnectCard=11,
	CardTypeAircraft=12,
	CardTypeAircraftCard=13,
	CardTypeAircraftWing=14,
	CardTypeBombCard=15,
	CardTypeKingBombCard=16,
	CardTypeBombFourCard=17,
	CardTypeBombTwoStraightCard=18,
	CardTypeBombFourStraightCard=19,
}

export enum Event_Type {
	TypeNil=0,
	TypeMatchPlayer=1,
	TypeWantLandlord=2,
	TypePlayCardIn=3,
}

export enum Webhook_Type {
	TypeNil=0,
	TypeRegisterResult=1,
	TypeLoginResult=2,
	TypeMatchResult=3,
	TypeGiveCard=4,
	TypeAckWantLandlord=5,
	TypeWantLandlordOutput=6,
	TypeWantLandlordResult=7,
	TypeStateChange=8,
	TypeAckPlayCard=9,
	TypeAckPlayCardFail=10,
	TypeAckPlayCardOut=11,
	TypeException=999,
}


