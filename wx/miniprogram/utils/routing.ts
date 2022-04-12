export namespace routing{
    export interface DrivingOpts{
        trip_id:string,
    }
    export interface LockOpts{
        car_id:string,
    }
    export interface RegisterOpts{
        redirect?:string,
    }
    export interface RegisterParams{
        redirectURL:string,
    }
    
    export function driving(o:DrivingOpts){
        return `../driving/driving?trip_id=${o.trip_id}`
    }
    export function lock(o:LockOpts){
        return`../lock/lock?car_id=${o.car_id}`
    }
    export function register(p?:RegisterParams){
        const page = '../register/register';
        if(!p){
            return page
        }
        return `${page}?redirect=${encodeURIComponent(p.redirectURL)}`
    }
    export function mytrips(){
        return '../mytrips/mytrips'
    }
}