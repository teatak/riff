import {combineReducers} from 'redux'
import nodes from './nodes'
import logs from './logs'

const reducers = combineReducers({
    nodes,
    logs,
});

export default reducers