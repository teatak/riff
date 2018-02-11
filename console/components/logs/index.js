import React from 'react'
import {connect} from 'react-redux'
import {cancelLogs, getLogs} from "../../reducers/logs"
import Close from '../icons/close'
import SwapHoriz from '../icons/swapHoriz'
import Refresh from '../icons/refresh'

import './index.pcss'

const mapStateToProps = (state, ownProps) => {
    return {
        logs: state.logs
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getLogs: () => {
            dispatch(getLogs());
        },
    }
};

class Logs extends React.Component {
    constructor(props) {
        super(props);
    }

    componentWillMount() {
        this.props.getLogs();
    }

    componentDidMount() {
        this.el = document.getElementById("logs");
        this.handleScroll();
    }

    componentWillReceiveProps(nextProps) {
        this.handleScroll();
    }

    handleRefresh = () => {
        this.props.getLogs();
    };
    handleScroll = () => {
        if (this.el) {
            let scrolltop = this.el.scrollTop;
            let clientHeight = this.el.clientHeight;
            let scrollHeight = this.el.scrollHeight;
            setTimeout(() => {
                if (scrolltop + clientHeight === scrollHeight || scrolltop === 0) {
                    scrollHeight = this.el.scrollHeight;
                    this.el.scrollTop = scrollHeight
                }
            }, 17)
        }
    };

    handleClose = () => {
        if (this.props.onClose) {
            this.props.onClose()
        }
    };

    renderList() {
        const {logs} = this.props;
        return <ul id="logs">
            {logs.list.map((log, index) => {
                return <li className="item" key={index}>
                    <pre>
                        {log}
                    </pre>
                </li>
            })}
        </ul>
    }

    render() {
        const {logs} = this.props;
        let network = "network";
        let error = "";
        if (logs.fetchLogs.status === 500) {
            network += " error";
            error = logs.fetchLogs.error;
        }

        return <div className="logs">
            <div className="logs-toolbar">
                <SwapHoriz className={network}/>
                {error === "" ? null :
                    <div className="error">{error}<Refresh className="refresh" onClick={this.handleRefresh}/></div>}
                <Close className="close" onClick={this.handleClose}/>
            </div>
            {this.renderList()}
        </div>
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Logs)
