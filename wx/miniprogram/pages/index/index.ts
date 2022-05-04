// index.ts
// 获取应用实例
import {IAppOption} from "../../appoption"
import { routing } from "../../utils/routing";

Page({
  isPagesShowing: true,
  data: {
    setting: {
      skew: 0,
      rotate: 0,
      showLocation: true,
      showScale: false,
      subKey: '',
      layerStyle: 1,
      enableZoom: true,
      enableScroll: true,
      enableRotate: false,
      showCompass: false,
      enable3D: false,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,

    },
    scale: 10,
    hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
    userInfo: getApp<IAppOption>().globalData.userInfo,
    userInfoStr: getApp<IAppOption>().globalData.userInfoStr,
    avatarURL: getApp<IAppOption>().globalData.avatarURL,
    motto: "hello world from ts",
    location: {
      latitude: 33.099994,
      longitude: 113.324520,
    },
    markers: [
      {
        id: 0,
        iconPath: "../image/OIP.jpg",
        latitude: 23.099994,
        longitude: 113.324520,
        width: 50,
        height: 50
      },

    ]
  },
  getUserProfile: function () {
    wx.getUserProfile({
      desc: '用于微信账号与平台账号绑定', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: (res) => {
        console.log("获取用户信息成功");
        this.setData({
          hasUserInfo: true,
          userInfo: res,
          userInfoStr: JSON.stringify(res),
          avatarURL: res.userInfo.avatarUrl,
        })
        const app = getApp<IAppOption>();
        app.globalData.hasUserInfo = this.data.hasUserInfo;
        app.globalData.userInfo = this.data.userInfo;
        app.globalData.userInfoStr = this.data.userInfoStr;
        app.globalData.avatarURL = res.userInfo.avatarUrl;
      },
      fail: (res) => {
        console.log("获取用户个人信息失败: ", res);
        const app = getApp<IAppOption>();
        app.globalData.hasUserInfo = false;
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
  onMyLocationTap() {
    wx.getLocation({
      type: 'gcj02',
      success: res => {
        this.setData({
          location: {
            latitude: res.latitude,
            longitude: res.longitude,
          },
        })
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: "请求失败",
        })
      }
    })
  },
  moveCars() {
    const map = wx.createMapContext("map")
    const dest = {
      latitude: 23.099994,
      longitude: 113.324520,
    }
    const moveCar = () => {
      dest.latitude += 0.1,
        dest.longitude += 0.1,
        map.translateMarker({
          destination: {
            latitude: dest.latitude,
            longitude: dest.longitude,
          },
          markerId: 0,
          autoRotate: false,
          rotate: 0,
          duration: 5000,
          animationEnd: () => {
            if (this.isPagesShowing) {
              moveCar()
            }
          },
        })
    }
    moveCar()
  },
  onShow() {
    this.isPagesShowing = true;
    this.setData({
      hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
      userInfo: getApp<IAppOption>().globalData.userInfo,
      userInfoStr: getApp<IAppOption>().globalData.userInfoStr,
      avatarURL: getApp<IAppOption>().globalData.avatarURL,
    })
  },
  onHide() {

    this.isPagesShowing = false;
  },
  onScanCliked() {
    wx.scanCode({
      success: () => {
        wx.showModal({
          title: '提示',
          content: '进行身份认证',
          success: () => {
            const carID = 'car123';
            //const redirectURL=`../lock/lock?car_id=${carID}`;
            const redirectURL = routing.lock({
              car_id: carID,
            })
            wx.navigateTo({
              //先注册，在解锁
              //url: `../register/register?redirect=${encodeURIComponent(redirectURL)}`
              url: routing.register({
                redirectURL: redirectURL,
              })
            })
          }
        })

      },
      fail: console.error,
    })
  },
  onMytripsTap() {
    if (!this.data.hasUserInfo) {
      this.getUserProfile();
    }
    if (this.data.hasUserInfo) {
      wx.navigateTo({
        url: routing.mytrips(),
      })
    }

  }
})
