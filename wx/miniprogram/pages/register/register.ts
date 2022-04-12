// pages/register/register.ts

import { routing } from "../../utils/routing";

//审查
Page({

  /**
   * 页面的初始数据
   */
  data: {
    licImageURL:'',
    licNo:'',
    licName:'',
    genderIndex:0,
    genders:["其他","男","女","其他"],
    birthDate:'1980-01-01',
    phone:'',
    state:'unsubmitted' as 'unsubmitted'|'pending'|'submittedFailed'|'submittedSucceed'
  },
  redirectURL:'',
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(opt:Record<'redirect',string>) {
    const o:routing.RegisterOpts = opt;
    if(o.redirect){
      this.redirectURL=decodeURIComponent(o.redirect);
      console.log(this.redirectURL);
    }
  },

  onUpload(){
    wx.chooseImage({
      success:res=>{
        if(res.tempFilePaths.length > 0){
          this.setData({
            licImageURL:res.tempFilePaths[0],
          })
          setTimeout(()=>{
            this.setData({
              licNo:'123456',
              licName:'张三',
              genderIndex:1,
              birthDate:'1999-01-21',
              phone:'9876412'
            })
          },2000)
        }
      },
    })
  },
  onGenderChange(e:any){
    this.setData({
      genderIndex:e.detail.value,
    })
  },
  onBirthdayChange(e:any){
    this.setData({
      birthDate:e.detail.value,
    })
  },
  onSubmit(){
    this.setData({
      state:'pending',
    })
    setTimeout(()=>{
      this.onVerified()
    },2000)
  },
  onReSubmit(){
    this.setData({
      state:"unsubmitted",
      licImageURL:'',
    })
  },
  onVerified(){
    this.setData({
      state:'submittedSucceed',
    })
  },
  onSubmittedSucceed(){
    if(this.redirectURL){
      wx.redirectTo({
        url:this.redirectURL,
      })
    }
    else{
      wx.redirectTo({
        url:routing.mytrips(),
      })
    }

  },
})