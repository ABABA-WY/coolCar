// pages/driving/driving.ts

import { routing } from "../../utils/routing";

const centPerSec = 9;//每秒多少钱
Page({
  timer:undefined as number|undefined,
  /**
   * 页面的初始数据
   */
  data: {
    setting : {
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
    scale:14,
    location :{
      latitude: 33.099994,
      longitude: 113.324520,
    },
    elapsed:'00:00:00',
    fee:''
  },

  onLoad(opt:Record<'trip_id',string>){
    //进行保护
    const o:routing.DrivingOpts = opt
    console.log("current trip:",o.trip_id)
    this.setupLocationUpdator();
    this.setUpdateTimer();
  },
  onUnload(){
    wx.stopLocationUpdate();
    if(this.timer){
      clearInterval(this.timer);
    }
  },
//跟新位置
  setupLocationUpdator(){
    //更新
    wx.startLocationUpdate({
        fail:console.error,
    })
    wx.onLocationChange(loc=>{
      console.log('loc:',loc);
      this.setData({
        location:{
          latitude:loc.latitude,
          longitude:loc.longitude,
        }
      })
    })
  },
  setUpdateTimer(){
    let elapsedSec = 0;
    let cents = 0;
    //setInterval() 方法会不停地调用函数，直到 clearInterval() 
    this.timer = setInterval(()=>{
      elapsedSec++;
      cents = elapsedSec*centPerSec;
      this.setData({
        elapsed:formatDuration(elapsedSec),
        fee:famatFee(cents),
      })
    },1000)

  },
})

//计算时间
function formatDuration(sec:number):string{
  const padString = (n:number)=>{
    return n < 10?'0'+n.toFixed(0):n.toFixed(0);
  }
  const h = Math.floor(sec/3600);
  sec -= h*3600;
  const m = Math.floor(sec/60);
  sec -= m*60;
  const s = Math.floor(sec);
  return `${padString(h)}:${padString(m)}:${padString(s)}`
}

//计算金钱
function famatFee(cents:number){
  return (cents/100).toFixed(2);
}