import Config from 'config'
import React from 'react';

let isIe = false;

if (/MSIE/i.test(navigator.userAgent)) {
    isIe = true;
}
if (/rv:11.0/i.test(navigator.userAgent)) {
    isIe = true;
}
if (/Edge\/\d./i.test(navigator.userAgent)) {
    isIe = false;
}

class Common {
    initRequest = {
        loading: false,
        status: 0,
        lastUpdated: "2000-1-1",    //请求的结果时间
    };
    checkStatus = (response) => {
        if (response.status >= 200 && response.status < 300) {
            return response
        } else if (response.status === 500) {
            return response
        } else {
            let error = new Error(response.statusText);
            error.response = response;
            throw error;
        }
    };
    parseJSON = (response) => {
        return response.json()
    };
    fetch = (cmd, cb) => {
        fetch(Config.api, {
            method: 'post',
            headers: {'Content-Type': 'application/json'},
            credentials: "include",
            body: JSON.stringify(cmd),
        })
            .then(this.checkStatus)
            .then(this.parseJSON)
            .then(json => {
                if (json.errors) {
                    if (json.errors[0].message === "NOT_FOUND") {
                        cb(json, json.errors[0].message, 404)
                    } else {
                        cb(json, json.errors[0].message, 500)
                    }
                } else {
                    cb(json, null, 200)
                }
            })
            .catch(error => {
                if (error.response) {
                    cb(null, error.message, error.response.status);
                } else {
                    cb(null, error.message, 500);
                }
            });
    };
    utf8ArrayToStr = (array) => {
        let out, i, len, c;
        let char2, char3;

        out = "";
        len = array.length;
        i = 0;
        while (i < len) {
            c = array[i++];
            switch (c >> 4) {
                case 0:
                case 1:
                case 2:
                case 3:
                case 4:
                case 5:
                case 6:
                case 7:
                    // 0xxxxxxx
                    out += String.fromCharCode(c);
                    break;
                case 12:
                case 13:
                    // 110x xxxx   10xx xxxx
                    char2 = array[i++];
                    out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
                    break;
                case 14:
                    // 1110 xxxx  10xx xxxx  10xx xxxx
                    char2 = array[i++];
                    char3 = array[i++];
                    out += String.fromCharCode(((c & 0x0F) << 12) |
                        ((char2 & 0x3F) << 6) |
                        ((char3 & 0x3F) << 0));
                    break;
            }
        }

        return out;
    };
    isIe = () => {
        return isIe
    };
}

export default new Common()