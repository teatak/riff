import Common from '../common'
import Config from 'config'

export const LOG_RESET = 'LOG_RESET';
export const LOG_REQUEST = 'LOG_REQUEST';
export const LOG_SUCCESS = 'LOG_SUCCESS';
export const LOG_FAILURE = 'LOG_FAILURE';

let initReader = null;

export const cancelLogs = () => (dispatch, getState) => {
    let state = getState();
    if (state.logs.fetchLogs.loading) {
        if (initReader !== null) {
            initReader.cancel();
            dispatch({type: LOG_RESET});
        }
    }
};

export const getLogs = () => (dispatch, getState) => {
    let state = getState();
    if (state.logs.fetchLogs.loading) {
        return
    }
    dispatch({type: LOG_REQUEST});

    let consume = (reader) => {
        initReader = reader;
        let total = 0;
        return new Promise((resolve, reject) => {
            function pump() {
                reader.read().then(({done, value}) => {
                    if (done) {
                        resolve();
                        return
                    }
                    total += value.byteLength;
                    //console.log(`received ${value.byteLength} bytes (${total} bytes in total)`)
                    let text = Common.utf8ArrayToStr(value);
                    dispatch({
                        type: LOG_SUCCESS,
                        status: 200,
                        text: text,
                        receivedAt: Date.now()
                    });
                    pump()
                }).catch(reject)
            }

            pump()
        })
    };

    fetch(Config.api + "/logs", {
        method: 'get',
        headers: {'connection': 'keep-alive'},
        credentials: 'include'
    }).then((response) => {
        return consume(response.body.getReader())
    }).then(() => {
        dispatch({
            type: LOG_FAILURE,
            status: 500,
            error: "Server connect closed",
            receivedAt: Date.now()
        });
    }).catch((error) => {
        dispatch({
            type: LOG_FAILURE,
            status: 500,
            error: error.message,
            receivedAt: Date.now()
        });
    })
};

const logs = (state = {
    fetchLogs: Common.initRequest,
    list: [],                   //数据
}, action) => {
    switch (action.type) {
        case LOG_RESET:
            return {
                ...state,
                fetchLogs: Common.initRequest,
                list: [],
            };
        case LOG_REQUEST:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case LOG_SUCCESS:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
                    status: 200,
                    error: null,
                    lastUpdated: action.receivedAt
                },
                list: [...state.list, action.text]
            };
        case LOG_FAILURE:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
                    loading: false,
                    status: action.status,
                    error: action.error,
                    lastUpdated: action.receivedAt
                },
            };
        default:
            return state
    }
};

export default logs