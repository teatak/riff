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
        cancelLogs: () => {
            dispatch(cancelLogs());
        },
    }
};

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

class Logs extends React.Component {
    constructor(props) {
        super(props);

    }

    componentWillMount() {
        if (isIe) {
            return
        }
        this.props.getLogs();
    }

    componentDidMount() {
        this.el = document.getElementById("logs");
        if (isIe) {

        } else {
            this.handleScroll();
        }
    }

    componentWillReceiveProps(nextProps) {
        if (isIe) {
            return
        }
        this.handleScroll();
    }

    handleRefresh = () => {
        if (isIe) {
            return
        }
        this.props.getLogs();
    };
    handleScroll = () => {
        if (isIe) {
            return
        }
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
        if (isIe) {

        } else {
            this.props.cancelLogs();
        }
        if (this.props.onClose) {
            this.props.onClose()
        }
    };

    renderList() {
        if (isIe) {
            return <iframe className="iframe-logs" src="/api/logs" id="logs"/>
        } else {
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
    }

    render() {
        const {logs} = this.props;
        let network = "network";
        let error = "";
        if (logs.fetchLogs.status === 500) {
            network += " error";
            error = logs.fetchLogs.error;
        }
        let renderIe = "";
        if (isIe) {
            renderIe = "This is an iframe under IE or Edge";
        }
        return <div className="logs">
            <div className="logs-toolbar">
                <SwapHoriz className={network}/>
                {error === "" ? null :
                    <div className="error">{error}<Refresh className="refresh" onClick={this.handleRefresh}/></div>}
                {renderIe === "" ? null :
                    <div className="error ie">{renderIe}</div>}
                <Close className="close" onClick={this.handleClose}/>
            </div>
            {this.renderList()}
        </div>
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Logs)
