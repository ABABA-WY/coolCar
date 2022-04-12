// pages/lock/lock.ts
//开锁
import { routing } from "../../utils/routing";

const newShare = "Share";
const newavatarUrl = "newavatarUrl";
Page({

  /**
   * 页面的初始数据
   */
  data: {
    hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
    avatarURL: getApp<IAppOption>().globalData.avatarURL,
    userInfo: getApp<IAppOption>().globalData.userInfo,
    userInfoStr: getApp<IAppOption>().globalData.userInfoStr,
    Share: false,
  },


  onLoad(opt: Record<'car_id', string>) {
    const o: routing.LockOpts = opt;
    console.log('unlock carID:', o.car_id);
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
    this.setData({
      hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
      userInfo: getApp<IAppOption>().globalData.userInfo,
      userInfoStr: getApp<IAppOption>().globalData.userInfoStr,
      Share: wx.getStorageSync(newShare)
    })

  },

  getUserProfile: function () {
    wx.getUserProfile({
      desc: '用于微信账号与平台账号绑定', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: (res) => {
        console.log("获取到的用户信息成功: ", JSON.stringify(res));
        wx.setStorageSync(newavatarUrl, res.userInfo.avatarUrl);
        this.setData({
          hasUserInfo: true,
          userInfo: res,
          userInfoStr: JSON.stringify(res),
        })
        if (this.data.Share) {
          this.setData({
            avatarURL: res.userInfo.avatarUrl,
          })
        }
        const app = getApp<IAppOption>();
        app.globalData.hasUserInfo = this.data.hasUserInfo;
        app.globalData.userInfo = this.data.userInfo;
        app.globalData.userInfoStr = this.data.userInfoStr;
        app.globalData.avatarURL = res.userInfo.avatarUrl;
      },
      fail: (res) => {
        console.log("获取用户个人信息失败: ", res);
        const app = getApp<IAppOption>();
        app.globalData.hasUserInfo = this.data.hasUserInfo;
        app.globalData.userInfo = this.data.userInfo;
        app.globalData.userInfoStr = this.data.userInfoStr;
        app.globalData.avatarURL = '';
        //用户按了拒绝按钮
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

  onShare() {
    const share = !this.data.Share;
    if (!share) {
      this.setData({
        Share: share,
        avatarURL: '',
      })
    } else {
      this.setData({
        Share: share,
        avatarURL: wx.getStorageSync(newavatarUrl),
      })
    }
    try {
      wx.setStorageSync(newShare, share);
    } catch (err) {
      console.log(err)
    }

  },
  onUnlock() {
    if (!this.data.hasUserInfo) {
      wx.showModal({
        title: '提示',
        content: '获取个人信息',
        success: () => {
          this.getUserProfile()
          return
        }
      })
    }
    else {
      wx.getLocation({
        type: 'gcj02',
        success: loc => {

          //向服务器发送请求
          //参数
          //协议：http
          //地址 api.coolcar.cn
          //路径/trip
          //数据类型 json
          //编码
          //安全性 头部
          //错误处理 http status
          // wx.request({
          //   url: 'test.php', //仅为示例，并非真实的接口地址
          //   data: {
          //     x: '',
          //     y: ''
          //   },
          //   header: {
          //     'content-type': 'application/json' // 默认值
          //   },
          //   success(res) {
          //     if(res.statusCode === 200){

          //     }
          //     else if(res.statusCode === 404){

          //     }
          //     console.log(res.data)
          //   }
          // })
          console.log('starting a trip', {
            location: {
              latitude: loc.latitude,
              longtitude: loc.longitude,
            },

            //todo双向绑定
            avatarURL: this.data.Share ? this.data.avatarURL : '',
          })
          const tripID = 'trip123';
          wx.showLoading({
            title: '开锁中',
            mask: true,//透明门层
          })
          setTimeout(() => {
            wx.redirectTo({
              //url: `../driving/driving?trip_id=${tripID}`,
              url: routing.driving({
                trip_id: tripID,
              }),
              complete: () => {
                wx.hideLoading();
              }
            })
          }, 2000)
        },
        fail: () => {
          wx.showToast({
            title: '请前往设置页授权位置信息',
            icon: 'none',
          });
        }
      })

    }


  },

  // onLoad(){
  //   const app =  getApp<IAppOption>();
  //   app.onLoad();
  //   this.setData({
  //     hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
  //     userInfo: getApp<IAppOption>().globalData.userInfo,
  //     userInfoStr :getApp<IAppOption>().globalData.userInfoStr,
  //     Share:wx.getStorageSync(newShare)
  //   })
  // },
  // async getUserProfile(res: any){
  //   await getApp<IAppOption>().getUserProfile(res);
  //   this.setData({
  //     hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
  //     userInfo: getApp<IAppOption>().globalData.userInfo,
  //     userInfoStr :getApp<IAppOption>().globalData.userInfoStr,
  //   })
  //   if(this.data.Share){
  //     this.setData({
  //       avatarURL: getApp<IAppOption>().globalData.avatarURL,
  //     })
  //   }
  //   if(this.data.avatarURL){
  //     wx.setStorageSync(newavatarUrl,this.data.avatarURL);
  //   }
  // },


  // async onLoad() {
  //   const userInfo = await getApp<IAppOption>().globalData.userInfo;
  //   this.setData({
  //     avatarURL: userInfo.avatarUrl,
  //   })
  // },
  // onGetUserInfo(e: any) {
  //   const userInfo: WechatMiniprogram.UserInfo = e.detail.value;
  //   if (!userInfo) {
  //     return
  //   }
  //   const app = getApp<IAppOption>();
  //   app.resolveUserInfo(userInfo);
  // },
})