import React, { Component } from "react";
import { Toast } from "react-bootstrap";
import { cap } from "../util";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTired } from "@fortawesome/free-solid-svg-icons";

class AftError extends Component {
  constructor(props) {
    super(props);
    this.state = { error: "" };
    this.handleDismiss = this.handleDismiss.bind(this);
  }

  handleDismiss() {
    this.setState({ error: "" });
  }

  componentDidUpdate(prevProps) {
    if (prevProps.error !== this.props.error) {
      this.setState({ error: this.props.error });
    }
  }

  render() {
    if (this.state.error !== "") {
      return (
        <Toast
          style={{
            position: "absolute",
            top: 80,
            left: "50%",
            marginLeft: "-150px",
            width: "300px"
          }}
          onClose={() => this.handleDismiss()}
        >
          <Toast.Header>
            <FontAwesomeIcon icon={faTired} />
            <strong className="mr-auto" style={{ marginLeft: "10px" }}>
              {cap(this.state.error.name)}
            </strong>
          </Toast.Header>
          <Toast.Body>{this.state.error.message}</Toast.Body>
        </Toast>
      );
    }
    return <div></div>;
  }
}

export default AftError;
