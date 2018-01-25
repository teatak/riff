import React from 'react';

const Common = {
    initRequest : {
        loading: false,
        status: 0,
        lastUpdated: "2000-1-1",    //请求的结果时间
    },
    checkStatus : (response) => {
        if (response.status >= 200 && response.status < 300) {
            return response
        } else {
            let error = new Error(response.statusText);
            error.response = response;
            throw error;
        }
    },
    parseJSON : (response) => {
        return response.json()
    }
};

export default Common