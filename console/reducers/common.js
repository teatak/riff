import Config from 'config'
import React from 'react';

class Common  {
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
    fetch = (cmd,cb) => {
        fetch(Config.api, {
            method: 'post',
            headers: { 'Content-Type': 'application/json' },
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
                if(error.response) {
                    cb(null, error.message, error.response.status);
                } else {
                    cb(null, error.message, 500);
                }
            });
    }
}

export default new Common()