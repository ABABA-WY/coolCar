

export interface IAppOption {
    globalData: {
        //userInfo:Promise<WechatMiniprogram.UserInfo>,
        hasUserInfo:boolean,
        avatarURL:string,
        userInfo:{},
        userInfoStr :string,

    }
    //resolveUserInfo(userInfo:WechatMiniprogram.UserInfo):void
    onLoad():void,
    getUserProfile():void,

}