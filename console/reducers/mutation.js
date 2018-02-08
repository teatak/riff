import Common from './common'

export const MUTATIONSERVICE_REQUEST = 'MUTATIONSERVICE_REQUEST';
export const MUTATIONSERVICE_SUCCESS = 'MUTATIONSERVICE_SUCCESS';
export const MUTATIONSERVICE_FAILURE = 'MUTATIONSERVICE_FAILURE';

//获取Product
export const mutationService = (services, cmd) => (dispatch, getState) => {
    let servicesList = [];
    Object.keys(services).map((key, index) => {
        let service = services[key];
        servicesList.push(`    {name:"` + service.name + `",ip:"` + service.ip + `",port:` + service.port + `,cmd:` + cmd + `},`);
    });
    const query = `mutation {
  mutationService(services:[
` + servicesList.join('\n') + `
  ]){
    name,
    ip,
    port,
    error,
    cmd,
    success
  }
}`;
    dispatch({type: MUTATIONSERVICE_REQUEST});
    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: MUTATIONSERVICE_SUCCESS,
                status: status,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: MUTATIONSERVICE_FAILURE,
                status: status,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};


const mutation = (state = {
    mutationService: Common.initRequest,
}, action) => {
    switch (action.type) {
        case MUTATIONSERVICE_REQUEST:
            return {
                ...state,
                mutationService: {
                    ...state.mutationService,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case MUTATIONSERVICE_SUCCESS:
            return {
                ...state,
                mutationService: {
                    ...state.mutationService,
                    loading: false,
                    status: 200,
                    error: null,
                    lastUpdated: action.receivedAt
                }
            };
        case MUTATIONSERVICE_FAILURE:
            return {
                ...state,
                mutationService: {
                    ...state.mutationService,
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

export default mutation