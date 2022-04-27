//import { getSetting, getUserInfo } from "./utils/util"
import {IAppOption} from "./appoption"
//import {auth} from "./service/proto_gen/auth/auth_pb";
//import camelcaseKeys from "camelcase-keys";
//import {rental} from "./service/proto_gen/rental/rental_pb";
import {Coolcar} from "./service/request";
// import {coolcar} from "./service/proto_gen/trip_pb";
// import camelcaseKeys from "camelcase-keys";
// app.ts
App<IAppOption>({
    globalData: {
        // userInfo:new Promise((resolve,reject)=>{
        //   resolveUserInfo = resolve
        //   rejectUserInfo = reject
        // }),
        hasUserInfo: false,
        avatarURL: '',
        userInfo: {},
        userInfoStr: '',
    },

    onLaunch() {

        Coolcar.login()
        // 登录
        // wx.login({
        //     success: res => {
        //         console.log(res.code)
        //         // 发送 res.code 到后台换取 openId, sessionKey, unionId
        //         wx.request({
        //             url: "http://127.0.0.1:8000/v1/auth/login",
        //             method: "POST",
        //             data: {
        //                 code: res.code,
        //             } as auth.v1.ILoginRequest,
        //             success: res => {
        //                 //camelcaseKeys将对象命名改为驼峰形式
        //                 const loginResponse: auth.v1.LoginResponse = auth.v1.LoginResponse.fromObject(camelcaseKeys(res.data as object))
        //                 console.log(loginResponse)
        //                 wx.request({
        //                     url: "http://127.0.0.1:8000/v1/trip",
        //                     method: "POST",
        //                     data: {
        //                         code:'abc',
        //                     } as rental.v1.ICreateTripRequest,
        //                     header: {
        //                         authorization: 'Bearer ' + loginResponse.accessToken,
        //                     },
        //                     success: res => {
        //                         const createTripRequest: rental.v1.CreateTripResponse = rental.v1.CreateTripResponse.fromObject(camelcaseKeys(res.data as object))
        //                         console.log(createTripRequest)
        //                     },
        //                     fail: console.error,
        //                 })
        //             },
        //             fail: console.error,
        //         })
        //     },
        // })
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
                this.globalData.userInfo = res;
                this.globalData.avatarURL = res.userInfo.avatarUrl;
            },
            fail: (res) => {
                console.log("获取用户个人信息失败: ", res);
                //用户按了拒绝按钮
                this.globalData.hasUserInfo = false;
                this.globalData.userInfo = {};
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