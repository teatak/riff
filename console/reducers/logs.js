import Common from './common'
import Config from 'config'

export const LOG_REQUEST = 'LOG_REQUEST';
export const LOG_SUCCESS = 'LOG_SUCCESS';
export const LOG_FAILURE = 'LOG_FAILURE';

export const getLogs = () => (dispatch, getState) => {
    let state = getState();
    if (state.logs.fetchLogs.loading) {
        return
    }
    dispatch({type: LOG_REQUEST});

    // let reader = null;
    // let loop = (response) => {
    //     setTimeout(() => {
    //         reader.read().then((b) => {
    //             let text = new TextDecoder("utf-8").decode(b.value);
    //             dispatch({
    //                 type: LOG_SUCCESS,
    //                 status: response.status,
    //                 text,
    //                 receivedAt: Date.now()
    //             });
    //             if (b.done) {
    //                 init = false;
    //             }
    //             loop(response);
    //         }).catch((error) => {
    //             dispatch({
    //                 type: LOG_FAILURE,
    //                 status: 500,
    //                 error: error,
    //                 receivedAt: Date.now()
    //             });
    //         })
    //     }, 0);
    // };

    // $.ajax({
    //     type: "get",
    //     url: Config.api + "/logs",
    //     dataType: "text",
    //     xhr: function () {
    //         let xhr = $.ajaxSettings.xhr();
    //         xhr.onprogress = (e) => {
    //             // For downloads
    //             dispatch({
    //                 type: LOG_SUCCESS,
    //                 status: e.currentTarget.status,
    //                 text: e.currentTarget.responseText,
    //                 receivedAt: Date.now()
    //             });
    //         };
    //         return xhr;
    //     }
    // }).fail((xhr) => {
    //     dispatch({
    //         type: LOG_FAILURE,
    //         status: 500,
    //         error: xhr.statusText,
    //         receivedAt: Date.now()
    //     });
    // });

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
        let pump = () => {
            return reader.read().then(({done, value}) => {
                if (done) {
                    return
                }
                dispatch({
                    type: LOG_SUCCESS,
                    status: 200,
                    text: Utf8ArrayToStr(value),
                    receivedAt: Date.now()
                });
                return pump();
            })
        };
        return pump();
    };

    fetch(Config.api + "/logs", {
        method: 'get',
        headers: {'connection': 'keep-alive'},
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