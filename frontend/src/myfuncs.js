import { loading } from "./store"
import {  toast } from '@zerodevx/svelte-toast'



export const warning_notify = (metin) => {
    toast.push(metin, { theme: { 
        '--toastBackground': 'yellow',
        '--toastColor': 'black',
        '--toastBarBackground': 'olive'
     } })
}
export const success_notify = (metin) => {
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

/* 
* copied from 
* https://stackoverflow.com/questions/15900485/correct-way-to-convert-size-in-bytes-to-kb-mb-gb-in-javascript
*/
export function formatBytes(bytes, decimals = 2) {
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}
/* 
copied from
https://css-tricks.com/snippets/javascript/htmlentities-for-javascript/
*/
export function htmlEntities(str) {
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}