import React, {Component} from 'react'

class Icon extends Component {

    static defaultProps = {
        onMouseEnter: () => {
        },
        onMouseLeave: () => {
        },
        viewBox: '0 0 24 24',
    };

    state = {
        hovered: false,
    };

    handleMouseLeave = (event) => {
        this.setState({hovered: false});
        this.props.onMouseLeave(event);
    };

    handleMouseEnter = (event) => {
        this.setState({hovered: true});
        this.props.onMouseEnter(event);
    };

    render() {
        const {
            children,
            color,
            hoverColor,
            onMouseEnter, // eslint-disable-line no-unused-vars
            onMouseLeave, // eslint-disable-line no-unused-vars
            style,
            viewBox,
            title,
            ...other
        } = this.props;


        const offColor = color ? color : 'currentColor';
        const onColor = hoverColor ? hoverColor : offColor;

        return (
            <svg
                {...other}
                onMouseEnter={this.handleMouseEnter}
                onMouseLeave={this.handleMouseLeave}
                style={style}
                viewBox={viewBox}
            >
                {title?<title>{title}</title>:null}
                {children}
            </svg>
        );
    }
}

export default Icon;