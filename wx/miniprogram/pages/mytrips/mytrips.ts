// pages/mytrips/mytrips.ts

import { routing } from "../../utils/routing";
interface Trip{
  id:string,
  start:string,
  end:string,
  duration:string,
  fee:string,
  distance:string,

}

//我的行程
Page({

  /**
   * 页面的初始数据
   */
  data: {
    hasUserInfo: getApp<IAppOption>().globalData.hasUserInfo,
    avatarURL: getApp<IAppOption>().globalData.avatarURL,
    userInfo: getApp<IAppOption>().globalData.userInfo,
    userInfoStr: getApp<IAppOption>().globalData.userInfoStr,
    isCertification:true,//是否认证
    indicatorDots:true,//是否显示面板指示点
    indicatorColor:'',//指示点颜色
    autoplay:false,//是否自动切换
    current:0,//当前所在滑块的 index
    currentItemId:0,
    interval:2000,//自动切换时间间隔
    duration:500,//滑动动画时长
    circular:true,//是否采用衔接滑动
    vertical:false,//滑动方向是否为纵向
    previousMargin:'0rpx',//前边距，可用于露出前一项的一小部分，接受 px 和 rpx 值
    nextMargin:'0rpx',//后边距，可用于露出后一项的一小部分，接受 px 和 rpx 值
    displayMultipleItems:1,//同时显示的滑块数量
    skipHiddenItemLayout:false,
    promotionItems:[
      {
        img:'../image/done.png',
        promotionId:1,
      },
      {
        img:'../image/lic.png',
        promotionId:2,
      },
      {
        img:'../image/OIP.jpg',
        promotionId:3,
      },
    ],
    trips:[] as Trip[],
    tripsHeight:0,
  },
  onLoad: function () {
    //查看是否授权
    this.populateTrips();
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
      avatarURL: getApp<IAppOption>().globalData.avatarURL,
    })

  },
  onReady(){
    let h:number = 0;//窗口大小
    let this_= this
    wx.createSelectorQuery().select('#heading').boundingClientRect(ret=>{
      wx.getSystemInfoAsync({
        success(res){
          h=res.windowHeight
          this_.setData({
            tripsHeight:h-ret.height,
          })
        }
      })
    }).exec();
  },
  populateTrips(){
    const trips:Trip[]=[];
    for(let i=0; i < 100;i++){
      trips.push({
        id:(1001+i).toString(),
        start:'东方明珠',
        end:'迪士尼',
        distance:'222.0公里',
        fee:'100元',
        duration:'0时44分',
      })
    }
    this.setData({
      trips:trips,
    })
  },
  onSwiperChange(){
  },
  onPromotionItemTap(e:any){
    console.log(e)
    const promotionId = e.currentTarget.dataset.promotionId;
    if(promotionId){

    }
  },
  onRegisterTap(){
    wx.navigateTo({
      url:routing.register(),
    })
  },
})