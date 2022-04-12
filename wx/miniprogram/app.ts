//import { getSetting, getUserInfo } from "./utils/util"
let resolveUserInfo :(value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void
let rejectUserInfo:(reason?: any) => void
// app.ts
App<IAppOption>({
  globalData: {
    // userInfo:new Promise((resolve,reject)=>{
    //   resolveUserInfo = resolve
    //   rejectUserInfo = reject
    // }),
    hasUserInfo:false,
    avatarURL: '',
    userInfo:{},
    userInfoStr :'',
  },

  onLaunch() {
    // 展示本地存储能力

    // 登录
    wx.login({
      success: res => {
        console.log(res.code)
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
    })
  //   getSetting().then(res => {
  //     if(res.authSetting['scope.userInfo']){
  //       return getUserInfo()
  //     }
  //     return undefined
  //   }).then(res=>{
  //     if(!res){
  //       return
  //     }
  //     resolveUserInfo(res.userInfo)
      
  //   }).catch(rejectUserInfo)

  // },
  // resolveUserInfo(userInfo:WechatMiniprogram.UserInfo){
  //   resolveUserInfo(userInfo)
  // },
  },
  onLoad() {
    //查看是否授权
    wx.getSetting({
      success: function (res) {
        if (res.authSetting['scope.userInfo']) {
          console.log("用户授权了");
        } else {
          //用户没有授权
          console.log("用户没有授权");
        }
      }
    });
  
  },

  getUserProfile() {
    wx.getUserProfile({
      desc: '用于微信账号与平台账号绑定', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: (res) => {
        console.log("获取到的用户信息成功: ", JSON.stringify(res));
        //wx.setStorageSync(newavatarUrl,res.userInfo.avatarUrl);
        this.globalData.hasUserInfo = true;
        this.globalData.userInfo= res;
        this.globalData.avatarURL = res.userInfo.avatarUrl;
      },
      fail: (res) => {
        console.log("获取用户个人信息失败: ", res);
        //用户按了拒绝按钮
        this.globalData.hasUserInfo = false;
        this.globalData.userInfo= {};
        this.globalData.avatarURL = '';
        wx.showModal({
          title: '提示',
          content: '您点击了拒绝授权!!!',
          showCancel: false,
          confirmText: '返回授权',
          success: function (res) {
            // 用户没有授权成功，不需要改变 isHide 的值
            if (res.confirm) {
              console.log('用户点击了“返回授权”');
            }
          }
        });
      }
    })
  },
})