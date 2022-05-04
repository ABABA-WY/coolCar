import camelcaseKeys from "camelcase-keys";
import {auth} from "./proto_gen/auth/auth_pb";

export namespace Coolcar{
    const serverAddr = 'http://127.0.0.1:8000'
    const AUTH_ERR ='AUTH_ERR'
    const authData = {
        token :'',
        expiresMs:0,
    }
    export interface RequestOption<REQ,RES> {
        method: 'GET'|'PUT'|'POST'|'DELETE'
        path:string
        data:REQ
        resMarshaller:(res:object)=>RES
    }
    //控制头部是否输入
    export interface AuthOption {
        attachAuthHeader:boolean
        retryAuthERR:boolean
    }
    //流程管理函数
    export async function sendRequestWithAuthRetry<REQ,RES>(req :RequestOption<REQ,RES>,a?:AuthOption) :Promise<RES>{
        const authOpt = a ||{
            attachAuthHeader : true,
            retryAuthERR:true,
        }
        try {
            await login()
            //console.log("req:",req)
            return sendRequest(req,authOpt)
        }catch (err) {
            if (err === AUTH_ERR && authOpt.retryAuthERR) {
                authData.token = ''
                authData.expiresMs = 0
                return  sendRequestWithAuthRetry(req, {
                    attachAuthHeader: authOpt.attachAuthHeader,
                    retryAuthERR: false,
                })
            }else {
                throw err
            }
        }

    }
    //登录管理
    export async function login(){
        //如果token有效
        if (authData.token && authData.expiresMs >= Date.now()) {
            return
        }
        const wxResp = await wxLogin()
        const reqTime = Date.now()
        const response = await sendRequest<auth.v1.ILoginRequest,auth.v1.ILoginResponse>({
            path: "/v1/auth/login",
            method: "POST",
            data: {
                code: wxResp.code,
            },
            resMarshaller:auth.v1.LoginResponse.fromObject,
        },{
            //登录啥都不用取
            attachAuthHeader:false,
            retryAuthERR:false,
        })
        authData.token = response.accessToken!
        authData.expiresMs = reqTime + response.expiresIn! * 1000//ms
    }

    //发送请求
    function sendRequest<REQ,RES>(req :RequestOption<REQ,RES>,a:AuthOption) :Promise<RES>{

        return new Promise((resolve,reject)=>{
            const  header: Record<string, any>= {}
            if (a.attachAuthHeader ){
                if(authData.token && authData.expiresMs >= Date.now()){
                    header.authorization = 'Bearer '+authData.token
                    //console.log(header.authorization)
                }else{
                    reject(AUTH_ERR)
                    return
                }
            }
            console.log(serverAddr+req.path)
            wx.request({
                url :serverAddr+req.path,
                method:req.method,
                data:req.data,
                header,
                success:res=>{
                    //console.log("res:",res)
                    //头部无效||不加头部的token
                    if(res.statusCode === 401){
                        //清除状态
                        authData.token=''
                        authData.expiresMs=0
                        reject(AUTH_ERR)
                    }else if(res.statusCode >= 400){
                        reject(res)
                    }else{
                        resolve(req.resMarshaller(camelcaseKeys(res.data as object,{deep:true})))
                    }
                },
                fail :reject,

            })
        })
    }
    function wxLogin():Promise<WechatMiniprogram.LoginSuccessCallbackResult>{
        return new Promise((resolve,reject)=>{
            wx.login({
                success:resolve,
                fail : reject
            })
        })
    }
}