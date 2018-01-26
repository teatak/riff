import React from 'react'
import {connect} from 'react-redux'
import {getLogs} from "../../reducers/logs";

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
        this.props.getLogs()
    }

    componentDidMount() {
        this.el = document.getElementById("logs");
    }

    componentWillReceiveProps(nextProps) {
        if (this.el) {
            let scrolltop = this.el.scrollTop;
            let clientHeight = this.el.clientHeight;
            let scrollHeight = this.el.scrollHeight;
            if (scrolltop + clientHeight === scrollHeight && scrolltop !== 0) {
                setTimeout(() => {
                    this.el.scrollTop = scrollHeight
                }, 17)
            }
        }
    }

    handleClose = () => {
      if(this.props.onClose) {
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
        return <div className="logs">
            <div onClick={this.handleClose} className="close">close</div>
            {this.renderList()}
        </div>
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Logs)
