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
                cb(json, null)
            })
            .catch(error => {
                cb(null, error)
            });
    }
}

export default new Common()