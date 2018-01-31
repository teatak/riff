import Common from './common'
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

    let Utf8ArrayToStr = (array) => {
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
                    let text = Utf8ArrayToStr(value);
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