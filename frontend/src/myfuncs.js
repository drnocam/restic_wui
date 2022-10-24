import { loading } from "./store"
import {  toast } from '@zerodevx/svelte-toast'



const warning_notify = (metin) => {
    toast.push(metin, { theme: { 
        '--toastBackground': 'yellow',
        '--toastColor': 'black',
        '--toastBarBackground': 'olive'
     } })
}
const success_notify = (metin) => {
    toast.push(metin, { theme: { 
        '--toastBackground': 'green',
        '--toastColor': 'white',
        '--toastBarBackground': 'olive'
     } })
}


export const myfetch_json = (func,...params) => {
    loading.set(1) 
    return func(...params)
    // for debug  output
    .then(r=>{
        console.log(r);
        return r;
    }) 
    // debug end
    .then(r=>JSON.parse(r))
    .then(r=>{ 
        if("message" in r) {
            if( r.message.Type == 0 ) {
                warning_notify(r.message.Message)
            } else if( r.message.Type == 1 ) {
                success_notify(r.message.Message)
            }
        }

        return r.data
    })
    .catch((err) => {
        console.error(err);
      })
    .finally(()=>{
        loading.set(0)
    })
}