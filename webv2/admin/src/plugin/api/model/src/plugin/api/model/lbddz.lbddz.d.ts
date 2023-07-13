// Code generated by rpc_gen. DO NOT EDIT.

declare namespace lbddz {
    export const enum ErrCode {
        Success = 0,
        // lbddz 120000 - 130000
        ErrAlreadyRegister = 120001, // 用户名已被注册
        ErrPasswordMistake = 120002, // 密码错误
        ErrPlayerNotFound = 120003, //
        ErrRoomNotFound = 120004, //
        ErrGameNotFound = 120005, //
        ErrGamePlayerNotFound = 120006, //
        ErrPlayCardNotFound = 120007, //
    }

    export const enum Event_Type {
        TypeNil = 0,
        TypeMatchPlayer = 1, // 匹配玩家
        TypePlayGame = 2, // 开始游戏
        TypeWantLandlord = 3, // 玩家叫/抢地主
        TypePlayCard = 4, // 出牌
    }

    export interface Event {
        type?: number;
        match_player?: MatchPlayer;
    }

    export const enum GameStateChange {
        GameStateChangeNil = 0,
        GameStateChangeWantLandlord = 1, // 抢地主阶段
        GameStateChangeGaming = 2, // 游戏中
        GameStateChangeGameOver = 3, // 游戏结束
    }

    export const enum Gender {
        GenderNil = 0,
        GenderMale = 1, // 男
        GenderFemale = 2, // 女
    }

    export const enum Webhook_Type {
        TypeNil = 0,
        TypeRegisterResult = 1, // 注册结果
        TypeLoginResult = 2, // 登录结果
        TypeMatchResult = 3, // 匹配结果
        TypeGiveCard = 4, // 发牌
        TypeStateChange = 5, // 游戏状态变更
        TypeGameOver = 6, // 对局结束
        TypeRoomExit = 7, // 房间关闭
        TypeException = 999, // 异常
    }

    export interface Webhook {
        type?: number;
        exception?: Exception;
        register?: RegisterResult;
        login?: LoginResult;
        match?: MatchResult;
        give_card?: GiveCard;
        state_change?: StateChange;
    }

    export interface ModelPlayer {
        id?: string;
        created_at?: number;
        updated_at?: number;
        deleted_at?: number;

        // @desc: 玩家昵称
        nickname?: string;

        // @desc: 登录用户名
        username?: string;

        // @desc: 玩家密码
        password?: string;

        // @desc: 玩家头像
        avatar?: string;

        // @desc: 最后登录时间
        last_login_at?: number;

        // @desc: 是否在线
        is_online?: boolean;

        // @desc: ip
        cur_ip_addr?: string;
        last_ip_addr?: string;
    }

    export interface ModelRoom {
        id?: string;
        created_at?: number;
        updated_at?: number;
        deleted_at?: number;

        // @desc: 房间创建者
        // @ref_to: ModelPlayer.id
        // type: uint64
        creator_id?: string;

        // @desc: 房间名称
        name?: string;

        // @desc: 房间玩家列表
        // @gotags: gorm:"json"
        // type: uint64
        player_ids?: Array<string>;
    }

    export interface ModelGame {
        id?: string;
        created_at?: number;
        updated_at?: number;
        deleted_at?: number;

        // @desc: 房间号
        // type: uint64
        room_id?: string;

        // @desc: 地主的位置
        landlord_seq?: number;

        // @desc: 当前出牌玩家的座位
        current_player_seq?: number;

        // @desc: 上一位出牌玩家的座位
        last_player_seq?: number;

        // @desc: 叫地主次数
        want_di_zhu_times?: number;

        // @desc: 当前叫地主分数
        cur_landlord_score?: number;

        // @desc: 地主牌
        // @gotags: gorm:"json"
        landlord_cards?: Array<number>;

        // @desc: 上一位出的牌，表示上一次出的牌
        // @gotags: gorm:"json"
        last_cards?: Array<number>;

        // @desc: 房间玩家列表
        // @gotags: gorm:"json"
        // type: uint64
        player_ids?: Array<string>;
    }

    export interface ModelGamePlayer {
        id?: string;
        created_at?: number;
        updated_at?: number;
        deleted_at?: number;
        room_id?: string;

        // @desc: 牌局ID
        // type: uint64
        game_id?: string;

        // @desc: 玩家ID
        // type: uint64
        player_id?: string;

        // @desc: 玩家在牌局中的位置，1-3
        seq?: number;

        // @desc: 玩家在牌局中的角色，地主或农民
        is_landlord?: boolean;

        // @desc: 玩家手牌
        // @gotags: gorm:"json"
        cards?: Array<number>;

        // @desc: 玩家当前手牌
        // @gotags: gorm:"json"
        cur_cards?: Array<number>;
    }

    export interface BaseGame {
        g?: ModelGame;
        gps?: Array<ModelGamePlayer>;
    }

    export interface Register {
        username?: string;
        password?: string;
    }

    export interface Login {
        username?: string;
        password?: string;
    }

    export interface MatchPlayer {
        player_id?: string;
    }

    export interface RegisterResult {
        player?: ModelPlayer;
    }

    export interface LoginResult {
        player?: ModelPlayer;
    }

    export interface MatchResult {
        players?: Array<ModelPlayer>;
        room?: ModelRoom;
    }

    export interface GiveCard {
        base_game?: BaseGame;
    }

    export interface StateChange {
        state_change?: number;
        base_game?: BaseGame;
    }

    export interface Exception {
        code?: number;
        message?: string;
    }

    export interface lbddzService {
    }
}
